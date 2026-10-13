# Oracle MCP — Especificação do Projeto

**Versão:** 1.0  
**Data:** 09/10/2026  
**Destino:** `C:\workspace\go-projects\oracle-mcp`  
**Linguagem:** Go  
**Status:** especificação para planejamento e implementação no Cursor

## 1. Objetivo

Desenvolver um servidor MCP (Model Context Protocol) independente para uma empresa terceira, que exponha de maneira controlada consultas e relatórios de um ambiente de Business Intelligence armazenado em Oracle. O código-fonte poderá ser entregue à empresa contratante.

O projeto **não é um produto da Ideia Tecnologia**. O `ideia-mcp` e o `ideia-copilot` podem servir de inspiração conceitual, mas **não devem ter código, modelos, regras de negócio, bibliotecas privadas, arquivos ou estruturas proprietárias copiados** sem autorização expressa e verificação dos direitos de reutilização.

C:\workspace\go-projects\ideia-go\apps\ideia-copilot



O desenho deve permitir que, no futuro, um projeto independente `sqlserver-mcp` seja implementado com os mesmos princípios, sem criar agora um framework genérico excessivo ou uma dependência entre entregáveis de clientes diferentes.

## 2. Escopo



### Incluído na primeira versão

- Aplicação Go independente, executável e configurável.
- Servidor HTTP com rotas versionadas e endpoints de saúde.
- Servidor MCP com transporte HTTP compatível com os clientes-alvo.
- Autenticação para API administrativa e mecanismo OAuth compatível com clientes MCP remotos.
- Gerenciamento de usuários e permissões por API REST.
- Armazenamento local de identidades em SQLite, sem instalação de servidor de banco.
- Comando local para bootstrap do primeiro administrador.
- Adaptador Oracle exclusivamente para leitura.
- Catálogo explícito de datasets e relatórios autorizados.
- Ferramentas MCP para descoberta e execução de relatórios.
- Auditoria, logs, testes e documentação de deploy.



### Fora do escopo inicial

- Cópia ou migração do código do Ideia MCP.
- Interface web administrativa.
- Geração livre de SQL pela IA.
- Escrita no Oracle, inclusive criação de tabelas de usuários ou tokens.
- Integração obrigatória com AD/LDAP/SSO, até descobrir a identidade utilizada no BI atual.
- Framework SQL universal, ORM próprio, multi-tenant distribuído e clusterização.
- Implantação do SQL Server MCP nesta entrega.



## 3. Princípios de arquitetura

1. **Independência:** nenhuma dependência de código privado da Ideia.
2. **Simplicidade:** usar bibliotecas públicas e abstrações apenas quando houver uma necessidade real.
3. **Separação de responsabilidades:** identidade, autorização, transporte MCP, domínio de BI e acesso Oracle isolados.
4. **Read-only:** o banco Oracle é uma fonte de dados, não o repositório de identidades.
5. **Menor privilégio:** cada usuário e cada ferramenta têm acesso apenas ao necessário.
6. **Contratos estáveis:** ferramentas e APIs com entradas e saídas documentadas.
7. **Evolução incremental:** primeiro servidor/autenticação, depois Oracle e ferramentas de BI.
8. **Sem inferência de permissões:** ter acesso ao BI existente não significa automaticamente ter o mesmo escopo no MCP.



## 4. Arquitetura lógica

```text
Cliente MCP / Agente de IA
          |
       HTTPS
          |
   HTTP Router / MCP Endpoint
          |
  AuthN (OAuth) + AuthZ (ACL)
          |
   MCP Tool Registry
          |
   BI Catalog / Reports
          |
  Oracle Query Executor
          |
 Oracle (usuário somente leitura)

API administrativa ---- Auth Service ---- IdentityStore ---- SQLite local
                                      \\-- outro provider (futuro)
```

