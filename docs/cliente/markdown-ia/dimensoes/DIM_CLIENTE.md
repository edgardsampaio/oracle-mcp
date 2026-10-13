# DIM_CLIENTE

**Tipo:** Dimensão

## Descrição

Dimensão de clientes com dados cadastrais (razão social, nome fantasia, CNPJ/CPF, inscrições, endereço, município, UF e região).

## Parâmetros

Todos os filtros são opcionais, salvo indicação: quando o parâmetro vem nulo, o filtro é ignorado. Buscas do tipo "contém" não diferenciam maiúsculas de minúsculas.

| Parâmetro | Tipo | Obrigatório | Coluna | Filtro | Descrição |
|---|---|---|---|---|---|
| `id_cliente` | inteiro | Não | `ID Cliente` | igual | ID do cliente. |
| `codigo_cliente` | texto | Não | `Código Cliente` | igual | Código do cliente. |
| `loja_cliente` | texto | Não | `Loja Cliente` | igual | Loja do cliente. |
| `razao_social` | texto | Não | `Razão Social` | contém | Busca parcial pela razão social. |
| `nome_fantasia` | texto | Não | `Nome Fantasia` | contém | Busca parcial pelo nome fantasia. |
| `cnpj` | texto | Não | `CNPJ` | igual | CNPJ/CPF completo. |
| `raiz_cnpj` | texto | Não | `Raíz CNPJ` | igual | Raiz do CNPJ (8 primeiros dígitos) — agrupa matriz e filiais. |
| `sigla_uf` | texto | Não | `Sigla UF` | igual | UF (ex.: RJ, SP). |
| `nome_municipio` | texto | Não | `Nome Município` | contém | Busca parcial pelo município. |
| `regiao` | texto | Não | `Região` | igual | Região geográfica. |
| `codigo_tipo_pessoa` | texto | Não | `Código Tipo Pessoa` | igual | Tipo de pessoa (física/jurídica). |
| `registro_ativo` | texto | Não | `Registro Ativo Cliente` | igual | Flag de cliente ativo. |
| `limite` | inteiro | Não | — | — | Máximo de linhas retornadas (padrão sugerido: 500). |

## Query

```sql
SELECT
    "ID Cliente",
    "Código Cliente",
    "Loja Cliente",
    "Razão Social",
    "Nome Fantasia",
    "Cliente Razão Social",
    "Cliente Nome Fantasia",
    "Cliente CNPJ Razão Social",
    "Cliente CNPJ Nome Fantasia",
    "Raíz CNPJ",
    "CNPJ",
    "Inscrição Estadual",
    "Inscrição Municipal",
    "Endereço",
    "Bairro",
    "CEP",
    "Código Tipo Pessoa",
    "Descrição Tipo Pessoa",
    "Tipo Pessoa",
    "ID Município",
    "Código Município",
    "Nome Município",
    "Município",
    "ID UF",
    "Sigla UF",
    "Nome UF",
    "UF",
    "Região",
    "Sigla País",
    "Registro Bloqueado Cliente",
    "Registro Ativo Cliente",
    "Registro Bloqueado Pessoa",
    "Registro Ativo Pessoa",
    "Registro Ativo Município",
    "Registro Ativo UF"
FROM DIM_CLIENTE
WHERE (:id_cliente IS NULL OR "ID Cliente" = :id_cliente)
  AND (:codigo_cliente IS NULL OR "Código Cliente" = :codigo_cliente)
  AND (:loja_cliente IS NULL OR "Loja Cliente" = :loja_cliente)
  AND (:razao_social IS NULL OR UPPER("Razão Social") LIKE '%' || UPPER(:razao_social) || '%')
  AND (:nome_fantasia IS NULL OR UPPER("Nome Fantasia") LIKE '%' || UPPER(:nome_fantasia) || '%')
  AND (:cnpj IS NULL OR "CNPJ" = :cnpj)
  AND (:raiz_cnpj IS NULL OR "Raíz CNPJ" = :raiz_cnpj)
  AND (:sigla_uf IS NULL OR "Sigla UF" = :sigla_uf)
  AND (:nome_municipio IS NULL OR UPPER("Nome Município") LIKE '%' || UPPER(:nome_municipio) || '%')
  AND (:regiao IS NULL OR "Região" = :regiao)
  AND (:codigo_tipo_pessoa IS NULL OR "Código Tipo Pessoa" = :codigo_tipo_pessoa)
  AND (:registro_ativo IS NULL OR "Registro Ativo Cliente" = :registro_ativo)
ORDER BY "Razão Social"
FETCH FIRST :limite ROWS ONLY;
```

## Objetivo

Localizar clientes por nome, documento ou localização e obter o ID usado em Contas a Receber e movimentações bancárias.
