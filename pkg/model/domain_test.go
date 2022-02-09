package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSha2Hash(t *testing.T) {
	assert.EqualValues(t, "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25", getSha2Hash("Test"))
}