**Importante:** o OAuth do protocolo MCP e o login REST administrativo são fluxos relacionados, mas não equivalentes. O servidor deve adotar os requisitos da versão do MCP escolhida e verificar compatibilidade com os clientes concretos antes da implementação. Não criar um protocolo OAuth artesanal.

## 5. Estrutura de diretórios sugerida

```text
oracle-mcp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── api/
│   │   ├── router.go
│   │   ├── handlers/
│   │   └── middleware/
│   ├── auth/
│   │   ├── service.go
│   │   ├── password.go
│   │   ├── tokens.go
│   │   └── oauth/
│   ├── identity/
│   │   ├── store.go
│   │   ├── models.go
│   │   └── sqlite/
│   ├── authorization/
│   ├── database/
│   │   ├── executor.go
│   │   └── oracle/
│   │       ├── connection.go
│   │       ├── executor.go
│   │       └── metadata.go
│   ├── bi/
│   │   ├── catalog/
│   │   └── reports/
│   ├── mcp/
│   │   ├── server.go
│   │   └── tools/
│   └── audit/
├── migrations/
│   └── identity/
├── configs/
│   └── config.example.yaml
├── docs/
│   ├── architecture.md
│   ├── security.md
│   ├── api.md
│   ├── deployment.md
│   └── decisions.md
├── data/                 # ignorado no Git
├── .env.example
├── .gitignore
├── AGENTS.md
├── go.mod
└── README.md
```

A estrutura é uma referência, não uma exigência de criar todos os arquivos antecipadamente. Evitar pacotes e interfaces vazios.

## 6. Tecnologias e dependências


| Responsabilidade  | Opção recomendada                              | Observação                                                            |
| ----------------- | ---------------------------------------------- | --------------------------------------------------------------------- |
| Linguagem         | Go estável                                     | Fixar versão no `go.mod`                                              |
| HTTP              | `net/http` ou router público leve              | Evitar framework pesado sem necessidade                               |
| MCP               | SDK Go público compatível com versão escolhida | Validar transporte, autenticação e clientes-alvo                      |
| Oracle            | `database/sql` + `go-ora`                      | Avaliar compatibilidade de versões e tipos; `godror` como alternativa |
| SQL Server futuro | `go-mssqldb`                                   | Em projeto separado                                                   |
| SQLite            | `modernc.org/sqlite`                           | Go puro; não exige instalação de serviço SQLite                       |
| Senhas            | Argon2id (`golang.org/x/crypto`)               | Salt individual, parâmetros versionados                               |
| Identificadores   | UUID                                           | Usar biblioteca pública consolidada                                   |
| Logs              | `log/slog`                                     | Estruturados e sem segredos                                           |
| Configuração      | Variáveis de ambiente + YAML opcional          | Segredos fora do Git                                                  |
| Testes            | `testing`                                      | Integração Oracle opcional, com credenciais próprias                  |


**Nota:** não decidir definitivamente entre `go-ora` e `godror` sem conhecer versão do Oracle, método de autenticação, requisitos de TLS, tipos utilizados e ambiente de deploy.

## 7. IdentityStore e SQLite

O IdentityStore deve encapsular operações de usuários, papéis, permissões, sessões e estado de OAuth, sem conhecer Oracle. Inicialmente será implementado em SQLite local (`data/identity.db`).

Tabelas iniciais sugeridas:

- `users`: id, username único, display_name, password_hash, active, must_change_password, created_at, updated_at.
- `roles`: id, code único, description.
- `permissions`: id, code único, description.
- `user_roles`: user_id, role_id.
- `role_permissions`: role_id, permission_id.
- `refresh_tokens`: id, user_id, client_id, token_hash, expires_at, revoked_at, created_at, rotated_from_id.
- `oauth_clients`: informações de registro e configuração de clientes, conforme fluxo adotado.
- `oauth_authorization_codes`: hash/código, cliente, usuário, desafio PKCE, escopos, validade e uso único, se aplicável.
- `audit_log`: id, timestamp, actor_id, action, target, outcome, correlation_id e metadados não sensíveis.

