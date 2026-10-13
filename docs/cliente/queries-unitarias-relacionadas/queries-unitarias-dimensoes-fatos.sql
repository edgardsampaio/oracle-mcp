--DIMENSÕES

SELECT "ID Banco", "Código Banco", "Nome Banco", "Banco", "ID Conta Bancária", "Código Agência", "Número Conta Bancária", "Nome Conta Bancária", "Nome Reduzido Conta Bancária", "Tipo Conta Bancária", "Conta Bancária", "Banco, Agência e Conta Bancária", "ID Empresa", "ID Filial", "Código Filial", "Nome Filial", "Filial", "Registro Ativo Banco", "Reg. Bloqueado Conta Bancária", "Registro Ativo Conta Bancária" FROM DIM_BANCO_CONTA_BANCARIA;

SELECT "ID Gestor", "Código Gestor", "Nome Gestor", "Gestor", "ID Segmento", "Código Segmento", "Nome Segmento", "Segmento", "ID Grupo Contrato", "Código Grupo Contrato", "Nome Grupo Contrato", "Grupo Contrato", "ID Contrato", "Código Contrato", "Nome Contrato", "Contrato", "ID Centro de Custo", "Código Centro de Custo", "Nome Centro de Custo", "Centro de Custo", "ID Cliente Final", "Código Cliente Final", "Nome Cliente Final", "Cliente Final", "ID Tipo Contrato", "Código Tipo Contrato", "Nome Tipo Contrato", "Tipo Contrato", "Registro Ativo Gestor", "Registro Ativo Segmento", "Registro Ativo Grupo Contrato", "Registro Ativo Contrato", "Registro Ativo Centro de Custo", "Registro Ativo Cliente Final" FROM DIM_CENTRO_CUSTO;

SELECT "ID Cliente", "Código Cliente", "Loja Cliente", "Razão Social", "Nome Fantasia", "Cliente Razão Social", "Cliente Nome Fantasia", "Cliente CNPJ Razão Social", "Cliente CNPJ Nome Fantasia", "Raíz CNPJ", "CNPJ", "Inscrição Estadual", "Inscrição Municipal", "Endereço", "Bairro", "CEP", "Código Tipo Pessoa", "Descrição Tipo Pessoa", "Tipo Pessoa", "ID Município", "Código Município", "Nome Município", "Município", "ID UF", "Sigla UF", "Nome UF", "UF", "Região", "Sigla País", "Registro Bloqueado Cliente", "Registro Ativo Cliente", "Registro Bloqueado Pessoa", "Registro Ativo Pessoa", "Registro Ativo Município", "Registro Ativo UF" FROM DIM_CLIENTE;

SELECT "ID Colaborador", "Matrícula Colaborador", "Nome Colaborador", "Data Nascimento", "Data Admissão", "Data Transferência", "Data Desligamento", "Empresa Atual", "Filial Atual", "Seção Atual", "Função Atual", "Situação Atual", "Registro Ativo Colaborador" FROM DIM_COLABORADOR;

SELECT "ID Conta Orçamentária", "Código Conta Orçamentária", "Nome Conta Orçamentária", "Conta Orçamentária", "Código Conta Orçamentária N1", "Nome Conta Orçamentária N1", "Conta Orçamentária N1", "Código Conta Orçamentária N2", "Nome Conta Orçamentária N2", "Conta Orçamentária N2", "Código Conta Orçamentária N3", "Nome Conta Orçamentária N3", "Conta Orçamentária N3", "Nível", "Tipo", "Condição", "Código Natureza", "Reg. Ativo Conta Orçamentária" FROM DIM_CONTA_ORCAMENTARIA;

SELECT "ID Conta Orçamentária", "ID Natureza" FROM DIM_CONTA_ORCM_NATUREZA;

SELECT "ID Empresa", "Código Empresa", "Nome Empresa", "Empresa", "ID Filial", "Código Filial", "Nome Filial", "Filial", "Raíz CNPJ", "CNPJ", "Registro Ativo Empresa", "Registro Ativo Filial" FROM DIM_EMPRESA_FILIAL;

SELECT "ID Evento", "Código Evento", "Nome Evento", "Evento", "ID Tipo Evento", "Código Tipo Evento", "Nome Tipo Evento", "Tipo Evento", "Registro Ativo Evento", "Registro Ativo Tipo Evento" FROM DIM_EVENTO;

SELECT "ID Fornecedor", "Código Fornecedor", "Loja Fornecedor", "Razão Social", "Nome Fantasia", "Fornecedor Razão Social", "Fornecedor Nome Fantasia", "Fornecedor CNPJ Razão Social", "Fornecedor CNPJ Nome Fantasia", "Raíz CNPJ", "CNPJ", "Inscrição Estadual", "Inscrição Municipal", "Endereço", "Bairro", "CEP", "Código Tipo Pessoa", "Descrição Tipo Pessoa", "Tipo Pessoa", "ID Município", "Código Município", "Nome Município", "Município", "ID UF", "Sigla UF", "Nome UF", "UF", "Região", "Sigla País", "ID País", "Código País", "Nome País", "País", "Registro Bloqueado Fornecedor", "Registro Ativo Fornecedor", "Registro Bloqueado Pessoa", "Registro Ativo Pessoa", "Registro Ativo Município", "Registro Ativo UF", "Registro Ativo País" FROM DIM_FORNECEDOR;

