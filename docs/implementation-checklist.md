# Checklist de implementação

Cada fase termina com `go test ./...`, uma nota do que ficou pendente e uma parada para autorizar a próxima. Commit só quando for pedido.

Estado atual: Fase 4 em pacotes de dimensão, sem fato e sem dados reais de BI. A conexão Oracle de leitura respondeu ao ping e ao `SELECT` em `DUAL` com `go-ora`, SID e TCP simples. O aceite dos lookups depende de P3 e P4. A versão do Oracle (P1) ainda não foi registrada.

## Fase 0 — Descoberta

- [ ] Registrar respostas de P1 a P10 em `docs/decisions.md` quando a empresa responder.
- [ ] Conferir licença das bibliotecas antes de cada entrada no `go.mod`.
- [ ] Confirmar que nenhum arquivo de produto privado entrou neste repositório.

A Fase 1 não espera a Fase 0. P1 bloqueia o aceite da Fase 3. P3 e P4 bloqueiam dados reais na Fase 4. P7 bloqueia o aceite do OAuth contra cliente real na Fase 2.

## Fase 1 — Base

Autorização para começar: concedida.

- [x] `go.mod` com path `oracle-mcp` e `go 1.25.0` fixado.
- [x] `.gitignore` cobrindo `data/` e `.env`.
- [x] Configuração YAML opcional mais variáveis de ambiente, validada na subida.
- [x] Logs com `log/slog`, sem segredos.
- [x] `GET /api/v1/health/live` e `GET /api/v1/health/ready` (ready da Fase 1 responde o processo).
- [x] `POST /mcp` com Streamable HTTP, sessão stateful.
- [x] Teste de health e de handshake (`initialize` / `tools/list`) com o cliente do SDK.
- [x] Aceite: o processo sobe em localhost e um cliente MCP de teste completa o handshake.

Fora desta fase: ferramenta de negócio, SQLite, OAuth, Oracle.

## Fase 2 — Identidade e auth

Autorização para começar: concedida.

- [x] Migrações: `users`, `roles`, `permissions`, `user_roles`, `role_permissions`, `refresh_tokens`, `oauth_clients`, `oauth_authorization_codes`, `audit_log`. Também `access_tokens` e `login_attempts`, necessários para token opaco e limite de login.
- [x] `oracle-mcp admin create` com senha no terminal, sem argumento e sem endpoint de bootstrap. Recusa um segundo administrador. A primeira execução do servidor, sem administrador, faz o mesmo pedido antes de escutar.
- [x] API `/api/v1` da spec, soft delete e troca obrigatória de senha temporária.
- [x] Argon2id, access token curto, refresh com rotação. Só o hash é gravado.
- [x] Authorization Code + PKCE S256.
- [x] Metadados RFC 8414 com `client_id_metadata_document_supported: true` e `registration_endpoint`.
- [x] Metadados RFC 9728 e bearer obrigatório em `/mcp`.
- [x] CIMD: URL `https` com path, `client_id` igual à URL, redirect na lista, cache, timeout, tamanho máximo, SSRF recusado.
- [x] DCR: POST persiste em `oauth_clients`, aceita `application_type: native`, `client_id` vale só neste issuer.
- [x] Papéis `admin`, `analyst`, `viewer`. Permissões `users:read` e `users:manage`.
- [x] Testes: hash, migração, CRUD, login, refresh, logout, reuse de refresh revogado, `/mcp` sem bearer responde 401, PKCE via CIMD, PKCE após DCR.
- [ ] Aceite contra cliente real: depende de P7.

## Fase 3 — Oracle

Autorização para começar: concedida.

- [x] Spike do driver no ambiente autorizado. `go-ora` v2.9.0 conectou por SID, porta 1526, TCP simples, sem wallet. `godror` não foi necessário nesta sessão.
- [x] Pool `database/sql`, timeout, máximo de linhas, máximo de bytes.
- [x] Normalização de `NUMBER`, `DATE`, `TIMESTAMP`, `CLOB` e `NULL`.
- [x] `health/ready` inclui ping do Oracle sem detalhe interno. Sem host, o ready continua só do processo.
- [x] Teste de integração opcional, atrás da variável `ORACLE_MCP_ORACLE_INTEGRATION=1`. Só ping e `SELECT` em `DUAL`. Sem sonda de escrita. Passou neste ambiente.
- [ ] Aceite de versão e sistema operacional: P1 ainda não foi respondido. A conexão em si funcionou.

## Fase 4 — Dimensões e fatos

Autorização para começar: concedida. Sem senha Oracle e sem view real.

- [x] `internal/dimensao/centrocusto` e `internal/dimensao/natureza`, com a query e a descrição no pacote.
- [x] `internal/fato` sem tool. O primeiro relatório entra como subpacote e chama `Contains` da dimensão que citar.
- [x] `dim_centro_custo` e `dim_natureza_financeira` com busca opcional e a folha: id, código, nome e rótulo.
- [x] Identidade vem do token. Escopo aplicado no servidor.
- [x] Sem Oracle, o lookup responde consulta indisponível, sem SQL na resposta.
- [x] `dim_centro_custo` e `dim_natureza_financeira` devolvem no máximo 200 linhas da folha.
- [x] Auditoria com ferramenta, ator, duração e outcome, sem linhas retornadas.
- [ ] Aceite com dados reais: depende de P3 e P4, e da comparação com o BI oficial.

## Fase 5 — Entrega

- [ ] Limites de payload e concorrência.
- [ ] Backup consistente do SQLite, incluindo WAL.
- [ ] Procedimento de rotação de segredo e recuperação administrativa.
- [ ] `docs/security.md`, `docs/api.md`, `docs/deployment.md`.
- [ ] README de build e execução.
- [ ] Lista de licenças de terceiros.
- [ ] Aceite: a empresa compila, configura, sobe, administra usuários e consulta os lookups sem dependência privada.