Índices, restrições de integridade e migrações versionadas são obrigatórios. Evitar armazenar tokens em texto puro. Não incluir o arquivo do SQLite em repositórios, imagens públicas ou backups desprotegidos.

Configurações operacionais:

- Criar o diretório de dados com permissões restritas.
- Ativar `foreign_keys`, avaliar WAL e definir `busy_timeout`.
- Fazer backup consistente (não copiar somente o `.db` enquanto WAL estiver ativo).
- Executar em disco local persistente, não em diretório temporário.
- Não compartilhar o mesmo arquivo por várias máquinas via rede.
- Se houver múltiplas instâncias, migrar o IdentityStore para um repositório centralizado.



## 8. Autenticação e autorização



### API administrativa

- Login com usuário e senha, hash Argon2id.
- Sessões ou access tokens de curta duração; refresh tokens com rotação e revogação quando utilizados.
- Rate limiting, proteção contra brute force e auditoria de tentativas.
- Senha temporária com troca obrigatória no primeiro acesso.
- Não expor senha, hash, segredo OAuth ou tokens em logs e respostas.
- Separar alteração da própria senha de redefinição administrativa.
- Nunca ter usuário/senha padrão em produção.



### MCP remoto

- Implementar o mecanismo de autorização exigido pela versão MCP adotada, preferindo bibliotecas consolidadas e fluxos OAuth com Authorization Code + PKCE para clientes públicos.
- Validar issuer, audience/resource, escopos, expiração, cliente e assinatura dos tokens conforme arquitetura escolhida.
- Se houver integração com provedor externo OIDC, definir claramente os papéis de resource server e authorization server; não presumir que JWT local é suficiente.
- Não usar fluxo de senha do usuário como substituto do OAuth do MCP.
- Considerar redirecionamentos permitidos, registro de clientes e metadados de descoberta conforme o perfil de cliente suportado.



### Papéis iniciais


| Papel     | Capacidade                                                   |
| --------- | ------------------------------------------------------------ |
| `admin`   | Gerenciar usuários, permissões e configuração administrativa |
| `analyst` | Executar ferramentas e relatórios explicitamente autorizados |
| `viewer`  | Consultar relatórios limitados autorizados                   |


O papel `admin` **não** deve conceder acesso irrestrito aos dados do Oracle. A autorização deve ser avaliada por usuário, ferramenta, dataset, relatório e, quando aplicável, dimensão de negócio (empresa, filial, departamento etc.). Se o BI original aplicar segurança por linha, documentar como ela será preservada ou reproduzida antes do piloto.

## 9. Bootstrap do administrador

Criar comando CLI interativo:

```powershell
.\oracle-mcp.exe admin create
```

Requisitos:

- Solicitar senha sem ecoar caracteres.
- Não receber senha por argumento de linha de comando.
- Restringir a execução ao host e ao acesso ao arquivo local.
- Permitir bootstrap somente quando não existir administrador, com tratamento explícito para recuperação controlada.
- Registrar a operação sem registrar a senha.
- Não expor endpoint público de bootstrap.



## 10. API REST administrativa

Prefixo: `/api/v1`.


| Método | Rota                   | Acesso                                   | Objetivo                        |
| ------ | ---------------------- | ---------------------------------------- | ------------------------------- |
| POST   | `/auth/login`          | Público com rate limit                   | Autenticar                      |
| POST   | `/auth/refresh`        | Refresh válido                           | Renovar sessão                  |
| POST   | `/auth/logout`         | Autenticado                              | Revogar sessão                  |
| GET    | `/users`               | `users:read`                             | Listar usuários com paginação   |
| POST   | `/users`               | `users:manage`                           | Criar usuário                   |
| GET    | `/users/{id}`          | `users:read`                             | Consultar usuário               |
| PATCH  | `/users/{id}`          | `users:manage`                           | Atualizar usuário               |
| PUT    | `/users/{id}/password` | Próprio usuário ou admin, conforme fluxo | Alterar/redefinir senha         |
| DELETE | `/users/{id}`          | `users:manage`                           | Desativar usuário (soft delete) |
| GET    | `/roles`               | `users:read`                             | Listar papéis                   |
| PUT    | `/users/{id}/roles`    | `users:manage`                           | Atualizar papéis                |
| GET    | `/health/live`         | Política de exposição definida           | Liveness                        |
| GET    | `/health/ready`        | Restrito ou resposta mínima              | Readiness                       |