SELECT "ID Função Colaborador", "Código Função Colaborador", "Descrição Função Colaborador", "Função Colaborador", "Reg. Ativo Função Colaborador" FROM DIM_FUNCAO_COLABORADOR;

SELECT "ID Moeda", "Código Moeda", "Nome Moeda", "Moeda", "Registro Ativo Moeda" FROM DIM_MOEDA;

SELECT "ID Natureza N1", "Código Natureza N1", "Nome Natureza N1", "Natureza N1", "ID Natureza N2", "Código Natureza N2", "Nome Natureza N2", "Natureza N2", "ID Natureza N3", "Código Natureza N3", "Nome Natureza N3", "Natureza N3", "ID Natureza N4", "Código Natureza N4", "Nome Natureza N4", "Natureza N4", "ID Natureza", "Código Natureza", "Nome Natureza", "Natureza", "Registro Bloqueado Natureza", "Registro Ativo Natureza" FROM DIM_NATUREZA;

SELECT "E-mail", "ID Contrato", "Código Contrato", "ID Centro de Custo", "Código Centro de Custo" FROM DIM_RLS_CENTRO_CUSTO;

SELECT "ID Seção", "Código Seção", "Nome Seção", "Seção", "Registro Ativo Seção" FROM DIM_SECAO;

SELECT "ID Situação Colaborador", "Código Situação Colaborador", "Descrição Situação Colaborador", "Situação Colaborador", "Reg. Ativo Situação Colab." FROM DIM_SITUACAO_COLABORADOR;

SELECT "Data", "Dia Útil", "Dia da Semana", "Nome do Dia da Semana", "Semana do Mês", "Descrição da Semana do Mês", "Semana do Ano", "Descrição da Semana do Ano", "Descrição Curta do Período", "Descrição Longa do Período", "Ano e Mês", "Mês e Ano", "Ano", "Mês", "Nome do Mês", "Trimestre", "Descrição do Trimestre", "Descrição Curta do Trimestre", "Descrição Longa do Trimestre", "Semestre", "Descrição do Semestre", "Descrição Curta do Semestre", "Descrição Longa do Semestre" FROM DIM_TEMPO;

SELECT "ID Tipo Mov. Bancária", "Código Tipo Mov. Bancária", "Descrição Tipo Mov. Bancária", "Tipo Mov. Bancária", "Reg. Ativo Tipo Mov. Bancária" FROM DIM_TIPO_MOV_BANCARIA;

SELECT "ID Tipo Processo", "Código Tipo Processo", "Descrição Tipo Processo", "Tipo Processo", "Registro Ativo Tipo Processo" FROM DIM_TIPO_PROCESSO;

SELECT "ID Tipo Recebimento", "Código Tipo Recebimento", "Descrição Tipo Recebimento", "Tipo Recebimento", "Reg. Ativo Tipo Recebimento" FROM DIM_TIPO_RECEBIMENTO_COLAB;

SELECT "ID Tipo Título", "Código Tipo Título", "Descrição Tipo Título", "Tipo Título", "Registro Ativo Tipo Título" FROM DIM_TIPO_TITULO;

--FATOS

SELECT "Operação Mov. Bancária Bx. Tit. Imposto", "ID Filial", "Data Mov. Bancária Bx. Tit. Imposto", "ID Tipo Mov. Bancária Bx. Tit. Imposto", "Código Numerário", "Número Documento", "Situação Documento", "Histórico", "Beneficiário", "ID Moeda", "Valor Mov. Bancária Bx. Tit. Imposto", "ID Conta Bancária", "Banco Virtual", "ID Natureza", "ID Centro de Custo", "Rateio Financeiro", "Rateio Contábil", "ID Empresa Filial Origem", "ID Empresa Filial Destino", "ID Conta Bancária Origem", "ID Conta Bancária Destino", "Prefixo Título", "Número Título", "Parcela Título", "ID Tipo Título", "ID Fornecedor Título", "ID Município Fornecedor Título", "ID País Fornecedor Título", "ID Cliente Título", "ID Município Cliente Título", "Data Emissão Título", "Data Vencimento Título", "Número Lote", "Motivo Baixa", "Sequência Baixa", "Número Cheque", "ID Processo", "Valor Título", "Registro Ativo" FROM FATO_BAIXA_TITULO_IMPOSTO;

SELECT "ID Fornecedor", "ID Contrato", "Valor Saldo", "Data Último Pagamento" FROM FATO_CONFISSAO_DIVIDA;

