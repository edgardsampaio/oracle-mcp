package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cimdMaxBytes     = 64 << 10
	cimdFetchTimeout = 5 * time.Second
	cimdDefaultTTL   = 5 * time.Minute
)

// CIMDDocument é o Client ID Metadata Document.
type CIMDDocument struct {
	ClientID     string
	ClientName   string
	RedirectURIs []string
}

type cimdCacheEntry struct {
	doc     CIMDDocument
	expires time.Time
}

type cimdCache struct {
	mu      sync.Mutex
	entries map[string]cimdCacheEntry
}

func (c *cimdCache) get(id string, now time.Time) (CIMDDocument, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[id]
	if !ok || !e.expires.After(now) {
		return CIMDDocument{}, false
	}
	return e.doc, true
}

func (c *cimdCache) put(id string, doc CIMDDocument, expires time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = map[string]cimdCacheEntry{}
	}
	c.entries[id] = cimdCacheEntry{doc: doc, expires: expires}
}

// CIMDFetcher busca e valida um documento de cliente.
type CIMDFetcher struct {
	Client *http.Client
	Lookup func(ctx context.Context, host string) ([]net.IP, error)
}

// ValidateCIMDURL recusa URL que não seja https com path, ou que aponte para rede interna.
func ValidateCIMDURL(raw string, lookup func(ctx context.Context, host string) ([]net.IP, error)) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("client_id inválido")
	}
	if u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("client_id CIMD deve ser uma URL https com path")
	}
	if u.Path == "" || u.Path == "/" {
		return fmt.Errorf("client_id CIMD deve conter um path")
	}
	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		if !publicIP(ip) {
			return fmt.Errorf("client_id CIMD aponta para endereço interno")
		}
		return nil
	}
	if lookup == nil {
		lookup = func(ctx context.Context, host string) ([]net.IP, error) {
			return net.DefaultResolver.LookupIP(ctx, "ip", host)
		}
	}
	ips, err := lookup(context.Background(), host)
	if err != nil || len(ips) == 0 {
		return fmt.Errorf("não foi possível resolver o client_id CIMD")
	}
	for _, ip := range ips {
		if !publicIP(ip) {
			return fmt.Errorf("client_id CIMD aponta para endereço interno")
		}
	}
	return nil
}

func publicIP(ip net.IP) bool {
	if ip == nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil && ip4[0] == 169 && ip4[1] == 254 {
		return false
	}
	return true
}

// Fetch baixa o documento e confere o client_id.
func (f CIMDFetcher) Fetch(ctx context.Context, clientID string) (CIMDDocument, time.Time, error) {
	lookup := f.Lookup
	if err := ValidateCIMDURL(clientID, lookup); err != nil {
		return CIMDDocument{}, time.Time{}, err
	}
	client := f.Client
	if client == nil {
		client = &http.Client{Timeout: cimdFetchTimeout, CheckRedirect: func(*http.Request, []*http.Request) error {
			return fmt.Errorf("redirecionamento de CIMD recusado")
		}}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, clientID, nil)
	if err != nil {
		return CIMDDocument{}, time.Time{}, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return CIMDDocument{}, time.Time{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return CIMDDocument{}, time.Time{}, fmt.Errorf("documento CIMD indisponível")
	}
	if ct := res.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		return CIMDDocument{}, time.Time{}, fmt.Errorf("documento CIMD sem JSON")
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, cimdMaxBytes+1))
	if err != nil {
		return CIMDDocument{}, time.Time{}, err
	}
	if len(body) > cimdMaxBytes {
		return CIMDDocument{}, time.Time{}, fmt.Errorf("documento CIMD grande demais")
	}
	var raw struct {
		ClientID     string   `json:"client_id"`
		ClientName   string   `json:"client_name"`
		RedirectURIs []string `json:"redirect_uris"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return CIMDDocument{}, time.Time{}, err
	}
	if raw.ClientID != clientID || raw.ClientName == "" || len(raw.RedirectURIs) == 0 {
		return CIMDDocument{}, time.Time{}, fmt.Errorf("documento CIMD inválido")
	}
	doc := CIMDDocument{ClientID: raw.ClientID, ClientName: raw.ClientName, RedirectURIs: raw.RedirectURIs}
	return doc, time.Now().Add(cacheTTL(res.Header.Get("Cache-Control"))), nil
}

func cacheTTL(header string) time.Duration {
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "no-store" || part == "no-cache" {
			return 0
		}
		if strings.HasPrefix(part, "max-age=") {
			n, err := strconv.Atoi(strings.TrimPrefix(part, "max-age="))
			if err == nil && n > 0 {
				return time.Duration(n) * time.Second
			}
		}
	}
	return cimdDefaultTTL
}

func isCIMDID(clientID string) bool {
	u, err := url.Parse(clientID)
	return err == nil && u.Scheme == "https" && u.Host != "" && u.Path != "" && u.Path != "/"
}