Exemplo de criação:

```http
POST /api/v1/users
Authorization: Bearer <access-token>
Content-Type: application/json
```

```json
{
  "username": "joao.silva",
  "name": "João Silva",
  "password": "<senha-temporaria>",
  "role": "analyst",
  "active": true
}
```

Retornar identificador, dados não sensíveis e indicação de troca obrigatória. Adotar respostas de erro consistentes, paginação, validação de entrada e códigos HTTP adequados. Documentar se a senha temporária será informada pelo administrador ou gerada por canal seguro; nunca enviá-la em logs.

## 11. Adaptador Oracle

Responsabilidades:

- Abrir e gerenciar pool de conexões usando `database/sql`.
- Configurar host, porta, service name, usuário, senha, timeouts e TLS quando aplicável.
- Validar conexão e apresentar health check sem expor dados internos.
- Executar somente consultas previamente cadastradas, parametrizadas e autorizadas.
- Tratar `NUMBER`, `DATE`, `TIMESTAMP`, `CLOB`, `NULL` e demais tipos necessários.
- Converter resultados para JSON previsível, preservando precisão numérica quando necessário.
- Impor timeout, máximo de linhas, limite de bytes e paginação adequada.
- Registrar tempo, identificador da consulta e resultado sem dados sensíveis.

**Controle de escrita:** utilizar credencial Oracle com privilégios efetivamente somente de leitura, preferencialmente em views autorizadas. Não confiar em filtros de strings ou regex para bloquear SQL perigoso. Não disponibilizar ferramenta de SQL arbitrário ao agente.

O acesso Oracle pode não permitir criação de objetos. Todas as tabelas de identidade, OAuth e auditoria ficam fora do Oracle, no SQLite.

## 12. Catálogo semântico de BI

A primeira etapa do domínio BI deve mapear os conceitos que já existem no sistema da empresa, por exemplo:

- Datasets/views disponíveis e seus campos.
- Significado das métricas, como faturamento e receita líquida.
- Campos de período e granularidade.
- Dimensões como filial, cliente, produto e região, quando existirem.
- Relacionamentos aprovados e filtros permitidos.
- Regras de segurança por linha/dimensão.
- Identificadores de relatórios e suas versões.

Não presumir nomes de tabelas nem fórmulas de negócio. Validar os números produzidos contra relatórios oficiais do BI.

## 13. Ferramentas MCP iniciais


| Ferramenta            | Descrição                                                      |
| --------------------- | -------------------------------------------------------------- |
| `bi_list_datasets`    | Lista datasets autorizados para o usuário                      |
| `bi_describe_dataset` | Descreve campos, tipos e filtros de um dataset                 |
| `bi_list_reports`     | Lista relatórios disponíveis e autorizados                     |
| `bi_describe_report`  | Descreve parâmetros, métricas e limites do relatório           |
| `bi_run_report`       | Executa consulta previamente aprovada com parâmetros validados |


Todas as ferramentas devem:

- Obter a identidade validada do contexto de autenticação, não de parâmetros enviados pelo agente.
- Validar autorização antes de consultar dados.
- Aplicar filtros de escopo do lado servidor.
- Limitar volume e tempo de execução.
- Produzir erros claros sem vazar SQL, credenciais ou estruturas internas.
- Gerar eventos de auditoria com identificador da ferramenta, usuário e duração.



## 14. Configuração

Exemplo ilustrativo de `configs/config.example.yaml`:

