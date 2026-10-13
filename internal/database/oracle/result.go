package oracle

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
	"unicode/utf8"

	go_ora "github.com/sijms/go-ora/v2"
)

// Result é o JSON previsível de uma consulta cadastrada.
type Result struct {
	Columns   []string
	Rows      []map[string]any
	Truncated bool
}

// rowSource é o pedaço de *sql.Rows usado na coleta.
type rowSource interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Err() error
}

// Collect aplica normalização, máximo de linhas e máximo de bytes.
func Collect(src rowSource, maxRows, maxBytes int) (Result, error) {
	if maxRows <= 0 || maxBytes <= 0 {
		return Result{}, fmt.Errorf("limite de consulta inválido")
	}
	cols, err := src.Columns()
	if err != nil {
		return Result{}, err
	}
	out := Result{Columns: cols, Rows: []map[string]any{}}
	used := 0
	for src.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := src.Scan(ptrs...); err != nil {
			return Result{}, err
		}
		row := make(map[string]any, len(cols))
		for i, col := range cols {
			name := col
			if name == "" {
				name = fmt.Sprintf("col_%d", i+1)
			}
			value, err := Normalize(vals[i])
			if err != nil {
				return Result{}, err
			}
			row[name] = value
		}
		raw, err := json.Marshal(row)
		if err != nil {
			return Result{}, err
		}
		if used+len(raw) > maxBytes {
			if len(out.Rows) == 0 {
				return Result{}, fmt.Errorf("a primeira linha excede o limite de bytes")
			}
			out.Truncated = true
			break
		}
		if len(out.Rows) >= maxRows {
			out.Truncated = true
			break
		}
		out.Rows = append(out.Rows, row)
		used += len(raw)
	}
	if err := src.Err(); err != nil {
		return Result{}, err
	}
	return out, nil
}

// Normalize converte um valor Oracle para JSON previsível.
// NUMBER que não cabe em inteiro seguro sai como string.
func Normalize(v any) (any, error) {
	switch t := v.(type) {
	case nil:
		return nil, nil
	case go_ora.Clob:
		if !t.Valid && t.String == "" {
			return nil, nil
		}
		return t.String, nil
	case *go_ora.Clob:
		if t == nil {
			return nil, nil
		}
		return Normalize(*t)
	case go_ora.NClob:
		if !t.Valid && t.String == "" {
			return nil, nil
		}
		return t.String, nil
	case *go_ora.NClob:
		if t == nil {
			return nil, nil
		}
		return Normalize(go_ora.NClob(*t))
	case go_ora.Number:
		return t.String()
	case *go_ora.Number:
		if t == nil {
			return nil, nil
		}
		return t.String()
	case time.Time:
		return t.UTC().Format(time.RFC3339Nano), nil
	case []byte:
		if utf8.Valid(t) {
			return string(t), nil
		}
		return strconv.QuoteToASCII(string(t)), nil
	case string:
		return t, nil
	case int:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case float64:
		return numberFromFloat(t), nil
	case float32:
		return numberFromFloat(float64(t)), nil
	case *big.Rat:
		if t == nil {
			return nil, nil
		}
		return t.RatString(), nil
	case json.Number:
		return t.String(), nil
	default:
		return nil, fmt.Errorf("tipo de coluna não suportado")
	}
}

func numberFromFloat(f float64) any {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}
	if f == math.Trunc(f) && f >= -1<<53 && f <= 1<<53 {
		return int64(f)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
