package browser

import (
	"os/exec"

	"github.com/chromedp/chromedp"
)

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
