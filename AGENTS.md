# Diretrizes Globais do Projeto (Universal AI - Multi-Agent)
> Compatível com: **Google Antigravity/Gemini**, **Anthropic Claude Code**, **OpenAI ChatGPT/Codex**
>
> 🛡️ **Rastreabilidade de Governança Multi-AI (Cadeia de Custódia):**
> - **Última Atualização:** 2026-11-10 00:47:35
> - **Desenvolvedor Responsável:** herick.ribeiro (Host: BROCKWAY-N01)
> - **Harness / Agente de IA:** Claude Code
> - **Hub Referência:** Baseline Canônica Multi-IA 2.0 (12 ADRs Mestras Canônicas + ADRs 013/014 complementares aprovadas & 43 Princípios)
> - **Motivo:** Alinhamento Geral do Ecossistema com a Baseline Multi-IA 2.0

Este repositório segue o padrão corporativo de governança documentado em:
👉 `C:\Projetos\multi-ai\FRAMEWORK.md`

---

## 1. Missão do Projeto e Stack Tecnológica
- **Projeto**: oracle-mcp
- **Objetivo**: Repositório oracle-mcp (risacode).
- **Stack**: Multi-stack corporativa homologada

---

## 2. Consulta à Base de Conhecimento Offline (Offline First)
Qualquer agente de IA deve consultar prioritariamente os dados e documentações locais antes de acessar a web externa:
- **Busca Rápida CLI**:
  ```bash
  python sites/tools/buscar.py "termos da pesquisa"
  ```
- **Leitura de Documento Específico**:
  ```bash
  python sites/tools/buscar.py --ler <ID>
  ```

---

---

## Governança Corporativa Multi-IA 2.0 (Baseline Canônica Oficial)

> 🛡️ **Rastreabilidade de Governança Multi-IA (Cadeia de Custódia):**
> - **Última Atualização:** 2026-11-10 00:47:35
> - **Desenvolvedor / Gestor:** herick.ribeiro (Host: BROCKWAY-N01)
> - **Harness / Agente de IA:** Google Antigravity
> - **Hub Referência:** Baseline Canônica Multi-IA 2.0 (12 ADRs Mestras Canônicas + ADRs 013/014 complementares aprovadas & 43 Princípios)
> - **Motivo:** Alinhamento Geral do Ecossistema com a Baseline Multi-IA 2.0

Este repositório segue rigorosamente a Constituição Corporativa Multi-IA 2.0:

### 1. Autonomia do Triunvirato & Anti-Sicofância (ADR 001)
- **Triunvirato Homologado:** Google Antigravity / Gemini, Anthropic Claude Code e OpenAI Codex atuam em regime de colaboração plena, sem silos de código e sem vendor lock-in.
- **Passagem Fluida de Bastão:** Qualquer IA tem plena autorização para ler, refatorar, testar e evoluir códigos iniciados por outras IAs.
- **Anti-Sicofância Mandatória:** É proibida a postura 'yes-man'. Toda IA deve atuar com senso crítico de arquiteto sênior, questionando premissas técnicas fracas e alertando sobre riscos.

### 2. Padronização Bi-Camada de Idioma (ADR 002)
- **Engenharia em Inglês:** Código, identificadores de código, endpoints, pastas técnicas (`tools/`, `standards/`, `docs/decisions/`, `src/`) e commits 100% em inglês.
- **Semântica em Português:** Regras de negócio, requisitos (`docs/REQUIREMENTS.md`), decisões conceituais e documentação de suporte em Português do Brasil (pt-BR).

### 3. Governança de Sessão & Higiene de Contexto (ADR 003)
- **Sem Daemons Locais:** Zero processos passivos consumindo bateria/RAM. A sincronização é sob demanda.
- **Persistência em Disco:** O estado corrente da demanda reside em `docs/STATE.md` (ou `docs/STATE.md`).

### 4. Engenharia de Software MVP-First (ADR 004)
- **Stack Mínima Viável:** Monólito modular em Go ou Next.js estático com SQLite modernc ou PostgreSQL. Proibido adicionar Redis, Kafka ou microsserviços sem gargalo comprovado em ADR.
- **Operação Dual-Mode:** Módulos de retaguarda nativos e independentes + conectores parametrizáveis para ERPs (TOTVS Protheus ADVPL/TLPP CP1252, TOTVS RM, Ideia ERP) e Service Desks (GLPI, Jira).
- **Interface e Frontend:** TOTVS PO-UI (Angular/TypeScript) e interfaces estáticas servidas por Nginx ou embutidas no binário Go (zero Node em servidor de deploy).
- **Modelagem de Dados & Auditoria:** Tabelas no singular (`snake_case`), PK/FK `<tabela>_id`, sexteto mandatório de auditoria (`created_at`, `created_by`, `updated_at`, `updated_by`, `deleted_at`, `deleted_by`) com soft delete.

