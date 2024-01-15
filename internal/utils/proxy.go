package utils

import (
	"net/http/httputil"
	"net/url"
)

// proxyRequest takes in a targetUrl and proxies the entire request to that service
func ProxyRequest(targetUrl string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}
