# Diretrizes Globais do Projeto (Universal AI - Multi-Agent)
> Compatível com: **Google Antigravity/Gemini**, **Anthropic Claude Code**, **OpenAI ChatGPT/Codex**
>
> 🛡️ **Rastreabilidade de Governança Multi-AI (Cadeia de Custódia):**
> - **Última Atualização:** 2026-11-05 20:54:59
> - **Desenvolvedor Responsável:** herick.ribeiro (Host: BROCKWAY-N01)
> - **Harness / Agente de IA:** Claude Code
> - **Hub Referência:** Framework Multi-AI (Ancorado na ADR-008)
> - **Motivo da Alteração:** Inicialização com Governança Multi-AI

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

## 3. Padrão de Coexistência entre IAs (Google, Anthropic, OpenAI)
- **Não Sobrescrever Instruções Nativas**:
  - `CLAUDE.md` / `.claude/` $\rightarrow$ Exclusivo do Claude Code.
  - `GEMINI.md` / `.antigravity/` / `.gemini/` $\rightarrow$ Exclusivo do Antigravity/Gemini.
  - `AGENTS.md` $\rightarrow$ Este arquivo unificado, consumido por todas as IAs (incluindo OpenAI/Codex).
- **Scripts e Ferramentas (Decisão 003)**:
  - Scripts compartilhados (que qualquer IA ou pessoa roda diretamente) ficam na pasta única `ferramentas/` (ou `tools/`) com o cabeçalho declarando o mantenedor:
    `# mantido_por: [antigravity | claude | codex | time]`
  - Apenas scripts que dependem de hooks, skills exclusivas ou rotinas nativas de uma IA ficam em `scripts/<ia>/`. Alterar o mantenedor atualiza o cabeçalho, nunca a pasta.

---

