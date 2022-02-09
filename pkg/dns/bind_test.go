package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	b   Bind
	rec Record
)

func TestMain(m *testing.M) {
	b = Bind{
		Debug:   true,
		Enabled: true,
		Host:    "localhost",
		Sync:    true,
	}
	rec = Record{
		Hostname: "test.example.com",
		TTL:      "1800",
		Data:     "1234567890",
		Type:     "TXT",
	}
}

func TestUpdate(t *testing.T) {
	r := b.Update(rec)
	exp := "DEBUG MODE\n\n would invoke:\n\nserver localhost\nzone example.com\nupdate delete test.example.com\nupdate add test.example.com 1800 TXT 1234567890\nsend\n"
	assert.Equal(t, nil, r.Error)
	assert.EqualValues(t, exp, r.Message)
}

func TestDelete(t *testing.T) {
	r := b.Delete(rec)
	exp := "DEBUG MODE\n\n would invoke:\n\nserver localhost\nzone example.com\nupdate delete test.example.com\nsend\n"
	assert.Equal(t, nil, r.Error)
	assert.EqualValues(t, exp, r.Message)
}
