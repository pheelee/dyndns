package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetZone(t *testing.T) {
	assert.Equal(t, "example.com", getZone("www.example.com"))
	assert.Equal(t, "first.example.com", getZone("www.first.example.com"))
	assert.Equal(t, "second.first.example.com", getZone("www.second.first.lab.example.com"))
}
