package oracle

import (
	"strings"
	"testing"
	"time"

	"oracle-mcp/internal/config"
)

func TestDataSourceModes(t *testing.T) {
	base := config.Oracle{
		Host:         "172.16.0.15",
		Port:         1526,
		User:         "CONSULTOR",
		Password:     "segredo-local",
		QueryTimeout: 30 * time.Second,
	}
	sid := base
	sid.SID = "bi"
	gotSID := dataSource(sid)
	if !strings.Contains(gotSID, "SID=bi") || strings.Contains(gotSID, "/BI") {
		t.Fatal("dsn de sid incorreto")
	}
	service := base
	service.ServiceName = "BI"
	gotService := dataSource(service)
	if strings.Contains(gotService, "SID=") || !strings.Contains(gotService, "/BI") {
		t.Fatal("dsn de service name incorreto")
	}
}
