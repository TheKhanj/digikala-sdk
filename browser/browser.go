package browser

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/thekhanj/digikala-sdk/api/http"
)

var ErrInShutdown = errors.New("browser: In shutdown")

type Proxy interface{}

type BrowserOption = func(o *Browser)

func WithHeadless(headless bool) BrowserOption {
	return func(o *Browser) {
		o.headless = headless
	}
}

func WithHttpProxy(proxy http.HttpProxy) BrowserOption {
	return func(o *Browser) {
		o.proxy = proxy
	}
}

type Browser struct {
	browserCtx context.Context
	cancels    []context.CancelFunc
	inShutdown atomic.Bool

	headless bool
	proxy    Proxy
}

func (this *Browser) Run(ctx context.Context) {
	this.inShutdown.Store(false)

	allocOpts := this.getAllocOptions()
	allocCtx, c1 := chromedp.NewExecAllocator(ctx, allocOpts...)
	this.appendCancel(c1)

	browserCtx, c2 := chromedp.NewContext(allocCtx)
	this.browserCtx = browserCtx
	this.appendCancel(c2)

	go func() {
		<-ctx.Done()
		this.shutdown()
	}()
}

func (this *Browser) getAllocOptions() []chromedp.ExecAllocatorOption {
	allocOpts := []chromedp.ExecAllocatorOption{
		SensibleExecPath(),

		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,

		// After Puppeteer's default behavior.
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),

		chromedp.Flag("enable-automation", false),
		chromedp.Flag("window-size", "1080,524"),
	}
	if httpProxy, ok := this.proxy.(http.HttpProxy); ok {
		allocOpts = append(allocOpts, chromedp.ProxyServer(httpProxy.Address))
	}
	if this.headless {
		allocOpts = append(allocOpts, chromedp.Headless)
	}

	return allocOpts
}

func (this *Browser) shutdown() error {
	if this.inShutdown.Load() {
		return ErrInShutdown
	}

	this.inShutdown.Store(true)

	for _, cancel := range this.cancels {
		cancel()
	}
	this.cancels = this.cancels[:0]

	return nil
}

func (this *Browser) NewTab() (context.Context, error) {
	if this.inShutdown.Load() {
		return nil, ErrInShutdown
	}

	tab, cancel := chromedp.NewContext(this.browserCtx)
	_ = this.appendCancel(cancel)

	if httpProxy, ok := this.proxy.(http.HttpProxy); ok {
		this.addProxyAuth(tab, httpProxy.Username, httpProxy.Password)
		err := chromedp.Run(tab, fetch.Enable().WithHandleAuthRequests(true))
		if err != nil {
			return nil, err
		}
		return tab, nil
	} else {
		err := chromedp.Run(tab)
		if err != nil {
			return nil, err
		}
		return tab, nil
	}
}

func (this *Browser) appendCancel(c context.CancelFunc) error {
	if this.inShutdown.Load() {
		return ErrInShutdown
	}

	this.cancels = append(this.cancels, c)
	return nil
}

func (this *Browser) addProxyAuth(ctx context.Context, user, password string) {
	lctx, lcancel := context.WithCancel(ctx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
			}()
		case *fetch.EventAuthRequired:
			if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				go func() {
					_ = chromedp.Run(ctx,
						fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
							Response: fetch.AuthChallengeResponseResponseProvideCredentials,
							Username: user,
							Password: password,
						}),
						fetch.Disable(),
					)
					lcancel()
				}()
			}
		}
	})
}

func NewBrowser(opts ...BrowserOption) *Browser {
	b := &Browser{
		browserCtx: nil,
		cancels:    make([]context.CancelFunc, 0),
		inShutdown: atomic.Bool{},

		headless: true,
		proxy:    nil,
	}
	b.inShutdown.Store(true)

	for _, o := range opts {
		o(b)
	}

	return b
}
