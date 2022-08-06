package client

import (
	http "github.com/Danny-Dasilva/fhttp"

	"time"

	"golang.org/x/net/proxy"
)

type Browser struct {
	// Return a greeting that embeds the name in a message.
	JA3       string
	UserAgent string
	Cookies   []Cookie
}

var disabledRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func clientBuilder(browser Browser, dialer proxy.ContextDialer, timeout int, disableRedirect bool) *http.Client {
	//if timeout is not set in call default to 15
	if timeout == 0 {
		timeout = 15
	}
	client := &http.Client{
		Transport: newRoundTripper(browser, dialer),
		Timeout:   time.Duration(timeout) * time.Second,
	}
	//if disableRedirect is set to true httpclient will not redirect
	if disableRedirect {
		client.CheckRedirect = disabledRedirect
	}
	return client
}

// newClient creates a new http client
func NewClient(browser Browser, timeout int, disableRedirect bool, UserAgent string, proxyURL ...string) (*http.Client, error) {
	//fix check PR
	if len(proxyURL) > 0 && len(proxyURL[0]) > 0 {
		dialer, err := newConnectDialer(proxyURL[0], UserAgent)
		if err != nil {
			return &http.Client{
				Timeout:       time.Duration(timeout) * time.Second,
				CheckRedirect: disabledRedirect, //fix this fallthrough issue (test for incorrect proxy)
			}, err
		}
		return clientBuilder(
			browser,
			dialer,
			timeout,
			disableRedirect,
		), nil
	}

	return clientBuilder(
		browser,
		proxy.Direct,
		timeout,
		disableRedirect,
	), nil

}
