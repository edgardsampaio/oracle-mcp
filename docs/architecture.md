# Arquitetura

Servidor MCP em Go para consultas e relatórios de BI em Oracle, somente leitura. Identidades, OAuth e auditoria ficam em SQLite local. O código é independente e pode ser entregue a terceiros: não importa projetos privados nem copia tools, SQL ou schemas de outros produtos.

A especificação de origem é `oracle-mcp-project.md`. Este documento fixa a arquitetura mínima. Pacotes só nascem na fase que os usa.

## Fluxo

```text
Cliente MCP / Agente
        |
      HTTPS
        |
   HTTP Router
        |
  AuthN (bearer) + AuthZ (ACL)
        |
   MCP Tool Registry
        |
   Dimensão ou fato (SQL no pacote)
        |
  Executor Oracle (SQL do pacote)
        |
 Oracle (usuário somente leitura)

API administrativa -- Auth REST -- IdentityStore -- SQLite
                                      \
                               Authorization server
                               (Authorization Code + PKCE,
                                CIMD e DCR)
```

O login REST administrativo e o OAuth do MCP são fluxos de emissão diferentes. Os dois produzem access token deste mesmo authorization server, e o endpoint `/mcp` aceita esse token tanto para cliente CIMD quanto para cliente DCR. O papel `admin` administra usuários e não libera dados do Oracle.

## Fronteiras

| Pacote | Responsabilidade | Não importa |
| --- | --- | --- |
| `internal/config` | YAML opcional, variáveis de ambiente, validação na subida | Oracle, MCP, identidade |
| `internal/api` | Router, health, handlers e middleware da API `/api/v1` | SQL de relatório |
| `internal/auth` | Senha Argon2id, tokens opacos, authorization server | Oracle, dimensão, fato |
| `internal/identity` | Store de usuários, papéis, OAuth e auditoria | Oracle, ferramentas |
| `internal/authorization` | ACL por usuário, ferramenta, dataset e relatório | Driver Oracle |
| `internal/database/oracle` | Pool, consulta parametrizada, normalização e limites | Usuários, OAuth |
| `internal/bi` | Auditoria da tool e a interface da consulta | Hash de senha, SQL de negócio |
| `internal/dimensao` | Uma tool por dimensão, com a query no pacote | Fato, outra dimensão |
| `internal/fato` | Um subpacote por relatório, com a query sem rótulo e a query com JOIN | SQL de dimensão |
| `internal/mcp` | Transporte Streamable HTTP e registro de tools | SQL solto, detalhe de tabela de identidade |
| `internal/audit` | Evento com ator, ação, alvo, outcome e duração | Senha, token, SQL, linhas retornadas |

Identidade não conhece Oracle. Oracle não conhece usuários. Ferramenta de negócio não recebe a identidade do agente como argumento: o ator sai do token já validado.

Pacotes da tabela entram na fase correspondente. Config, router, health, transporte MCP, identidade, o adaptador Oracle, as dimensões e os fatos já existem. A linha visível de cada fato sai de `DIM_RLS_CENTRO_CUSTO`, com o e-mail gravado no usuário.

## HTTP

Desenvolvimento escuta em `127.0.0.1`. HTTPS remoto termina em proxy reverso confiável, com origem e cabeçalhos encaminhados validados.

`read_timeout` limita a leitura de cabeçalho e corpo no `http.Server`. `write_timeout` fica reservado: o transporte stateful mantém um GET de eventos aberto, e um prazo de escrita no servidor encerraria essa resposta.

| Método | Rota | Desde | Acesso |
| --- | --- | --- | --- |
| GET | `/api/v1/health/live` | Fase 1 | Processo no ar |
| GET | `/api/v1/health/ready` | Fase 1 | Pronto para tráfego. Sem Oracle configurado responde o processo. Com host configurado, faz ping e, se falhar, responde 503 sem detalhe interno |
| POST | `/mcp` | Fase 1 | Streamable HTTP. A partir da Fase 2 exige bearer |
| GET | `/.well-known/oauth-protected-resource` | Fase 2 | Público. RFC 9728 |
| GET | `/.well-known/oauth-authorization-server` | Fase 2 | Público. RFC 8414 |
| GET | `/oauth/authorize` | Fase 2 | Dono do recurso, depois do login |
| POST | `/oauth/token` | Fase 2 | Authorization Code + PKCE |
| POST | `/oauth/register` | Fase 2 | DCR, RFC 7591 |
| POST | `/api/v1/auth/login` | Fase 2 | Público, com rate limit |
| POST | `/api/v1/auth/refresh` | Fase 2 | Refresh válido |
| POST | `/api/v1/auth/logout` | Fase 2 | Autenticado |
| GET, POST, PATCH, DELETE | `/api/v1/users` e sub-rotas da spec | Fase 2 | `users:read` ou `users:manage`, conforme a rota |
| GET | `/api/v1/roles` | Fase 2 | `users:read` |
| PUT | `/api/v1/users/{id}/roles` | Fase 2 | `users:manage` |

Não há endpoint público de bootstrap. Na primeira execução, sem administrador ativo, o processo pede usuário, senha e confirmação no terminal e só então escuta. O comando `oracle-mcp admin create` faz o mesmo pedido e recusa um segundo administrador.

## MCP

SDK oficial `github.com/modelcontextprotocol/go-sdk` v1.7.x. Transporte Streamable HTTP em `/mcp`. Stdio fica fora da v1.

O protocolo mais recente suportado pelo SDK é `2026-07-28`. Sobre HTTP, essa versão exige modo stateless. O padrão da v1 é sessão stateful, e clientes mais antigos negociam para `2025-11-25` ou anterior. O modo stateless só é ligado se um cliente-alvo exigir `2026-07-28`.