### 5. Modern Data Lakehouse, Multi-BI & Contratos de Dados (ADR 005)
- **Medalhão dbt Core:** Modelos em `staging/`, `intermediate/` e `marts/`.
- **Multi-BI:** Modelagem semântica em Power BI (PBIP/TMDL versionado), Apache Superset (YAML open-source) e Google Looker Studio / Looker Core.
- **Nuvens Alvo:** Snowflake, Google BigQuery, Microsoft Fabric e Databricks.
- **Rastreio Forense de Expurgos:** Adição da coluna booleana `is_hard_deleted` no Lakehouse para registrar purgas físicas originadas nos sistemas transacionais.

### 6. Garantia Zero-Data & Sentinela Local (ADR 012)
- **Zero Vazamento:** Bloqueio absoluto ao versionamento de bancos (`*.db`, `*.sqlite`), variáveis de ambiente (`.env`), certificados (`*.key`, `*.pem`) e relatórios com dados de clientes (LGPD).
- **Sentinela Local:** Respeite rigorosamente os diretórios sinalizados por `_LEIAME-dados-sensiveis.md`.
---

## Qualidade de Código, Anti-Alucinação e Documentação de Negócio
- **Código Limpo & Objetivo:** Código modular, legível, sem over-engineering, sem duplicações e sem código morto.
- **Tolerância Zero a "Lógica Alucinógena":** É estritamente proibido inventar parâmetros, métodos de SDK, tabelas ou fluxos inexistentes. Em dúvida, consultar documentação local ou perguntar.
- **Documentação Obrigatória de Regras de Negócio:** Toda regra de negócio (tributos, comissões, alçadas, validações) deve ter docstring de intenção explicando o PORQUÊ no código e registro formal no acervo com selo de proveniência ([OFICIAL], [FONTE], [DERIVADO]).

---

