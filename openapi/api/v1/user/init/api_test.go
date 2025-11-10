package init

import (
	"log"
	"testing"

	"github.com/thekhanj/digikala-sdk/api"
	"github.com/thekhanj/digikala-sdk/api/http"
	"github.com/thekhanj/digikala-sdk/test"
)

func TestInit(t *testing.T) {
	env, err := test.NewTestEnv()
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

	httpClient, err := http.NewClient(
		http.WithCookies(cookies),
		http.WithCurlLog(log.Default()),
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	client, err := api.NewClientWithResponses(api.Server, api.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
		return
	}

	res, err := client.GetV1UserInitWithResponse(t.Context())
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(string(res.Body))

	if res.JSON200.Status != 200 {
		t.Fatal("Expected 200 status code")
		return
	}

	t.Log(res.JSON200.Data)
}
