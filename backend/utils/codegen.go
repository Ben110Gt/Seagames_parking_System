package util

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateTicketCode returns TKT-XXXXXXXX
func GenerateTicketCode() string {
	id := strings.ToUpper(uuid.New().String()[:8])
	return "TKT-" + id
}

// GenerateMemberCode returns MBR-XXXXXXXX
func GenerateMemberCode() string {
	id := strings.ToUpper(uuid.New().String()[:8])
	return "MBR-" + id
}
