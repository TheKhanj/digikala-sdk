package auth

import (
	"os"
	"testing"

	"github.com/thekhanj/digikala-sdk/chromedp"
)

func TestAuth(t *testing.T) {
	tab, cancel := chromedp.NewTab(t.Context())
	defer cancel()

	auth := NewAuth(tab)

	username := os.Getenv("TEST_CUSTOMER_USERNAME")
	password := os.Getenv("TEST_CUSTOMER_PASSWORD")
	if username == "" {
		t.Fatal("Empty username was passed to the test")
		return
	}
	if password == "" {
		t.Fatal("Empty password was passed to the test")
		return
	}
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
}
