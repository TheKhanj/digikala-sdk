package http

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/moul/http2curl"
	"github.com/thekhanj/digikala-sdk/api"
)

type HttpProxy struct {
	Username string
	Password string
	Address  string
}

type ClientOption = func(c *http.Client) error

func WithHttpProxy(proxy HttpProxy) ClientOption {
	return func(c *http.Client) error {
		u, err := url.Parse(proxy.Address)
		if err != nil {
			return err
		}
		u.User = url.UserPassword(proxy.Username, proxy.Password)

		c.Transport = &http.Transport{
			Proxy:           http.ProxyURL(u),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return nil
	}
}

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

// Must be called after WithHttpProxy.
// I couldn't find a way to make the order independent,
// and honestly, it's not worth the effort.
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