```yaml
server:
  address: "127.0.0.1:8080"
  public_base_url: "https://mcp.exemplo.com"
  read_timeout: "15s"
  write_timeout: "30s"

auth:
  identity_provider: "sqlite"
  sqlite_path: "./data/identity.db"
  access_token_ttl: "15m"
  refresh_token_ttl: "168h"

oracle:
  host: "oracle.exemplo.local"
  port: 1521
  service_name: "BI"
  username_env: "ORACLE_MCP_DB_USER"
  password_env: "ORACLE_MCP_DB_PASSWORD"
  query_timeout: "30s"
  max_open_connections: 10
  max_idle_connections: 5

bi:
  max_rows: 1000
  max_response_bytes: 2097152

audit:
  enabled: true
```

Esses valores são exemplos, não configurações aprovadas para produção. Segredos devem vir de variáveis de ambiente ou gerenciador de segredos; validar configurações obrigatórias na inicialização. HTTPS pode ser terminado por proxy reverso confiável, com regras explícitas de cabeçalhos e origem.

## 15. Segurança operacional

- HTTPS obrigatório em acessos remotos; não expor HTTP sem TLS diretamente na Internet.
- Validar origem, redirecionamentos OAuth e cabeçalhos encaminhados pelo proxy.
- Não registrar senhas, tokens, parâmetros sensíveis ou resultados completos de consultas.
- Usar usuário Oracle dedicado e privilégios mínimos.
- Restringir permissões de leitura/escrita do diretório `data/`.
- Configurar limites de requisições, concorrência e tamanho de payload.
- Evitar consultas longas e sobrecarga no banco de BI.
- Aplicar política de retenção de auditoria e backup.
- Definir procedimento de rotação de segredos e recuperação administrativa.
- Verificar licenças de bibliotecas e termos de entrega do código-fonte.
- Realizar revisão de ameaças antes de expor o MCP publicamente.



## 16. Estratégia de testes



### Unitários

- Hash e verificação de senhas.
- Permissões e escopos de acesso.
- Validação de parâmetros de ferramentas.
- Normalização de resultados Oracle simulados.
- Expiração, rotação e revogação de tokens.



### Integração local

- Migrações SQLite.
- Bootstrap administrativo.
- CRUD de usuários e papéis.
- Login, refresh, logout e troca de senha.
- Transporte MCP e descoberta das ferramentas.



### Integração Oracle (ambiente autorizado)

- Conectividade e autenticação.
- Tipos numéricos, datas, CLOB e NULL.
- Limites de tempo/linhas e cancelamento de contexto.
- Consultas parametrizadas e comparação com relatórios do BI.
- Permissões de leitura e impossibilidade de escrita com a credencial dedicada.



### Segurança

- Tentativas de acesso não autorizado.
- Escalada de privilégio e IDOR.
- Reutilização de refresh token revogado.
- Falhas de validação de parâmetros e de escopo.
- Logs sem segredos.



## 17. Plano incremental de implementação



### Fase 0 — Descoberta e decisões

- Identificar versão Oracle, sistema operacional e ambiente de deploy.
- Identificar como usuários se autenticam no BI atual (AD/LDAP/SSO, login próprio etc.).
- Identificar quais views/relatórios podem ser acessados e quais regras de autorização já existem.
- Confirmar clientes MCP-alvo, transporte e versão do protocolo.
- Validar licença e titularidade de todos os componentes entregáveis.

**Aceite:** decisões registradas em `docs/decisions.md`, sem código proprietário copiado.

### Fase 1 — Base Go e servidor MCP

- Criar módulo Go, configuração, logs, router, health checks e transporte MCP.
- Testes mínimos e documentação de execução.

**Aceite:** aplicação inicializa, endpoints respondem e cliente MCP consegue estabelecer conexão em ambiente de teste.

### Fase 2 — Identidade e API administrativa

