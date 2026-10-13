# DIM_BANCO_CONTA_BANCARIA

**Tipo:** Dimensão

## Descrição

Dimensão de bancos e contas bancárias, com agência, número, tipo de conta e a empresa/filial à qual cada conta pertence.

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_empresa` | inteiro | Não | `ID Empresa` | igual | Empresa dona da conta. |
| `id_filial` | inteiro | Não | `ID Filial` | igual | Filial dona da conta. |
| `id_banco` | inteiro | Não | `ID Banco` | igual | ID do banco. |
| `codigo_banco` | texto | Não | `Código Banco` | igual | Código do banco (ex.: 001, 237, 341). |
| `nome_banco` | texto | Não | `Nome Banco` | contém | Busca parcial pelo nome do banco. |
| `id_conta_bancaria` | inteiro | Não | `ID Conta Bancária` | igual | ID da conta bancária. |
| `nome_conta` | texto | Não | `Nome Conta Bancária` | contém | Busca parcial pelo nome da conta. |
| `tipo_conta_bancaria` | texto | Não | `Tipo Conta Bancária` | igual | Tipo da conta. |
| `registro_ativo` | texto | Não | `Registro Ativo Conta Bancária` | igual | Flag de conta ativa. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Banco",
    "Código Banco",
    "Nome Banco",
    "Banco",
    "ID Conta Bancária",
    "Código Agência",
    "Número Conta Bancária",
    "Nome Conta Bancária",
    "Nome Reduzido Conta Bancária",
    "Tipo Conta Bancária",
    "Conta Bancária",
    "Banco, Agência e Conta Bancária",
    "ID Empresa",
    "ID Filial",
    "Código Filial",
    "Nome Filial",
    "Filial",
    "Registro Ativo Banco",
    "Reg. Bloqueado Conta Bancária",
    "Registro Ativo Conta Bancária"
FROM DIM_BANCO_CONTA_BANCARIA
WHERE (:id_empresa IS NULL OR "ID Empresa" = :id_empresa)
  AND (:id_filial IS NULL OR "ID Filial" = :id_filial)
  AND (:id_banco IS NULL OR "ID Banco" = :id_banco)
  AND (:codigo_banco IS NULL OR "Código Banco" = :codigo_banco)
  AND (:nome_banco IS NULL OR UPPER("Nome Banco") LIKE '%' || UPPER(:nome_banco) || '%')
  AND (:id_conta_bancaria IS NULL OR "ID Conta Bancária" = :id_conta_bancaria)
  AND (:nome_conta IS NULL OR UPPER("Nome Conta Bancária") LIKE '%' || UPPER(:nome_conta) || '%')
  AND (:tipo_conta_bancaria IS NULL OR "Tipo Conta Bancária" = :tipo_conta_bancaria)
  AND (:registro_ativo IS NULL OR "Registro Ativo Conta Bancária" = :registro_ativo)
ORDER BY "Nome Banco", "Nome Conta Bancária"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Permitir que a IA identifique bancos e contas bancárias (por nome, código ou filial) e obtenha os IDs usados para filtrar e rotular as movimentações bancárias e baixas de títulos.
