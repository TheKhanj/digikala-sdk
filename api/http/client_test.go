package http

import (
	"log"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	proxy := HttpProxy{
		Username: os.Getenv("HTTP_PROXY_USERNAME"),
		Password: os.Getenv("HTTP_PROXY_PASSWORD"),
		Address:  os.Getenv("HTTP_PROXY_ADDRESS"),
	}
	t.Log(proxy)
	c, err := NewClient(
		WithHttpProxy(proxy),
		WithCurlLog(log.Default()),
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	res, err := c.Get("https://google.com")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer res.Body.Close()
}
