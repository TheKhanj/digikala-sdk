package products

import (
	"testing"

	"github.com/thekhanj/digikala-sdk/api"
)

func TestFetchingComments(t *testing.T) {
	client, err := api.NewClientWithResponses(api.Server)
	if err != nil {
		t.Fatal(err)
		return
	}

	iphoneId := 20481189
	res, err := client.GetV1RateReviewProductsProductIdWithResponse(
		t.Context(), iphoneId, nil,
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	if res.JSON200.Status != 200 {
		t.Fatal("Expected 200 status code")
		return
	}

	comment := res.JSON200.Data.Comments[0]
	t.Log(comment)
}
