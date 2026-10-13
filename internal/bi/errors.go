package bi

import "errors"

var (
	// ErrIdentity é a ausência do usuário do token.
	ErrIdentity = errors.New("identidade ausente")
	// ErrUnavailable é falha de consulta sem detalhe do driver.
	ErrUnavailable = errors.New("consulta indisponível")
	// ErrEmail é usuário autenticado sem e-mail para a segurança por linha.
	ErrEmail = errors.New("usuário sem e-mail")
	// ErrParam é filtro rejeitado antes da consulta.
	ErrParam = errors.New("parâmetro inválido")
)
