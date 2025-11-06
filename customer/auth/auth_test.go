package auth

import (
	"errors"
	"testing"

	"github.com/thekhanj/digikala-sdk/browser"
	"github.com/thekhanj/digikala-sdk/common"
)

func TestAuth(t *testing.T) {
	proxy := browser.HttpProxy{
		User:     common.GetMandatoryEnv("TEST_HTTP_PROXY_USER"),
		Password: common.GetMandatoryEnv("TEST_HTTP_PROXY_PASSWORD"),
		Address:  common.GetMandatoryEnv("TEST_HTTP_PROXY_ADDRESS"),
	}
	b := browser.NewBrowser(
		browser.WithHeadless(false),
		browser.WithHttpProxy(proxy),
	)
	b.Run(t.Context())

	tab, err := b.NewTab()
	if err != nil {
		t.Fatal(err)
		return
	}

	auth := NewAuth(tab)

	username := common.GetMandatoryEnv("TEST_CUSTOMER_USERNAME")
	password := common.GetMandatoryEnv("TEST_CUSTOMER_PASSWORD")
	cookies, err := auth.Login(username, password)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(cookies) == 0 {
		t.Fatal("No cookies found")
		return
	}

	t.Log(cookies)

	found := false
	for _, cookie := range cookies {
		if cookie.Name == "DK_ACCESS_TOKEN" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal(errors.New("DK_ACCESS_TOKEN cookie not found!"))
		return
	}
}