## 4. Recursos Externos e de Rede (Cross-Drive & Network)
- O código-fonte reside canonicamente em `C:\Projetos\<projeto>\`.
- Acessos a unidades secundárias (`D:\`, `E:\`) ou compartilhamentos de rede UNC (`\\servidor\...`) devem ser parametrizados em arquivos `.env` ou variáveis de ambiente locais, nunca fixados (*hardcoded*) no código.
- NUNCA versione no Git arquivos binários pesados, bancos `*.db` ou dumps contendo dados de clientes.

---

## 5. Livre Circulação de Conhecimento e Zero Bloqueio (Passagem Fluida de Bastão)
- **Zero Restrição de Conhecimento**: Todo código, teste, script e documentação pertencem à organização e são abertos a todos os agentes.
- **Continuidade sem Barreiras**: Se uma IA iniciou um desenvolvimento, qualquer outra IA (Google, Anthropic ou OpenAI) tem total liberdade e autonomia para prosseguir na tarefa, refatorar, corrigir bugs e avançar, sem que uma barra a outra.
- Nenhuma IA deve se recusar a evoluir um código sob a alegação de ter sido feito por outro agente. O objetivo é a entrega contínua com máxima qualidade técnica.

---

## 6. Escopo Multidisciplinar nos 4 Pilares Corporativos
As IAs devem atuar de forma holística em 4 pilares estratégicos:
1. **Engenharia de Software, ERP, Modern Lakehouse & IA/ML:**
   - ERPs: TOTVS Protheus 12.1.2510 (ADVPL/TLPP MVC REST CP1252), TOTVS RM, TOTVS Fluig (BPM, ECM), TOTVS PO-UI.
   - Lakehouse & Dados: BigQuery, Microsoft Fabric (OneLake), Databricks (Delta Lake, Unity Catalog), dbt Core, Dual-BI (Power BI PBIP/TMDL + Superset YAML).
   - Ingestão & Orquestração: Apache Airflow + Astronomer Cosmos, Apache Hop, Airbyte, dlt, n8n.
   - Engenharia de IA & Machine Learning: Model Context Protocol (MCP), RAG Híbrido Corporativo (FTS5 + pgvector), MLOps, Vertex AI, Azure OpenAI.
   - Apps & SaaS: Golang (Hexagonal, Chi, GORM), Next.js (TypeScript, Tailwind, Radix/Shadcn, Recharts), Flutter (Offline-First, Sqflite, Geolocator).
   - Bancos: PostgreSQL (pgvector), SQL Server, Oracle (PL/SQL), MariaDB, SQLite (FTS5).
   - Infra & Observabilidade: Nginx, Certbot SSL, Healthchecks, Docker/Compose.
2. **Marketing (Digital & Trade), CRM, Vendas, SFA & E-commerce:**
   - Marketing Digital & Growth: SEO avançado (Semrush), copywriting persuasivo (AIDA, PAS), funis de conversão e LPs sem menção a concorrentes.
   - Trade Marketing & Comercial: Gestão de PDV, campanhas omnicanal, sell-in/sell-out, políticas comerciais, tabelas de preço, CRM B2B/B2C e regras de comissão.
   - SFA & E-commerce: Automação de força de vendas móvel offline-first (padrão Risa Vendex 360), e-commerce e integrações com Marketplaces (Mercado Livre, Amazon).
3. **Supply Chain, Operações, WMS, Fiscal, Contábil, HCM, GRC, QA & Finanças:**
   - Supply Chain & Estoque: Compras estratégicas, cotações, alçadas, homologação de fornecedores, curva ABC, inventários cíclicos e controle de custos.
   - WMS & Logística: Endereçamento de armazém, recebimento/conferência com coletores, picking, packing, expedição e roteirização.
   - Fiscal & Contabilidade: Apuração de impostos (ICMS, IPI, PIS, COFINS, ISS), Reforma Tributária (IBS/CBS/IS), SPEDs (Fiscal, Contribuições, Reinf, ECD/ECF) e DRE/Balanço.
   - HCM, DP & RH: Folha de pagamento, controle de jornada, eSocial, avaliação de desempenho (9-Box), cargos, salários e recrutamento.
   - Governança de TI, GRC, QA & Lean: ITIL 4, COBIT, FinOps, GRC (ISO 31000, ISO 27001, LGPD, SOX com mentoria ativa ao Herick), QA (Pytest, Go, Playwright, k6, dbt test), Lean IT (VSM, Kaizen, Kanban) e Unit Economics SaaS (LTV, CAC, Payback, Churn, OKRs).
4. **Pesquisa, Inovação & R&D:** Benchmarking competitivo profundo de mercado, prototipagem rápida e análise contínua de novas tecnologias.
👉 Padrão completo: `standards/BUSINESS_INNOVATION_AND_GROWTH.md`.

---

## 6. Sincronização Automática Contínua
- SEMPRE que você alterar qualquer arquivo do Hub Multi-AI (`C:\Projetos\multi-ai\`), execute imediatamente:
  ```bash
  python C:\Projetos\multi-ai\tools\atualizar_documentacao.py
  ```
  para manter toda a documentação visual e técnica 100% atualizada.

---

## 7. Qualidade de Código, Anti-Alucinação e Documentação de Regras de Negócio
- **Código Limpo & Objetivo:** O código deve ser enxuto, modular, bem estruturado, sem funções redundantes, sem código morto e sem abstrações desnecessárias (*over-engineering*).
- **Tolerância Zero a "Lógica Alucinógena":** É expressamente proibido inventar parâmetros, métodos de SDK, tabelas ou fluxos inexistentes. Em caso de dúvida sobre bibliotecas ou ERP, a IA DEVE consultar a documentação local, realizar testes mínimos ou perguntar ao Herick Ribeiro.
- **Documentação Obrigatória de Regras de Negócio:** Sempre que uma regra de negócio for implementada (regras fiscais, comissões, alçadas, validações de fluxo, cálculos de precificação):
  1. *No Código:* Comentário/docstring de intenção explicando o **"PORQUÊ"** da regra (o contexto corporativo ou legal).
  2. *Na Documentação:* Registrada em `knowledge/rules/` ou `.multi-ai/docs/` com seu selo de proveniência ([OFICIAL], [FONTE], [DERIVADO]).
---

## 8. Padrão Mandatório de Modelagem de Dados
- **Nomenclatura no Singular e `snake_case`:** Tabelas e colunas sempre em `snake_case` e no singular (`cliente`, `pedido`, `ordem_servico`). Proibido usar plural ou prefixos de tipos de dados (`str_`, `num_`, `tbl_`).
- **Identificadores com `_id`:** Tanto PK quanto FK seguem o padrão `<tabela>_id` (ex.: PK `cliente_id` na tabela `cliente`, FK `cliente_id` na tabela `pedido`), eliminando ambiguidades em joins e no Lakehouse. Usar UUID (v7/v4) para SaaS/Mobile Offline-First e BigInt Sequencial para alta densidade analítica.
- **Auditoria Universal & Deleção Lógica:** Toda tabela transacional deve conter:
  - `created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `deleted_at TIMESTAMP WITH TIME ZONE NULL` (Soft Delete)
  - Autoria forense e governança: `created_by`, `updated_by`, `deleted_by`, além de `tenant_id` (SaaS) e `version` (concorrência otimista).

---

## 9. Padrão Mandatório de Deploy e Infraestrutura
- **Deploy 100% Dockerizado via Docker Hub:** Builds e compilações ocorrem em pipelines CI/CD ou ambiente de desenvolvimento homologado. As imagens são enviadas ao Docker Hub com versionamento semântico.
- **Zero Compilação no Servidor de Destino:** O servidor de produção/staging nunca compila código nem clona repositórios de código-fonte cru.
- **Manipulação Mínima no Host:** No servidor de destino, manipula-se exclusivamente o `docker-compose.yml` (puxando imagens do Docker Hub) e o arquivo `.env` com os segredos locais do ambiente.
---

## 10. Padrão de Stack Mínima Viável (MVP-First)
- **Simplicidade Inicial Mandatória:** Inicie sempre com a menor stack viável (monólito modular, Go/Next + PostgreSQL ou SQLite). É proibido introduzir Redis, Kafka, Elasticsearch ou microsserviços no MVP sem necessidade comprovada.
- **Esgotamento da Stack Base:** Use tabelas do próprio PostgreSQL para filas/jobs e FTS5/pgvector para busca semântica antes de adicionar novos serviços.
- **Front-end Estático:** Gerar HTML/CSS/JS no build e distribuir como arquivos; nenhum Node em servidores de deploy (Node local só para dev, versão fixada; build de deploy sempre em container). Next.js apenas com `output: 'export'`. Ver `standards/FRONTEND_STATIC_FIRST.md`.
- **Expansão Sob Demanda Autorizada:** Novas tecnologias só entram com gargalo comprovado, registro em ADR (`docs/decisoes/`) e aprovação prévia do Herick Ribeiro.
