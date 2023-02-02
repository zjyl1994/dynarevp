package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/samber/lo"
)

type reverseProxy struct{}

func (reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check method if set allow method list
	if len(Config.Method) > 0 {
		if !lo.Contains(Config.Method, r.Method) {
			errorPage(w, http.StatusMethodNotAllowed, "only allow "+strings.Join(Config.Method, "/"))
			return
		}
	}
	// Get destination url
	proxyDestUrl := r.URL.RequestURI()
	if Config.Prefix > "" {
		if strings.HasPrefix(proxyDestUrl, Config.Prefix) {
			proxyDestUrl = strings.TrimPrefix(proxyDestUrl, Config.Prefix)
		} else {
			errorPage(w, http.StatusNotFound, proxyDestUrl)
			return
		}
	}
	proxyDestUrl = strings.TrimPrefix(proxyDestUrl, "/")

	// Proxy to destination
	destUrl, err := url.Parse(proxyDestUrl)
	if err != nil {
		errorPage(w, http.StatusBadRequest, err.Error())
		return
	}

	if !(destUrl.Scheme == "http" || destUrl.Scheme == "https") {
		errorPage(w, http.StatusForbidden, "only support http/https")
		return
	}

	// Block url in blacklist
	if len(Config.Blacklist) > 0 {
		if urlMatchBlackList(proxyDestUrl) {
			errorPage(w, http.StatusForbidden, "url in blacklist")
			return
		}
	}
	// Block url not in whitelist
	if len(Config.Whitelist) > 0 {
		if !urlMatchWhiteList(proxyDestUrl) {
			errorPage(w, http.StatusForbidden, "url not in whitelist")
			return
		}
	}

	// Proxy to destination
	proxy := newReverseProxy(destUrl)
	proxy.ServeHTTP(w, r)
}

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL = target
		req.Host = target.Host
	}}
}

func errorPage(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, statusCode, http.StatusText(statusCode), msg)
}