SELECT "Chave Título", "ID Filial", "Prefixo", "Número", "Parcela", "ID Tipo Título", "Código Tipo Título", "ID Fornecedor", "Código Fornecedor", "Código Loja Fornecedor", "ID Município", "ID País", "ID Natureza", "ID Moeda", "ID Centro de Custo", "Data Emissão", "Data Vencimento", "Data Vencimento Real", "Data Última Baixa", "Histórico", "Taxa Moeda", "Valor Título", "Valor Título Real", "Valor Última Baixa", "Valor Saldo", "Valor Saldo Real", "Valor Pago Real", "Rateio", "Protesto", "Valor Protesto", "Confissão Dívida", "Processo", "ID Processo", "Dação", "Adiantamento Cheque", "Adiantamento Compensado", "Registro Ativo", "Valor ISS", "Valor INSS", "Valor IRRF", "Valor PIS", "Valor COFINS", "Valor CSLL", "Valor Saldo Acréscimo", "Valor Saldo Decréscimo", "Chave Título Pai", "Valor Multa", "Valor Juros", "Valor Desconto", "Deduz INSS", "Valor SEST", "Valor Base", "Valor Acréscimo", "Valor Decréscimo", "Valor Saldo Real Final", "Valor Total Tributos Reter", "Valor Total Tributos Reter Emissão", "Valor Total Tributos Reter Baixa", "Valor Líquido Título Real", "Valor Abatimentos", "Valor Tributos Retidos Baixa", "Valor Tributos Pendentes Retenção Baixa", "Valor Líquido Pago Real", "Valor Líquido Saldo Real" FROM FATO_CONTA_PAGAR_RATEADO;

SELECT "ID Filial", "Prefixo", "Número", "Parcela", "ID Tipo Título", "Código Tipo Título", "ID Cliente", "Código Cliente", "Código Loja Cliente", "ID Município", "ID Natureza", "ID Moeda", "ID Centro de Custo", "Data Emissão", "Data Vencimento", "Data Vencimento Real", "Data Última Baixa", "Histórico", "Taxa Moeda", "Valor Título", "Valor Título Real", "Valor Última Baixa", "Valor Saldo", "Valor Saldo Real", "Valor Recebido Real", "Rateio", "Registro Ativo", "Valor ISS", "Valor INSS", "Valor IRRF", "Valor PIS", "Valor COFINS", "Valor CSLL", "Valor Saldo Acréscimo", "Valor Saldo Decréscimo", "Chave Título Pai", "Valor Total Tributos Reter", "Valor Líquido Título Real", "Valor Abatimentos", "Valor Tributos Retidos Baixa", "Valor Líquido Recebido Real", "Valor Líquido Saldo Real" FROM FATO_CONTA_RECEBER_RATEADO;

SELECT "Ano e Mês", "Data Competência", "Período Pagamento", "ID Colaborador", "ID Evento", "Valor Evento", "Quantidade Evento", "Registro Ativo" FROM FATO_FOLHA_PAGAMENTO;

SELECT "Ano e Mês", "Data Competência", "ID Colaborador", "ID Filial", "ID Seção", "ID Centro de Custo", "ID Função Colaborador", "ID Situação Colaborador", "ID Tipo Recebimento", "ID Banco", "Valor Salário Base", "Registro Ativo Colaborador" FROM FATO_HISTORICO_COLABORADOR;

SELECT "Operação Mov. Bancária", "ID Filial", "Data Mov. Bancária", "ID Tipo Mov. Bancária", "Código Numerário", "Número Documento", "Situação Documento", "Histórico", "Beneficiário", "ID Moeda", "Valor Mov. Bancária", "ID Conta Bancária", "Banco Virtual", "ID Natureza", "ID Centro de Custo", "Rateio Financeiro", "Rateio Contábil", "ID Empresa Filial Origem", "ID Empresa Filial Destino", "ID Conta Bancária Origem", "ID Conta Bancária Destino", "Prefixo Título", "Número Título", "Parcela Título", "ID Tipo Título", "ID Fornecedor Título", "ID Município Fornecedor Título", "ID País Fornecedor Título", "ID Cliente Título", "ID Município Cliente Título", "Data Emissão Título", "Data Vencimento Título", "Número Lote", "Motivo Baixa", "Sequência Baixa", "Número Cheque", "ID Processo", "Número Registro", "Registro Ativo" FROM FATO_MOV_BANCARIA_RATEADA;

SELECT "ID Orçamento", "Código Orçamento", "Descrição Orçamento", "Orçamento", "Data Início Período", "Data Fim Período", "Número Versão", "Número Item", "ID Conta Orçamentária", "ID Contrato", "Valor Total Orçado Original", "Valor Total Orçado Atual", "Valor Variação Orig. x Atual" FROM FATO_ORCAMENTO_TOTAL;

SELECT "ID Processo", "Código Processo", "Nome Autor", "ID Tipo Processo", "ID Centro de Custo", "ID Fornecedor", "Valor Total", "Registro Ativo" FROM FATO_PROCESSO;
