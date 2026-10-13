// Package boasvindas apresenta o agente e as logos configuradas.
// Não importa dimensão, fato nem Oracle.
package boasvindas

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"oracle-mcp/internal/bi"
)

const (
	// ToolName é o nome MCP desta tool.
	ToolName = "boas_vindas"
	fileName = "welcome.json"

	description = `Apresenta o agente no início da conversa. Chame esta tool antes de consultar dados e mostre ao usuário a mensagem e as duas logos: a da consultoria que implanta o BI e a da empresa cliente. Não invente logo que a tool não devolveu.`
)

// ErrUnavailable é a ausência do JSON de boas-vindas, sem caminho no texto.
var ErrUnavailable = errors.New("boas-vindas indisponível")

var errPath = errors.New("arquivo de logo recusado")

type fileDoc struct {
	Nome          string `json:"nome"`
	Personalidade string `json:"personalidade"`
	Mensagem      string `json:"mensagem"`
	Logos         []logo `json:"logos"`
}

type logo struct {
	Papel   string `json:"papel"`
	Arquivo string `json:"arquivo"`
	Legenda string `json:"legenda"`
}

// Image é uma logo lida da pasta de configuração.
type Image struct {
	Papel   string
	Legenda string
	MIME    string
	Dados   []byte
}

// Result é o texto e as logos encontradas. Os bytes ficam só no bloco de imagem do MCP.
type Result struct {
	Texto   string
	Imagens []Image
}

// saida é o resultado estruturado. Não leva bytes: o esquema de []byte não coincide com o base64 do JSON.
type saida struct {
	Texto string `json:"texto"`
}

// Register publica a tool. dir é a pasta que contém o JSON e as imagens.
func Register(server *sdk.Server, dir string, audit bi.Auditor, log *slog.Logger) {
	sdk.AddTool(server, &sdk.Tool{Name: ToolName, Description: description}, func(ctx context.Context, req *sdk.CallToolRequest, _ struct{}) (*sdk.CallToolResult, saida, error) {
		result, err := Greet(ctx, dir, audit, log, userID(req))
		if err != nil {
			return nil, saida{}, err
		}
		call, out := toolResult(result)
		return call, out, nil
	})
}

// Greet lê o JSON e as imagens da pasta. Logo ausente entra no texto, sem derrubar a mensagem.
func Greet(ctx context.Context, dir string, audit bi.Auditor, log *slog.Logger, userID string) (result Result, err error) {
	start := time.Now()
	defer bi.Finish(ctx, audit, log, userID, ToolName, start, &err)
	if userID == "" {
		return Result{}, bi.ErrIdentity
	}
	doc, err := readDoc(dir)
	if err != nil {
		return Result{}, ErrUnavailable
	}
	var b strings.Builder
	if doc.Mensagem != "" {
		b.WriteString(doc.Mensagem)
		b.WriteString("\n\n")
	}
	if doc.Nome != "" {
		b.WriteString(doc.Nome)
		b.WriteString("\n\n")
	}
	if doc.Personalidade != "" {
		b.WriteString(doc.Personalidade)
	}
	for _, item := range doc.Logos {
		img, readErr := readLogo(dir, item)
		if readErr != nil {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString("Logo ")
			b.WriteString(label(item))
			b.WriteString(" não encontrada: ")
			b.WriteString(filepath.Base(strings.TrimSpace(item.Arquivo)))
			continue
		}
		result.Imagens = append(result.Imagens, img)
	}
	result.Texto = b.String()
	return result, nil
}

func toolResult(result Result) (*sdk.CallToolResult, saida) {
	content := []sdk.Content{&sdk.TextContent{Text: result.Texto}}
	for _, img := range result.Imagens {
		content = append(content, &sdk.ImageContent{Data: img.Dados, MIMEType: img.MIME})
	}
	return &sdk.CallToolResult{Content: content}, saida{Texto: result.Texto}
}

func readDoc(dir string) (fileDoc, error) {
	if strings.TrimSpace(dir) == "" {
		return fileDoc{}, ErrUnavailable
	}
	raw, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		return fileDoc{}, err
	}
	var doc fileDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fileDoc{}, err
	}
	return doc, nil
}

func readLogo(dir string, item logo) (Image, error) {
	path, err := safeJoin(dir, item.Arquivo)
	if err != nil {
		return Image{}, err
	}
	mime, ok := mimeOf(path)
	if !ok {
		return Image{}, errPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Image{}, err
	}
	return Image{Papel: item.Papel, Legenda: item.Legenda, MIME: mime, Dados: data}, nil
}

func safeJoin(dir, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.Contains(name, "..") || strings.ContainsAny(name, `/\`) {
		return "", errPath
	}
	return filepath.Join(dir, name), nil
}

func mimeOf(path string) (string, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "image/png", true
	case ".jpg", ".jpeg":
		return "image/jpeg", true
	case ".webp":
		return "image/webp", true
	default:
		return "", false
	}
}

// ConsultoriaLogo devolve a imagem de papel consultoria. Ausência de JSON, papel ou arquivo devolve ok falso.
func ConsultoriaLogo(dir string) (mime, legenda string, data []byte, ok bool) {
	doc, err := readDoc(dir)
	if err != nil {
		return "", "", nil, false
	}
	for _, item := range doc.Logos {
		if item.Papel != "consultoria" {
			continue
		}
		img, err := readLogo(dir, item)
		if err != nil {
			return "", "", nil, false
		}
		legenda = img.Legenda
		if legenda == "" {
			legenda = "Consultoria"
		}
		return img.MIME, legenda, img.Dados, true
	}
	return "", "", nil, false
}

func label(item logo) string {
	if item.Papel != "" {
		return item.Papel
	}
	return "sem papel"
}

func userID(req *sdk.CallToolRequest) string {
	if req == nil || req.Extra == nil || req.Extra.TokenInfo == nil {
		return ""
	}
	return req.Extra.TokenInfo.UserID
}
