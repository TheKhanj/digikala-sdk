package test

import (
	"errors"
	"testing"
)

func TestCustomerAssertAuthorized(t *testing.T) {
	env, err := NewTestEnv()
	if err != nil {
		t.Fatal(err)
		return
	}

	err = env.CustomerAssertAuthorized(t.Context())
	if err != nil {
		t.Fatal(err)
		return
	}

	cookies, err := env.Config.Customer.Cache.GetCookies()
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