- Implementar SQLite, migrações, bootstrap CLI, usuários, papéis e permissões.
- Implementar autenticação REST e OAuth compatível com clientes escolhidos.

**Aceite:** administrador criado localmente; API permite cadastrar, consultar, alterar e desativar usuários; sessões podem ser revogadas; testes passam.

### Fase 3 — Oracle

- Implementar adaptador, pool, consultas parametrizadas, normalização e limites.

**Aceite:** consultas aprovadas funcionam no ambiente de teste sem privilégios de escrita.

### Fase 4 — BI/MCP tools

- Implementar catálogo de datasets/relatórios e as cinco ferramentas iniciais.
- Validar métricas e escopos com o responsável pelo BI.

**Aceite:** usuário autorizado consulta relatório correto; usuário não autorizado é bloqueado; resultados batem com BI de referência.

### Fase 5 — Homologação e entrega

- Testes de carga, segurança, backup/restauração, configuração e deploy.
- Entregar código-fonte, documentação, licenças de terceiros, scripts e instruções de operação.

**Aceite:** a empresa consegue compilar, configurar, iniciar, administrar usuários e consultar relatórios sem dependências privadas da Ideia.

## 18. Evolução para SQL Server

Um futuro `sqlserver-mcp` deverá ser outro repositório independente. Reaproveitar conceitos, contratos documentados e eventualmente componentes cuja licença e titularidade permitam a reutilização. O adaptador SQL Server será implementado com driver próprio, respeitando sintaxe, tipos, autenticação, paginação e particularidades desse banco.

Não antecipar uma biblioteca compartilhada obrigatória. Considerar extração de um módulo comum somente após existir uma segunda implementação real e acordo de propriedade intelectual que permita o compartilhamento.

## 19. Prompt mestre para o Cursor

> Atue como arquiteto e desenvolvedor Go sênior. Leia integralmente este arquivo `oracle-mcp-project.md` e trate-o como especificação inicial do projeto em `C:\workspace\go-projects\oracle-mcp`.
>
> O código será entregue a uma empresa terceira. Não copie nem importe código proprietário de `ideia-mcp`, `ideia-copilot` ou de outros projetos privados. Use somente bibliotecas públicas com licenças compatíveis e implementação independente.
>
> **Primeira tarefa: planejamento, sem implementar código.** Inspecione o workspace e o estado atual do projeto. Identifique dúvidas, decisões em aberto, dependências e riscos. Proponha a menor arquitetura viável, estrutura de pacotes, interfaces necessárias e plano de commits/fases. Produza `docs/architecture.md`, `docs/decisions.md` e um checklist de implementação.
>
> Não invente tabelas Oracle, permissões do BI ou requisitos de identidade. Quando faltar informação, registre a pendência. Não crie abstrações vazias nem um framework SQL genérico. Não implemente SQL arbitrário para agentes. Use SQLite local para identidade e Oracle somente leitura.
>
> Apresente o plano e aguarde minha aprovação antes de implementar a Fase 1. Ao final de cada fase, execute os testes relevantes, documente o que foi feito, o que ficou pendente e solicite autorização para a próxima fase.



## 20. Perguntas pendentes para a empresa

1. Qual é a versão do Oracle e o sistema operacional do servidor?
2. O BI utiliza AD, LDAP, SSO, usuários Oracle ou autenticação própria?
3. Existem regras de segurança por usuário, filial, departamento ou empresa?
4. Quais views e relatórios poderão ser disponibilizados ao MCP?
5. O MCP será acessado apenas na rede interna, por VPN ou publicamente por HTTPS?
6. Quais clientes MCP devem funcionar na primeira entrega?
7. Qual é o volume esperado de usuários e consultas simultâneas?
8. Quem será responsável por backup, atualização e administração do serviço?
9. Quais são os termos de titularidade, reutilização e manutenção do código entregue?

---

**Orientação final:** entregar primeiro um MCP simples, independente, seguro e funcional. Generalizar somente aquilo que uma necessidade concreta justificar.