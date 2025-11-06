package browser

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/thekhanj/digikala-sdk/common"
)

func TestBrowser(t *testing.T) {
	// t.Skip()

	notify, cancel := signal.NotifyContext(
		t.Context(), os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	ctx, cancel2 := context.WithTimeout(notify, time.Second*3)
	defer cancel2()

	proxy := HttpProxy{
		common.GetMandatoryEnv("TEST_HTTP_PROXY_USER"),
		common.GetMandatoryEnv("TEST_HTTP_PROXY_PASSWORD"),
		common.GetMandatoryEnv("TEST_HTTP_PROXY_ADDRESS"),
	}
	b := NewBrowser(
		WithHeadless(false),
		WithHttpProxy(proxy),
	)
	b.Run(ctx)

	for i := 0; i < 3; i++ {
		go func() {
			tab, err := b.NewTab()
			if err != nil {
				t.Fatal(err)
				return
			}
			chromedp.Run(tab, chromedp.Navigate("https://ifconfig.io"))
		}()
	}

	<-ctx.Done()
}
