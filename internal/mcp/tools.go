package mcp

import (
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"oracle-mcp/internal/bi"
	"oracle-mcp/internal/boasvindas"
	bancocontabancaria "oracle-mcp/internal/dimensao/bancocontabancaria"
	centrocusto "oracle-mcp/internal/dimensao/centrocusto"
	cliente "oracle-mcp/internal/dimensao/cliente"
	colaborador "oracle-mcp/internal/dimensao/colaborador"
	contaorcamentaria "oracle-mcp/internal/dimensao/contaorcamentaria"
	contaorcmnatureza "oracle-mcp/internal/dimensao/contaorcmnatureza"
	empresafilial "oracle-mcp/internal/dimensao/empresafilial"
	evento "oracle-mcp/internal/dimensao/evento"
	fornecedor "oracle-mcp/internal/dimensao/fornecedor"
	funcaocolaborador "oracle-mcp/internal/dimensao/funcaocolaborador"
	moeda "oracle-mcp/internal/dimensao/moeda"
	natureza "oracle-mcp/internal/dimensao/natureza"
	rlscentrocusto "oracle-mcp/internal/dimensao/rlscentrocusto"
	secao "oracle-mcp/internal/dimensao/secao"
	situacaocolaborador "oracle-mcp/internal/dimensao/situacaocolaborador"
	tempo "oracle-mcp/internal/dimensao/tempo"
	tipomovbancaria "oracle-mcp/internal/dimensao/tipomovbancaria"
	tipoprocesso "oracle-mcp/internal/dimensao/tipoprocesso"
	tiporecebimentocolab "oracle-mcp/internal/dimensao/tiporecebimentocolab"
	tipotitulo "oracle-mcp/internal/dimensao/tipotitulo"
	baixatituloimposto "oracle-mcp/internal/fato/baixatituloimposto"
	confissaodivida "oracle-mcp/internal/fato/confissaodivida"
	contapagar "oracle-mcp/internal/fato/contapagar"
	contareceber "oracle-mcp/internal/fato/contareceber"
	folhapagamento "oracle-mcp/internal/fato/folhapagamento"
	historicocolaborador "oracle-mcp/internal/fato/historicocolaborador"
	movbancaria "oracle-mcp/internal/fato/movbancaria"
	orcamento "oracle-mcp/internal/fato/orcamento"
	processo "oracle-mcp/internal/fato/processo"
)

// businessTools são as consultas de dimensão e fato, sem a boas-vindas.
var businessTools = []string{
	"dim_banco_conta_bancaria",
	"dim_centro_custo",
	"dim_cliente",
	"dim_colaborador",
	"dim_conta_orcamentaria",
	"dim_conta_orcm_natureza",
	"dim_empresa_filial",
	"dim_evento",
	"dim_fornecedor",
	"dim_funcao_colaborador",
	"dim_moeda",
	"dim_natureza",
	"dim_rls_centro_custo",
	"dim_secao",
	"dim_situacao_colaborador",
	"dim_tempo",
	"dim_tipo_mov_bancaria",
	"dim_tipo_processo",
	"dim_tipo_recebimento_colab",
	"dim_tipo_titulo",
	"fato_baixa_titulo_imposto",
	"fato_baixa_titulo_imposto_com_dimensoes",
	"fato_conta_pagar_rateado",
	"fato_conta_pagar_rateado_com_dimensoes",
	"fato_conta_receber_rateado",
	"fato_conta_receber_rateado_com_dimensoes",
	"fato_confissao_divida",
	"fato_confissao_divida_com_dimensoes",
	"fato_folha_pagamento",
	"fato_folha_pagamento_com_dimensoes",
	"fato_historico_colaborador",
	"fato_historico_colaborador_com_dimensoes",
	"fato_mov_bancaria_rateada",
	"fato_mov_bancaria_rateada_com_dimensoes",
	"fato_orcamento_total",
	"fato_orcamento_total_com_dimensoes",
	"fato_processo",
	"fato_processo_com_dimensoes",
}

// ToolCount é a boas-vindas mais as consultas de negócio.
func ToolCount() int { return len(businessTools) + 1 }

func registerTools(server *sdk.Server, env bi.ToolEnv, welcomeDir string) {
	boasvindas.Register(server, welcomeDir, env.Audit, env.Log)
	bancocontabancaria.Register(server, env)
	centrocusto.Register(server, env)
	cliente.Register(server, env)
	colaborador.Register(server, env)
	contaorcamentaria.Register(server, env)
	contaorcmnatureza.Register(server, env)
	empresafilial.Register(server, env)
	evento.Register(server, env)
	fornecedor.Register(server, env)
	funcaocolaborador.Register(server, env)
	moeda.Register(server, env)
	natureza.Register(server, env)
	rlscentrocusto.Register(server, env)
	secao.Register(server, env)
	situacaocolaborador.Register(server, env)
	tempo.Register(server, env)
	tipomovbancaria.Register(server, env)
	tipoprocesso.Register(server, env)
	tiporecebimentocolab.Register(server, env)
	tipotitulo.Register(server, env)
	baixatituloimposto.Register(server, env)
	contapagar.Register(server, env)
	contareceber.Register(server, env)
	confissaodivida.Register(server, env)
	folhapagamento.Register(server, env)
	historicocolaborador.Register(server, env)
	movbancaria.Register(server, env)
	orcamento.Register(server, env)
	processo.Register(server, env)
}
