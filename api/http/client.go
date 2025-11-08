package http

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/moul/http2curl"
	"github.com/thekhanj/digikala-sdk/api"
)

type ClientOption = func(c *http.Client) error

func WithCookies(cookies []*http.Cookie) ClientOption {
	return func(c *http.Client) error {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}

		c.Jar = jar

		for _, rawUrl := range api.CookiesUrls {
			url, err := url.Parse(rawUrl)
			if err != nil {
				return err
			}
			c.Jar.SetCookies(url, cookies)
		}

		return nil
	}
}

func WithCurlLog(log *log.Logger) ClientOption {
	return func(c *http.Client) error {
		baseTransport := c.Transport
		if baseTransport == nil {
			baseTransport = http.DefaultTransport
		}

		c.Transport = &curlRoundTripper{log, baseTransport}

		return nil
	}
}

type curlRoundTripper struct {
	log           *log.Logger
	baseTransport http.RoundTripper
}

func (this *curlRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cmd, err := http2curl.GetCurlCommand(req)
	if err == nil {
		this.log.Printf("[CURL] %s\n", cmd.String())
	}

	return this.baseTransport.RoundTrip(req)
}

func NewClient(opts ...ClientOption) (*http.Client, error) {
	c := &http.Client{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