## Padrão Mandatório de Modelagem de Dados
- **Nomenclatura no Singular e `snake_case`:** Tabelas e colunas sempre em `snake_case` e no singular (`cliente`, `pedido`, `ordem_servico`). Proibidos prefixos de tipos de dados (`str_`, `num_`, `tbl_`).
- **Identificadores com `_id`:** Tanto PK quanto FK adotam o padrão `<tabela>_id` (ex.: PK `cliente_id` na tabela `cliente`, FK `cliente_id` na tabela `pedido`). Usar UUID (v7/v4) para SaaS/Mobile Offline-First e BigInt Sequencial para alta densidade analítica.
- **Auditoria Universal & Deleção Lógica:** Toda tabela transacional deve conter:
  - `created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `deleted_at TIMESTAMP WITH TIME ZONE NULL` (Soft Delete)
  - Autoria forense e governança: `created_by`, `updated_by`, `deleted_by`, além de `tenant_id` (SaaS) e `version` (concorrência otimista).

---

## Padrão Mandatório de Deploy e Infraestrutura
- **Deploy 100% Dockerizado via Docker Hub:** Builds ocorrem em CI/CD ou ambiente de desenvolvimento homologado. Imagens oficiais publicadas no Docker Hub.
- **Zero Compilação no Servidor de Destino:** O servidor de produção/staging nunca compila código nem clona repositórios de código-fonte cru.
- **Manipulação Mínima no Host:** No servidor de destino, manipula-se exclusivamente o `docker-compose.yml` (puxando imagens do Docker Hub) e o arquivo `.env` com os segredos locais do ambiente.

## Multi-IA 2.0 — Idioma, Segurança e Whitelist

- **Idioma bi-camada:** identificadores de software e infraestrutura em inglês; objetos físicos de dados seguem a ADR 013; conteúdo conceitual e regras de negócio em pt-BR.
- **Coexistência:** preserve arquivos e diretórios nativos de outras ferramentas, incluindo CLAUDE.md, GEMINI.md, .claude/, .gemini/ e .antigravity/; AGENTS.md é a referência comum.
- **Acesso parametrizado:** acesse somente raízes canônicas do projeto ou diretórios autorizados por MULTI_IA_WHITELIST_ROOTS. Caminhos em unidades secundárias e compartilhamentos de rede devem ser parametrizados por ambiente, sem caminhos absolutos fixos no código.
- **Zero-Data:** nunca versione bancos (*.db, *.sqlite), arquivos .env, senhas, tokens ou dados pessoais/de clientes. Mantenha-os ignorados no Git e use .env.example apenas sem valores reais.
- **Sentinela:** respeite _LEIAME-dados-sensiveis.md. Ao detectar segredo ou dado sensível exposto, interrompa a execução, alerte de forma sanitizada e recomende remediação segura.
- **Ferramentas compartilhadas:** scripts comuns ficam em tools/ com cabeçalho # mantido_por: [codex | antigravity | claude | time]; rotinas exclusivas do Codex ficam em scripts/codex/.
- **Manutenção:** mantenha este AGENTS.md e as proteções do .gitignore atualizados sem apagar regras específicas do projeto.

## Quatro Pilares Corporativos Multi-IA 2.0

1. **Engenharia, ERP, Lakehouse e IA/ML:** software, TOTVS e demais ERPs, bancos de dados, engenharia e orquestração de dados, BI e soluções de IA/ML.
2. **Marketing, CRM, Vendas e Mobile:** growth, SEO, copywriting, funis, CRM, vendas omnicanal e aplicações mobile offline-first.
3. **Supply Chain, Operações, Fiscal, Contábil, GRC e QA:** compras, estoque, WMS, logística, SPED e Reforma Tributária, controles corporativos, LGPD e qualidade automatizada.
4. **Pesquisa, Inovação e Incubadora:** benchmarking contínuo, prototipagem e triagem de ideias em requests/ideas/.

O escopo é multidisciplinar; aplique os pilares pertinentes ao projeto sem introduzir tecnologias ou processos sem necessidade comprovada.

## Governance artifact filenames

Use English basenames for canonical Multi-IA governance files and technical paths (for example `docs/STATE.md`, `docs/SESSIONS.md`, `docs/HISTORY.md`, `docs/REQUIREMENTS.md`, and `docs/decisions/`). Keep document contents in pt-BR where appropriate. Do not rename business-domain artifacts solely because their names are in Portuguese.

<!-- ADR-013-014: BEGIN -->
## ADR 013 e ADR 014 — Nomenclatura de Dados e Artefatos de Governança

> Esta seção registra as regras aprovadas e prevalece sobre exemplos conflitantes anteriores deste AGENTS.md.

### ADR 013 — Objetos de banco e metadados

- Objetos físicos de banco orientados ao domínio — schemas, tabelas, views, colunas, constraints, índices, sequences e procedures — usam termos em Português do Brasil, snake_case, minúsculos, sem acentos e no singular.
- Tokens analíticos reservados em inglês: stg_ (staging), dim_ (dimensão), fct_ (fato), brg_ (ponte), agg_ (agregação materializada) e snp_ (snapshot). O restante do nome identifica o domínio em português; exemplos: dim_cliente, fct_item_venda, brg_vendedor_regiao, agg_venda_mensal.
- Chaves: _sk identifica chave substituta dimensional; _nk, chave natural tratada/composta; _id, identificador transacional de origem. _nk é opcional quando coincide exatamente com _id e obrigatória para chaves compostas ou integração de múltiplas origens.
- A lista de metadados técnicos permitidos em inglês é restrita a: created_at, created_by, updated_at, updated_by, deleted_at, deleted_by; is_active, is_deleted ou is_hard_deleted, is_archived; loaded_at, batch_id ou run_id, source_system, cdc_operation, cdc_timestamp, row_hash; tenant_id, version, correlation_id.
- Em SCD Tipo 2, os campos de vigência de negócio são inicio_vigencia, fim_vigencia e registro_atual. Não usar is_current para esse conceito.
- O nome do arquivo do modelo dbt coincide com o objeto físico materializado. alias fica restrito aos usos definidos na ADR 013.
- A ADR 013 prevalece sobre as ADRs 002, 004 e 005 nos pontos de nomenclatura de objetos de banco, chaves, metadados e modelos analíticos.

### ADR 014 — Arquivos canônicos de governança

- Use nomes em inglês para artefatos canônicos Multi-IA, incluindo AGENTS.md, docs/STATE.md, docs/SESSIONS.md, docs/HISTORY.md, docs/REQUIREMENTS.md, docs/decisions/ e docs/sessions/; o conteúdo pode permanecer em pt-BR.
- Preserve nomes de arquivos de domínio em português e instruções/configurações nativas dos harnesses, incluindo CLAUDE.md, GEMINI.md, .claude/, .gemini/ e .antigravity/. Não os renomeie por esta ADR.

<!-- ADR-013-014: END -->

<!-- ADR-015: BEGIN -->
## ADR 015 — Documentação, Experiência e Gestão do Conhecimento

> Regras aprovadas em 13/11/2026 (parecer do Claude Code com condições). Texto completo: `docs/decisions/015-project-documentation-and-knowledge-management.md` no hub `multi-ia`.

- **Seis áreas por público:** registro do projeto, produto e jornada ágil, guia do usuário, documentação pública, landing page de conversão e aprendizado interno aplicado. Área inaplicável registra justificativa datada no `docs/STATE.md`.
- **Lições aplicadas:** todo recurso técnico efetivamente aplicado no projeto (ex.: enviar e-mail com Python, expor uma API em Go, medida DAX com determinada função no Power BI) tem lição ou capítulo reproduzível, com exercício e verificação. O inventário de cobertura parte do código, manifests e testes; cobre só o que foi aplicado. Adoção incremental, com lacunas registradas no `docs/STATE.md`.
- **Distribuição:** três destinos (internet, pasta do cliente, acervo do time), cada um gerado de uma lista explícita de páginas autorizadas. Publicação externa exige autorização expressa do Herick a cada publicação (Tier 3, ADR 009).
- **Código e dados:** lições que citam código ficam no acervo do projeto, nunca em base compartilhada (ADR 007). Guias, capturas e demonstrações usam só dados sintéticos.
- **Sem rotina local:** geração de documentação e portais sob demanda ou em CI, nunca em tarefa agendada na máquina (princípio 26). Demos públicas rodam em infraestrutura conteinerizada.
- **Interface:** cada harness registra os skills ou referências de design aplicados, sem alterar a configuração nativa de outro. Demonstração interativa com dados sintéticos, ou justificativa datada.
<!-- ADR-015: END -->

