package config

import (
	"fmt"
	"net/http"
	"testing"
)

func testCustomerWriteCookies(t *testing.T, cookies []*http.Cookie) {
	var c Config
	err := ReadConfig(&c, GetConfigPath())
	if err != nil {
		t.Fatal(err)
		return
	}

	err = c.Customer.Cache.SetCookies(cookies)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = WriteConfig(&c, GetConfigPath())
	if err != nil {
		t.Fatal(err)
		return
	}
}

func testCustomerReadCookies(t *testing.T, cookies []*http.Cookie) {
	var c Config
	err := ReadConfig(&c, GetConfigPath())
	if err != nil {
		t.Fatal(err)
		return
	}

	savedCookies, err := c.Customer.Cache.GetCookies()
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, cookie := range cookies {
		found := false
		for _, savedCookie := range savedCookies {
			isEqual := cookie.Name == savedCookie.Name &&
				cookie.Value == savedCookie.Value
			if isEqual {
				found = true
				break
			}
		}

		if !found {
			t.Fatal(fmt.Errorf("Cookie %s was not found", cookie.Name))
		}
		t.Logf("Cookie %s found", cookie.Name)
	}
}

func TestCustomerCookies(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "NAME", Value: "VALUE"},
	}
	testCustomerWriteCookies(t, cookies)
	testCustomerReadCookies(t, cookies)
}
