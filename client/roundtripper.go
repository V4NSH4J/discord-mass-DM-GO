package client

import (
	"context"
	"errors"
	"fmt"
	"net"

	"strings"
	"sync"

	http "github.com/Danny-Dasilva/fhttp"
	http2 "github.com/Danny-Dasilva/fhttp/http2"
	utls "github.com/Danny-Dasilva/utls"
	"golang.org/x/net/proxy"
)

var errProtocolNegotiated = errors.New("protocol negotiated")

type roundTripper struct {
	sync.Mutex
	// fix typing
	JA3       string
	UserAgent string

	Cookies           []Cookie
	cachedConnections map[string]net.Conn
	cachedTransports  map[string]http.RoundTripper

	dialer proxy.ContextDialer
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Fix this later for proper cookie parsing
	for _, properties := range rt.Cookies {
		req.AddCookie(&http.Cookie{
			Name:       properties.Name,
			Value:      properties.Value,
			Path:       properties.Path,
			Domain:     properties.Domain,
			Expires:    properties.JSONExpires.Time, //TODO: scuffed af
			RawExpires: properties.RawExpires,
			MaxAge:     properties.MaxAge,
			HttpOnly:   properties.HTTPOnly,
			Secure:     properties.Secure,
			Raw:        properties.Raw,
			Unparsed:   properties.Unparsed,
		})
	}
	req.Header.Set("User-Agent", rt.UserAgent)
	addr := rt.getDialTLSAddr(req)
	if _, ok := rt.cachedTransports[addr]; !ok {
		if err := rt.getTransport(req, addr); err != nil {
			return nil, err
		}
	}
	return rt.cachedTransports[addr].RoundTrip(req)
}

func (rt *roundTripper) getTransport(req *http.Request, addr string) error {
	switch strings.ToLower(req.URL.Scheme) {
	case "http":
		rt.cachedTransports[addr] = &http.Transport{DialContext: rt.dialer.DialContext, DisableKeepAlives: true}
		return nil
	case "https":
	default:
		return fmt.Errorf("invalid URL scheme: [%v]", req.URL.Scheme)
	}

	_, err := rt.dialTLS(context.Background(), "tcp", addr)
	switch err {
	case errProtocolNegotiated:
	case nil:
		// Should never happen.
		panic("dialTLS returned no error when determining cachedTransports")
	default:
		return err
	}

	return nil
}

func (rt *roundTripper) dialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	rt.Lock()
	defer rt.Unlock()

	// If we have the connection from when we determined the HTTPS
	// cachedTransports to use, return that.
	if conn := rt.cachedConnections[addr]; conn != nil {
		delete(rt.cachedConnections, addr)
		return conn, nil
	}
	rawConn, err := rt.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	var host string
	if host, _, err = net.SplitHostPort(addr); err != nil {
		host = addr
	}
	//////////////////

	spec, err := StringToSpec(rt.JA3, rt.UserAgent)
	if err != nil {
		return nil, err
	}

	conn := utls.UClient(rawConn, &utls.Config{ServerName: host, InsecureSkipVerify: true}, // MinVersion:         tls.VersionTLS10,
		// MaxVersion:         tls.VersionTLS13,

		utls.HelloCustom)

	if err := conn.ApplyPreset(spec); err != nil {
		return nil, err
	}

	if err = conn.Handshake(); err != nil {
		_ = conn.Close()

		if err.Error() == "tls: CurvePreferences includes unsupported curve" {
			//fix this
			return nil, fmt.Errorf("conn.Handshake() error for tls 1.3 (please retry request): %+v", err)
		}
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	//////////
	if rt.cachedTransports[addr] != nil {
		return conn, nil
	}

	// No http.Transport constructed yet, create one based on the results
	// of ALPN.
	switch conn.ConnectionState().NegotiatedProtocol {
	case http2.NextProtoTLS:
		parsedUserAgent := parseUserAgent(rt.UserAgent)
		t2 := http2.Transport{DialTLS: rt.dialTLSHTTP2,
			PushHandler: &http2.DefaultPushHandler{},
			Navigator:   parsedUserAgent,
		}
		rt.cachedTransports[addr] = &t2
	default:
		// Assume the remote peer is speaking HTTP 1.x + TLS.
		rt.cachedTransports[addr] = &http.Transport{DialTLSContext: rt.dialTLS}

	}

	// Stash the connection just established for use servicing the
	// actual request (should be near-immediate).
	rt.cachedConnections[addr] = conn

	return nil, errProtocolNegotiated
}

func (rt *roundTripper) dialTLSHTTP2(network, addr string, _ *utls.Config) (net.Conn, error) {
	return rt.dialTLS(context.Background(), network, addr)
}

func (rt *roundTripper) getDialTLSAddr(req *http.Request) string {
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err == nil {
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort(req.URL.Host, "443") // we can assume port is 443 at this point
}

func newRoundTripper(browser Browser, dialer ...proxy.ContextDialer) http.RoundTripper {
	if len(dialer) > 0 {

		return &roundTripper{
			dialer: dialer[0],

			JA3:               browser.JA3,
			UserAgent:         browser.UserAgent,
			Cookies:           browser.Cookies,
			cachedTransports:  make(map[string]http.RoundTripper),
			cachedConnections: make(map[string]net.Conn),
		}
	}

	return &roundTripper{
		dialer: proxy.Direct,

		JA3:               browser.JA3,
		UserAgent:         browser.UserAgent,
		Cookies:           browser.Cookies,
		cachedTransports:  make(map[string]http.RoundTripper),
		cachedConnections: make(map[string]net.Conn),
	}
}
