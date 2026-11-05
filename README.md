# Risa DataMCP Oracle — Enterprise AI Data Gateway

[![Go Version](https://img.shields.io/badge/go-1.25.0-00ADD8?style=flat&logo=go)](https://golang.org)
[![MCP Protocol](https://img.shields.io/badge/MCP-2026--07--28%20%2F%20Streamable%20HTTP-7B2CBF?style=flat)](https://modelcontextprotocol.io)
[![Oracle Database](https://img.shields.io/badge/Oracle-19c%20%7C%2021c%20%7C%2023ai-F80000?style=flat&logo=oracle)](https://www.oracle.com/database/)
[![Security](https://img.shields.io/badge/Security-Argon2id%20%2B%20PKCE%20%2B%20Zero--SQL--Injection-success)](file:///C:/Projetos/risacode/oracle-mcp/docs/security.md)
[![Status](https://img.shields.io/badge/Build-Passing%20(100%25%20Unit%20Tests)-brightgreen)](file:///C:/Projetos/risacode/oracle-mcp/docs/implementation-checklist.md)

O **Risa DataMCP Oracle** é um servidor MCP (*Model Context Protocol*) corporativo independente, desenvolvido em **Golang** de alta performance, projetado para expor métricas, relatórios analíticos e dados multidimensionais de Business Intelligence (BI) armazenados em bancos **Oracle Database** diretamente para agentes de Inteligência Artificial (Claude Desktop, Cursor IDE, ChatGPT, Claude Code, LangChain, SDKs e agentes autônomos).

Diferente de abordagens ingênuas de IA baseadas em "Text-to-SQL" irrestrito — que geram riscos severos de injeção de SQL, vazamento de credenciais e alucinações de dados —, o **Risa DataMCP Oracle** implementa o conceito de **Catálogo Fechado de Ferramentas Parametrizadas**, onde a IA tem acesso estritamente a **39 ferramentas semânticas** (20 dimensões de negócio e 18 fatos analíticos), com controle de acesso em nível de linha (**RLS - Row Level Security**) injetado dinamicamente a partir da identidade do usuário autenticado.

---

## 🚀 Principais Diferenciais Corporativos

| Pilar | Abordagem Tradicional (Text-to-SQL) | Risa DataMCP Oracle |
| :--- | :--- | :--- |
| **Segurança contra Injeção de SQL** | ❌ Vulnerável a prompt injection e comandos destrutivos | ✅ **Zero SQL Livre**. Todas as queries são compiladas e parametrizadas no binário |
| **Segurança por Linha (RLS)** | ❌ A IA pode "esquecer" cláusulas `WHERE` | ✅ **RLS Determinístico Injetado**. O servidor faz o bind do `:email` do token de sessão |
| **Credenciais Oracle** | ❌ Expostas no prompt ou compartilhadas | ✅ **Usuário Oracle Dedicado 100% Read-Only**, sem privilégio de escrita ou DDL |
| **Isolamento de Identidade** | ❌ Cria tabelas no banco de dados corporativo | ✅ **IdentityStore Local em SQLite** (`identity.db`), sem alterar objetos no Oracle |
| **Padrões de Autenticação** | ❌ Tokens estáticos em `.env` | ✅ **OAuth2 RFC-Compliant** (Authorization Code + PKCE S256, CIMD e DCR RFC 7591) |
| **Performance e Footprint** | ❌ Pesados, exigem JVM ou runtimes de 500MB | ✅ **Binário Go Nativo**, consumo < 25MB de RAM, subida em < 50ms |

---

## 🏛️ Arquitetura Resumida

```text
  [ Cursor / Claude Desktop / Agente de IA ]
                     │
         HTTPS (Streamable HTTP /mcp)
         Bearer Token (OAuth2 / PKCE)
                     ▼
  ┌─────────────────────────────────────────────────────────────┐
  │                   RISA DATAMCP ORACLE                       │
  │                                                             │
  │  ┌──────────────────────┐      ┌─────────────────────────┐  │
  │  │  Identity & Auth     │      │   HTTP Router / Engine  │  │
  │  │  - Argon2id Password │      │   - /api/v1 (Admin API) │  │
  │  │  - PKCE S256 / CIMD  │      │   - /mcp (Streamable)   │  │
  │  │  - SQLite (WAL Mode) │      │   - Health Check & Ping │  │
  │  └──────────┬───────────┘      └────────────┬────────────┘  │
  │             │ Identidade Validada (Email)   │               │
  │             └───────────────┬───────────────┘               │
  │                             ▼                               │
  │              ┌─────────────────────────────┐                │
  │              │    MCP Tool Registry (39)   │                │
  │              │  - 1 Boas-Vindas            │                │
  │              │  - 20 Dimensões de Negócio  │                │
  │              │  - 18 Fatos (Puro + JOINs)  │                │
  │              └──────────────┬──────────────┘                │
  │                             │ Bindings Seguros              │
  │                             ▼                               │
  │              ┌─────────────────────────────┐                │
  │              │   Oracle Read-Only Pool     │                │
  │              │   (go-ora v2 / pure Go)     │                │
  │              └──────────────┬──────────────┘                │
  └─────────────────────────────┼───────────────────────────────┘
                                ▼ TCP:1521 / TCP:1526
                 ┌─────────────────────────────┐
                 │    Oracle Database (BI/DW)  │
                 │    - Views Parametrizadas   │
                 │    - Validação de RLS       │
                 └─────────────────────────────┘
```

---

## 📦 Catálogo de Ferramentas MCP (39 Tools)

### 1. Boas-Vindas & Orientação Semântica
- `boas_vindas`: Instrui o agente de IA sobre como navegar pelo modelo analítico, personalidades recomendadas, logos institucionais e diretrizes de interpretação sem alucinação.

### 2. 20 Dimensões Cadastrais & Estruturais
| Tool | Finalidade de Negócio |
| :--- | :--- |
| `dim_banco_conta_bancaria` | Contas correntes, agências, bancos e saldos operacionais |
| `dim_centro_custo` | Estrutura hierárquica e centros de custo ativos |
| `dim_cliente` | Base de clientes, CNPJ/CPF, segmento e região |
| `dim_colaborador` | Cadastro de colaboradores e vínculos ativos |
| `dim_conta_orcamentaria` | Plano de contas orçamentárias e classificação contábil |
| `dim_conta_orcm_natureza`| Matriz de correlação entre orçamento e natureza financeira |
| `dim_empresa_filial` | Unidades de negócio, empresas do grupo e filiais |
| `dim_evento` | Rubricas e eventos de folha de pagamento |
| `dim_fornecedor` | Fornecedores homologados, dados fiscais e categorias |
| `dim_funcao_colaborador` | Cargos e funções corporativas |
| `dim_moeda` | Moedas operacionais e taxas de conversão |
| `dim_natureza` | Naturezas financeiras (receitas, custos, despesas) |
| `dim_rls_centro_custo` | Centros de custo autorizados exclusivamente para o e-mail do usuário |
| `dim_secao` | Seções e divisões organizacionais de colaboradores |
| `dim_situacao_colaborador`| Situações funcionais (ativo, afastado, férias, desligado) |
| `dim_tempo` | Calendário analítico, dias úteis, bimestres, trimestres |
| `dim_tipo_mov_bancaria` | Tipos de movimentação de extrato bancário |
| `dim_tipo_processo` | Tipos de processos judiciais e administrativos |
| `dim_tipo_recebimento_colab` | Modalidades de pagamento/recebimento de pessoal |
| `dim_tipo_titulo` | Tipos de títulos a pagar e a receber (duplicatas, boletos, etc.) |

### 3. 18 Fatos Analíticos (9 Grãos × 2 Variantes)
Para cada grão de dados, o servidor fornece duas ferramentas: a versão **pura** (códigos e valores brutos para agregações ultrarrápidas) e a versão **com_dimensoes** (enriquecida com descrições, nomes e classificações hierárquicas completas).

1. `fato_baixa_titulo_imposto` / `fato_baixa_titulo_imposto_com_dimensoes`
2. `fato_conta_pagar_rateado` / `fato_conta_pagar_rateado_com_dimensoes`
3. `fato_conta_receber_rateado` / `fato_conta_receber_rateado_com_dimensoes`
4. `fato_confissao_divida` / `fato_confissao_divida_com_dimensoes`
5. `fato_folha_pagamento` / `fato_folha_pagamento_com_dimensoes`
6. `fato_historico_colaborador` / `fato_historico_colaborador_com_dimensoes`
7. `fato_mov_bancaria_rateada` / `fato_mov_bancaria_rateada_com_dimensoes`
8. `fato_orcamento_total` / `fato_orcamento_total_com_dimensoes`
9. `fato_processo` / `fato_processo_com_dimensoes`

---

## 🛠️ Instalação e Execução

### Pré-requisitos
- **Go 1.25.0+** (ou Go 1.23+ com toolchain automático)
- Acesso de rede à porta do Oracle (1521 ou 1526)
- Sistema Operacional: Windows, Linux ou macOS

### 1. Configuração de Variáveis de Ambiente
Copie o modelo de variáveis de ambiente:
```bash
cp .env.example .env
```

Edite o `.env` ou configure no seu gerenciador de segredos:
```env
ORACLE_MCP_ADDRESS=127.0.0.1:8080
ORACLE_MCP_PUBLIC_BASE_URL=https://mcp.suaempresa.com
ORACLE_MCP_SQLITE_PATH=./data/identity.db

# Credenciais Oracle (Read-Only)
ORACLE_MCP_DB_HOST=10.0.0.50
ORACLE_MCP_DB_PORT=1521
ORACLE_MCP_DB_SERVICE=BI_PROD
ORACLE_MCP_DB_USER=USR_BI_MCP_RO
ORACLE_MCP_DB_PASSWORD=SuaSenhaForteAqui
```

### 2. Bootstrap do Administrador Inicial
Na primeira execução, o servidor detecta a ausência de administradores e solicita a criação interativa no console de forma segura (sem salvar senhas em logs ou histórico):
```powershell
go run ./cmd/server admin create
```
*Ou execute o binário compilado:*
```powershell
.\oracle-mcp.exe admin create
```

### 3. Execução do Servidor
```powershell
go run ./cmd/server
```
O servidor inicializará escutando em `http://127.0.0.1:8080`:
- **Health Check Liveness:** `GET /api/v1/health/live`
- **Health Check Readiness (com ping Oracle):** `GET /api/v1/health/ready`
- **Endpoint MCP Streamable HTTP:** `POST /mcp`
- **OAuth2 Metadata:** `GET /.well-known/oauth-authorization-server`
- **Protected Resource Metadata:** `GET /.well-known/oauth-protected-resource`

---

## 🧪 Suíte de Testes Automatizados

O projeto conta com cobertura de testes unitários abrangente para todas as camadas de API, autenticação, drivers e ferramentas MCP:

```powershell
# Execução completa da suíte de testes
go test -v ./...
```

---

## 📚 Documentação Corporativa Completa

Acesse os artefatos de engenharia e planejamento:
- [📘 Blueprint de Arquitetura Hexagonal](file:///C:/Projetos/risacode/oracle-mcp/ARCHITECTURE_BLUEPRINT.md)
- [📊 Requisitos, Status e Benchmark Competitivo](file:///C:/Projetos/risacode/oracle-mcp/REQUISITOS_E_STATUS_CONSOLIDADO.md)
- [🖥️ Simulador & Protótipo Navegável](file:///C:/Projetos/risacode/oracle-mcp/prototipo-navegavel.html)
- [📖 Manual de Utilização Interativo](file:///C:/Projetos/risacode/oracle-mcp/manual.html)
- [🌐 Landing Page Comercial](file:///C:/Projetos/risacode/oracle-mcp/lp-oracle-mcp.html)
- [🏛️ Portal Unificado de Documentação](file:///C:/Projetos/risacode/oracle-mcp/docs/index.html)

---

## 🔒 Licença e Governança

Desenvolvido segundo as diretrizes de governança e segurança de software corporativo da **Risa Code**. Código auditável, dependências abertas homologadas (MIT / BSD-3 / Apache-2.0) e zero dependências de repositórios proprietários externos.
