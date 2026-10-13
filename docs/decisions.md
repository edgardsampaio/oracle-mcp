# Decisões

Registro do que já está fechado e do que continua pendente. Pendência não é preenchida com suposição. A fase indicada é a que fica sem aceite enquanto a pergunta estiver aberta.

## Fechadas

### D1 — Módulo e linguagem

Path do módulo Go: `oracle-mcp`, até existir repositório de entrega. O `go.mod` fixa `go 1.25.0` porque o SDK v1.7.0 exige Go >= 1.25. Na máquina da Fase 1 o comando `go` era 1.23.5 e baixou uma toolchain compatível para compilar. HTTP com `net/http` ServeMux.

`read_timeout` alimenta `ReadHeaderTimeout` e `ReadTimeout` do `http.Server`. `write_timeout` é validado e guardado, e não é aplicado ao servidor HTTP: o GET de eventos do transporte stateful permanece aberto.

### D2 — SDK e transporte MCP

`github.com/modelcontextprotocol/go-sdk` v1.7.0. Streamable HTTP em `/mcp`. Stdio fora da v1. Sessão stateful por padrão. `2026-07-28` sobre HTTP exige modo stateless; clientes antigos negociam para baixo. O modo stateless só é ligado se um cliente-alvo exigir essa versão.

Licenças conferidas antes de entrar no `go.mod`: o SDK está sob Apache-2.0, MIT e CC-BY-4.0; `gopkg.in/yaml.v3` v3.0.1 está sob MIT e Apache-2.0. Na Fase 2: `modernc.org/sqlite` BSD-3-Clause, `github.com/google/uuid` BSD-3-Clause, `golang.org/x/crypto` e `golang.org/x/term` BSD-3-Clause.

### D3 — Resource server e authorization server no mesmo processo

O processo publica metadados RFC 9728 e valida bearer em `/mcp`. Também é o authorization server local: Authorization Code + PKCE S256, tokens opacos com hash no SQLite, audience RFC 8707 igual à URL canônica do MCP, refresh com rotação. Não há password grant nem client credentials no lugar desse fluxo.

O login REST em `/api/v1/auth/*` administra usuários com Argon2id. Não registra cliente MCP e não substitui PKCE. O access token emitido por esse authorization server autoriza `/mcp` para cliente CIMD e para cliente DCR. O papel `admin` não libera dados Oracle.

Biblioteca de authorization server (fosite ou equivalente) só entra se implementar CIMD e DCR. Caso contrário, os endpoints desses dois fluxos seguem as RFCs, sem framework OAuth genérico.

### D4 — CIMD e DCR são obrigatórios

Ordem do cliente: pré-registro que ele já possua, depois CIMD, depois DCR. Pré-registro administrativo pode existir e não substitui os outros dois.

- CIMD: `client_id` é URL `https` com path. O servidor busca o JSON, exige `client_id`, `client_name` e `redirect_uris`, confere a URL e o redirect, anuncia `client_id_metadata_document_supported: true`, cacheia conforme os cabeçalhos HTTP e bloqueia SSRF. Não emite `client_id`.
- DCR (RFC 7591): `registration_endpoint` nos metadados. O POST grava em `oauth_clients` e devolve `client_id` deste issuer. Aceita `application_type`, inclusive `native`.

### D5 — Identidade em SQLite

`modernc.org/sqlite`, arquivo `data/identity.db`, migrações SQL embutidas. Sem ORM e sem servidor de banco. Oracle não guarda usuários, tokens nem auditoria.

O access token é opaco, então a tabela `access_tokens` guarda só o hash, a audiência e a expiração. `login_attempts` sustenta o limite de tentativas de login. O cache de CIMD fica na memória do processo. Cliente DCR desta versão é público (`token_endpoint_auth_method=none`).

Na primeira execução, sem administrador ativo, o processo pede usuário, senha e confirmação no terminal. O comando `admin create` faz o mesmo pedido e recusa um segundo administrador. A recuperação controlada não tem atalho neste comando.

