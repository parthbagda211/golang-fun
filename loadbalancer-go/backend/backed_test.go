package backend

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestBackendCreation(t *testing.T) {
	url, _ := url.Parse("http://localhost:8080")
	b := NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	assert.Equal(t, "http://localhost:8080", b.GetUrl().String())
	assert.Equal(t, true, b.IsAlive())
}

func TestBackend_IsAlive(t *testing.T) {
	url, _ := url.Parse("http://localhost:8080")
	b := NewBackend(url, httputil.NewSingleHostReverseProxy(url))
	b.SetAlive(b.IsAlive())
	assert.Equal(t, false, b.IsAlive())
}
