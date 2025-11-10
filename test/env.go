package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thekhanj/digikala-sdk/browser"
	"github.com/thekhanj/digikala-sdk/customer/auth"
	"github.com/thekhanj/digikala-sdk/test/config"
)

func NewTestEnv() (*TestEnv, error) {
	env := &TestEnv{
		Config: config.Config{},

		configPath: config.GetConfigPath(),
	}

	err := env.Read()
	if err != nil {
		return nil, err
	}

	return env, nil
}

type TestEnv struct {
	Config config.Config

	configPath string
}

func (this *TestEnv) CustomerAssertAuthorized(ctx context.Context) error {
	cache := &this.Config.Customer.Cache
	if cache.Cookies != nil {
		return nil
	}

	cookies, err := this.CustomerAuthorize(ctx)
	if err != nil {
		return err
	}

	err = cache.SetCookies(cookies)
	if err != nil {
		return err
	}

	return this.Write()
}

func (this *TestEnv) CustomerAuthorize(
	ctx context.Context,
) ([]*http.Cookie, error) {
	p := this.Config.Customer.Auth.Proxy
	var testProxy config.HttpProxy
	if p != nil {
		var ok bool
		testProxy, ok = this.Config.Proxies[*p]
		if !ok {
			return nil, fmt.Errorf("Proxy %s not found", *p)
		}
	}

	b := browser.NewBrowser(
		browser.WithHeadless(false),
		func(b *browser.Browser) {
			if p == nil {
				return
			}

			user := testProxy.Username
			password := testProxy.Password
			address := testProxy.Address

			proxy := browser.HttpProxy{
				User: user, Password: password, Address: address,
			}
			browser.WithHttpProxy(proxy)(b)
		},
	)
	b.Run(ctx)

	tab, err := b.NewTab()
	if err != nil {
		return nil, err
	}

	auth := auth.NewAuth(tab)

	username := this.Config.Customer.Auth.Username
	password := this.Config.Customer.Auth.Password
	cookies, err := auth.Login(username, password)
	if err != nil {
		return nil, err
	}

	return cookies, nil
}

func (this *TestEnv) Read() error {
	return config.ReadConfig(&this.Config, this.configPath)
}

func (this *TestEnv) Write() error {
	return config.WriteConfig(&this.Config, this.configPath)
}
