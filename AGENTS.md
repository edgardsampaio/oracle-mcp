# Diretrizes Globais do Projeto (Universal AI - Multi-Agent)
> Compatível com: **Google Antigravity/Gemini**, **Anthropic Claude Code**, **OpenAI ChatGPT/Codex**
>
> 🛡️ **Rastreabilidade de Governança Multi-AI (Cadeia de Custódia):**
> - **Última Atualização:** 2026-11-10 00:47:35
> - **Desenvolvedor Responsável:** herick.ribeiro (Host: BROCKWAY-N01)
> - **Harness / Agente de IA:** Claude Code
> - **Hub Referência:** Baseline Canônica Multi-IA 2.0 (12 ADRs Mestras Canônicas & 43 Princípios)
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
> - **Hub Referência:** Baseline Canônica Multi-IA 2.0 (12 ADRs Mestras Canônicas & 43 Princípios)
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

- **Idioma bi-camada:** identificadores de código, scripts, pastas técnicas, tokens analíticos (`dim_`, `fct_`, `_sk`, `_nk`) e metadados técnicos (auditoria, CDC, concorrência) em inglês; objetos de banco do domínio de negócio (schemas, tabelas, colunas, constraints) em pt-BR, `snake_case`, minúsculo, sem acentos, conforme a ADR 013; justificativas, regras de negócio e documentação conceitual ficam em Português do Brasil. Para objetos de banco, a ADR 013 prevalece sobre as ADRs 002, 004 e 005.
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
