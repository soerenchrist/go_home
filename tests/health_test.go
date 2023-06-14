package tests

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestHealthRoute(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/health")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"ok\"}", w.Body.String())
}