### D6 — Oracle somente leitura, driver provisório

O adaptador usa `github.com/sijms/go-ora/v2` v2.9.0 (MIT), Go puro, enquanto P1 não responde. `godror` só entra se a versão, o TLS ou os tipos do ambiente exigirem. Host vazio deixa o processo no ar e o ready só do processo. Com host, a conexão usa service name ou SID, nunca os dois. Usuário e senha vêm das variáveis nomeadas e o ready faz ping. Neste ambiente a conexão por SID na porta 1526, TCP simples, foi aceita pelo `go-ora`. Só SQL cadastrado e parametrizado. Sem tool de SQL livre. A credencial é somente leitura e o servidor não emite INSERT, UPDATE, DELETE nem DDL. Não há sonda de CREATE TABLE. A proteção de escrita é a credencial, não um filtro de texto.

### D7 — Dimensão e fato separados

Cada dimensão do contrato do cliente é um pacote em `internal/dimensao`. Cada fato é um subpacote em `internal/fato` com duas tools: `fato_<assunto>` e `fato_<assunto>_com_dimensoes`. O SQL fica no pacote. Não há motor que monte JOIN na hora, nem catálogo YAML, nem SQL livre.

A resposta da dimensão inclui a hierarquia e as flags do SQL entregue. Filtro omitido é `NULL`. O limite omitido é 500, cortado por `bi.max_rows`. O e-mail do RLS é a coluna `email` do usuário, injetada pelo servidor. `dim_rls_centro_custo` lista só os centros de custo desse e-mail. A tool de natureza se chama `dim_natureza`.

### D8 — SQL no pacote da tool

A query fica no pacote da dimensão ou do fato, não num arquivo de catálogo. Identidade da chamada vem do token. As tools genéricas de listar e executar relatório do catálogo não entram.

### D9 — Pacotes sob demanda

Não criar pacote nem interface vazios sem um lugar já decidido. `internal/database/oracle` entrou na Fase 3. `internal/bi` guarda auditoria, a interface da consulta e o bind dos filtros opcionais. `internal/fato` publica os nove relatórios. `internal/authorization` continua de fora: a regra de linha está no `EXISTS` de cada fato.

### D10 — Segredos e logs

Segredos em variável de ambiente ou gerenciador de segredos. `log/slog` não registra senha, token, SQL nem resultado de consulta.

## Abertas

| ID | Pergunta | Bloqueia |
| --- | --- | --- |
| P1 | Versão do Oracle, sistema operacional e TLS | Escolha do driver (Fase 3) |
| P2 | Como o BI autentica hoje (AD, LDAP, SSO, Oracle ou login próprio) | Qualquer provider além do SQLite. A v1 segue com SQLite |
| P3 | Segurança por linha, filial, empresa ou departamento | Fechado: `DIM_RLS_CENTRO_CUSTO` pelo e-mail do usuário. Ver D7 |
| P4 | Flag de registro ativo e o uso dos níveis da hierarquia como filtro | Fechado: a resposta segue o SQL do cliente, com hierarquia e flags |
| P5 | Qual o primeiro fato depois das duas dimensões | Fechado: os nove fatos do contrato, cada um em duas tools |
| P6 | Rede: interna, VPN ou HTTPS público | Endurecimento de deploy (Fase 5). Desenvolvimento fica em `127.0.0.1` |
| P7 | Clientes MCP da primeira entrega | Aceite do OAuth contra cliente real (Fase 2). Até lá o transporte negocia pelo SDK |
| P8 | Volume de usuários e consultas simultâneas | Não bloqueia a v1 de instância única |
| P9 | Quem faz backup, atualização e administração | Documentação de operação (Fase 5) |
| P10 | Titularidade, reutilização e manutenção do código entregue | Path definitivo do módulo e licença do repositório, antes da entrega |
