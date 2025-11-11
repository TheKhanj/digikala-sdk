package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/thekhanj/digikala-sdk/api"
)

const (
	LoginUrl               = "https://digikala.com/users/login"
	UsernameInputSelector  = "#username"
	UsernameButtonSelector = "#dk-form-buttons"
	PasswordInputSelector  = "#password"
	PasswordButtonSelector = "#dk-form-buttons"
)

type Auth struct {
	tab context.Context

	readyDelay time.Duration
}

func (this *Auth) LoginWithContext(
	ctx context.Context,
	username, password string,
) ([]*http.Cookie, error) {
	subctx, cancel := context.WithCancel(this.tab)

	returned := make(chan struct{})
	defer close(returned)

	go func() {
		select {
		case <-returned:
			break
		case <-ctx.Done():
			break
		}

		cancel()
	}()

	return this.login(subctx, username, password)
}

func (this *Auth) Login(username, password string) ([]*http.Cookie, error) {
	return this.LoginWithContext(context.Background(), username, password)
}

func (this *Auth) login(
	ctx context.Context,
	username, password string,
) ([]*http.Cookie, error) {
	err := this.openLoginPage(ctx)
	if err != nil {
		return nil, err
	}

	err = this.submitUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = this.submitPassword(ctx, password)
	if err != nil {
		return nil, err
	}

	return this.pullCookies(ctx)
}

func (this *Auth) openLoginPage(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.Navigate(LoginUrl),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
}

func (this *Auth) submitUsername(
	ctx context.Context, username string,
) error {
	return chromedp.Run(ctx,
		chromedp.WaitVisible(UsernameButtonSelector, chromedp.ByID),
		chromedp.WaitVisible(UsernameInputSelector, chromedp.ByID),
		chromedp.Sleep(this.readyDelay),
		chromedp.SendKeys(UsernameInputSelector, username, chromedp.ByID),
		chromedp.Click(UsernameButtonSelector, chromedp.ByID),
	)
}

func (this *Auth) submitPassword(
	ctx context.Context, password string,
) error {
	return chromedp.Run(ctx,
		chromedp.WaitVisible(PasswordButtonSelector, chromedp.ByID),
		chromedp.WaitVisible(PasswordInputSelector, chromedp.ByID),
		chromedp.Sleep(this.readyDelay),
		chromedp.SendKeys(PasswordInputSelector, password, chromedp.ByID),
		chromedp.Click(PasswordButtonSelector, chromedp.ByID),
	)
}

func (this *Auth) pullCookies(
	ctx context.Context,
) ([]*http.Cookie, error) {
	ret := make([]*http.Cookie, 0)

	err := chromedp.Run(ctx,
		chromedp.Sleep(this.readyDelay*3),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().
				WithURLs(api.CookiesUrls).
				Do(ctx)
			if err != nil {
				return err
			}

			for _, c := range cookies {
				adapted := http.Cookie{
					Name:   c.Name,
					Value:  c.Value,
					Path:   c.Path,
					Domain: c.Domain,
					Secure: c.Secure,
				}
				ret = append(ret, &adapted)
			}

			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type AuthOption = func(auth *Auth)

func WithReadyDelay(d time.Duration) AuthOption {
	return func(auth *Auth) {
		auth.readyDelay = d
	}
}

// TODO: assert dk_access_token cookie
func NewAuth(tab context.Context, opts ...AuthOption) *Auth {
	auth := &Auth{
		tab: tab,

		readyDelay: time.Second * 3,
	}

	for _, opt := range opts {
		opt(auth)
	}

	return auth
}