Na Fase 1 o handshake lista as ferramentas registradas. As ferramentas de negócio entram na Fase 4.

A partir da Fase 2, `/mcp` é resource server:

- Metadados de recurso protegido pelo helper do SDK.
- Todo pedido traz `Authorization: Bearer`.
- O verificador confere token opaco (hash no SQLite), expiração, revogação e audience RFC 8707 igual à URL canônica do MCP.
- A ACL da ferramenta roda depois da autenticação. Token válido não autoriza todas as tools.

## Authorization server

Um único grant: Authorization Code com PKCE S256. Não há password grant nem client credentials no lugar desse fluxo. Access token curto, refresh com rotação, ambos só como hash.

O cliente entra nesta ordem: pré-registro que já possua, CIMD, DCR. Pré-registro administrativo pode existir para um cliente conhecido. Cliente novo entra sem cadastro manual, por um dos dois registros abaixo.

### CIMD

Caminho preferido da spec MCP `2026-07-28`. O `client_id` é uma URL `https` com path. O servidor busca o JSON e exige `client_id`, `client_name` e `redirect_uris`. O `client_id` do documento tem de ser exatamente a URL buscada, e o `redirect_uri` do pedido tem de estar na lista. Os metadados anunciam `client_id_metadata_document_supported: true`. Cache segue os cabeçalhos HTTP. A busca tem timeout, tamanho máximo e bloqueio de SSRF: só HTTPS, sem IP privado, loopback, link-local ou metadados de nuvem. CIMD não emite `client_id`.

### DCR

RFC 7591, para clientes anteriores. Os metadados publicam `registration_endpoint`. O POST grava o cliente em `oauth_clients` e devolve um `client_id` deste issuer, que não vale em outro authorization server. `application_type` é aceito, inclusive `native`. O retorno pode ser https (Claude, ChatGPT e o callback web do Cursor), http em localhost, `127.0.0.1`, `::1` ou IP privado, ou esquema privado de aplicativo, inclusive `cursor://`. Na autorização, http local ignora a porta; `localhost`, `127.0.0.1` e `::1` contam como o mesmo host.

A biblioteca do authorization server entra na Fase 2 só se implementar CIMD e DCR. Caso contrário, esses fluxos são escritos no padrão das RFCs, sem um framework OAuth genérico.

## Identidade

SQLite via `modernc.org/sqlite` em `data/identity.db`. Migrações SQL embutidas, com tabela de versão. `foreign_keys` ligado, WAL avaliado e `busy_timeout` definido. O arquivo não entra no Git.

Tabelas: `users`, `roles`, `permissions`, `user_roles`, `role_permissions`, `access_tokens`, `refresh_tokens`, `oauth_clients`, `oauth_authorization_codes`, `audit_log`, `login_attempts`. `access_tokens` existe porque o bearer é opaco e só o hash pode ser guardado.

Papéis iniciais: `admin`, `analyst`, `viewer`. Permissões iniciais de API: `users:read`, `users:manage`. Permissões de relatório são códigos associados a uma ferramenta concreta; datasets não são inventados.

Senha com Argon2id, salt individual e parâmetros versionados. Senha temporária exige troca no primeiro acesso. A senha do bootstrap é lida no terminal, sem eco e sem argumento de linha de comando.

## Oracle e tools

O executor está em `internal/database/oracle`, com `go-ora` v2.9.0 (MIT) como driver provisório. `godror` entra se a versão, o TLS ou os tipos do ambiente exigirem. Credencial dedicada de leitura. Só SQL escrito no pacote da tool, parametrizado. Timeout, máximo de linhas e máximo de bytes. `NUMBER` que precise de precisão sai como string. `DATE`, `TIMESTAMP`, `CLOB` e `NULL` têm normalização explícita. O log da consulta traz o id, a duração, a quantidade de linhas e o outcome.

Não há catálogo YAML e não existe tool de SQL livre. Cada dimensão é um pacote em `internal/dimensao`. Cada fato é um subpacote em `internal/fato`, com as duas queries copiadas do contrato do cliente. Dimensão não importa fato, e fato não importa dimensão.

O e-mail da segurança por linha não entra no JSON da tool. O servidor lê a coluna `email` do usuário do token e faz o bind `:email`. Usuário sem e-mail não consulta fato.

## Dimensões e fatos

Há uma tool por dimensão do contrato, inclusive `dim_rls_centro_custo`, que só devolve os centros de custo do próprio usuário. `dim_natureza` substitui o nome antigo `dim_natureza_financeira`. A resposta traz a hierarquia e as flags que o SQL do cliente seleciona. Filtro omitido chega ao Oracle como `NULL`. O `limite` omitido vale 500 e não passa de `bi.max_rows`.

Cada fato publica duas tools. `fato_<assunto>` é o grão sem rótulos. `fato_<assunto>_com_dimensoes` é o mesmo grão com `LEFT JOIN` nas dimensões. As duas levam o `EXISTS` em `DIM_RLS_CENTRO_CUSTO`. A chave é `ID Centro de Custo`, `ID Contrato` no orçamento e na confissão de dívida, ou o histórico do colaborador na folha.

## O que esta arquitetura não faz

- Não copia código, modelo ou regra de negócio de produtos privados.
- Não cria framework SQL genérico nem biblioteca compartilhada para outro banco.
- Não grava no Oracle e não guarda usuário ou token lá.
- Não trata acesso ao BI existente como autorização no MCP.
- Não expõe o MCP sem bearer a partir da Fase 2.
