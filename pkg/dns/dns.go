package dns

import (
	"fmt"
	"strings"
)

// Result hold the data for a update session
type Result struct {
	Error   error
	Message string
}

// Record is the data needed to invoke update
type Record struct {
	Hostname string
	TTL      string
	Type     string
	Data     string
}

// Server implements the handler for a backend
type Server interface {
	Update(d Record) Result
	Delete(d Record) Result
}

func getZone(hostname string) string {
	domainParts := strings.Split(strings.TrimRight(hostname, "."), ".")
	return fmt.Sprintf("%s.%s", domainParts[len(domainParts)-2], domainParts[len(domainParts)-1])
}
