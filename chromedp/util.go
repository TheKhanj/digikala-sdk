package chromedp

import (
	"context"
	"os/exec"

	"github.com/chromedp/chromedp"
)

func NewTab(ctx context.Context) (context.Context, context.CancelFunc) {
	allocCtx, c1 := chromedp.NewExecAllocator(ctx, SensibleExecPath())

	tab, c2 := chromedp.NewContext(allocCtx)

	c := func() {
		c2()
		c1()
	}

	return tab, c
}

func SensibleExecPath() chromedp.ExecAllocatorOption {
	locations := []string{
		"chromium",
		"chromium-browser",
		"brave",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
	}

	for _, path := range locations {
		found, err := exec.LookPath(path)
		if err == nil {
			return chromedp.ExecPath(found)
		}
	}

	return chromedp.ExecPath("google-chrome")
}
