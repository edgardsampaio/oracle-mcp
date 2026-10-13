package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/term"

	"oracle-mcp/internal/auth"
	"oracle-mcp/internal/config"
	"oracle-mcp/internal/identity/sqlite"
)

var errPasswordMismatch = errors.New("as senhas não conferem")

func runAdmin(logger *slog.Logger, args []string) error {
	if len(args) == 0 || args[0] != "create" || len(args) != 1 {
		return fmt.Errorf("use: oracle-mcp admin create")
	}
	if err := prepareRuntime(); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := applyRuntimePaths(&cfg); err != nil {
		return err
	}
	store, err := sqlite.Open(cfg.Auth.SQLitePath)
	if err != nil {
		return err
	}
	defer store.Close()
	svc, err := auth.New(store, cfg, logger)
	if err != nil {
		return err
	}
	username, email, password, err := promptAdmin(os.Stdin, os.Stderr)
	if err != nil {
		return err
	}
	if err := svc.CreateInitialAdmin(context.Background(), username, email, password); err != nil {
		return err
	}
	logger.Info("administrador criado", "username", username)
	return nil
}

func promptAdmin(in *os.File, out io.Writer) (string, string, string, error) {
	if !term.IsTerminal(int(in.Fd())) {
		return "", "", "", fmt.Errorf("execute em um terminal; a senha não é aceita por argumento")
	}
	fmt.Fprint(out, "Usuário: ")
	reader := bufio.NewReader(in)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	username = strings.TrimSpace(username)
	fmt.Fprint(out, "E-mail: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	email = strings.TrimSpace(email)
	fmt.Fprint(out, "Senha: ")
	first, err := term.ReadPassword(int(in.Fd()))
	fmt.Fprintln(out)
	if err != nil {
		return "", "", "", err
	}
	fmt.Fprint(out, "Confirme a senha: ")
	second, err := term.ReadPassword(int(in.Fd()))
	fmt.Fprintln(out)
	if err != nil {
		return "", "", "", err
	}
	if string(first) != string(second) {
		return "", "", "", errPasswordMismatch
	}
	return username, email, string(first), nil
}

// ensureInitialAdmin pede o administrador só quando ainda não existe um ativo.
func ensureInitialAdmin(ctx context.Context, admins int, out io.Writer, prompt func() (string, string, string, error), create func(context.Context, string, string, string) error) (string, error) {
	if admins > 0 {
		return "", nil
	}
	return bootstrapAdmin(ctx, out, prompt, create)
}

// bootstrapAdmin pede usuário, senha e confirmação.
// Senha divergente ou rejeitada pela validação pede de novo. Falha de terminal encerra.
func bootstrapAdmin(ctx context.Context, out io.Writer, prompt func() (string, string, string, error), create func(context.Context, string, string, string) error) (string, error) {
	fmt.Fprintln(out, "Primeiro acesso. Informe o administrador.")
	for {
		username, email, password, err := prompt()
		if err != nil {
			if errors.Is(err, errPasswordMismatch) {
				fmt.Fprintln(out, err.Error())
				continue
			}
			return "", err
		}
		if err := create(ctx, username, email, password); err != nil {
			if errors.Is(err, auth.ErrValidation) {
				fmt.Fprintln(out, err.Error())
				continue
			}
			return "", err
		}
		return username, nil
	}
}
