package main

import "github.com/gobwas/glob"

var (
	urlBlackListGlobCache []glob.Glob
	urlWhiteListGlobCache []glob.Glob
)

func initRuleGlobCache() error {
	for _, rule := range Config.Blacklist {
		g, err := glob.Compile(rule)
		if err != nil {
			return err
		}
		urlBlackListGlobCache = append(urlBlackListGlobCache, g)
	}

	for _, rule := range Config.Whitelist {
		g, err := glob.Compile(rule)
		if err != nil {
			return err
		}
		urlWhiteListGlobCache = append(urlWhiteListGlobCache, g)
	}

	return nil
}

func checkMatchList(url string, rules []glob.Glob) bool {
	for _, rule := range rules {
		if rule.Match(url) {
			return true
		}
	}
	return false
}

func urlMatchBlackList(url string) bool {
	return checkMatchList(url, urlBlackListGlobCache)
}

func urlMatchWhiteList(url string) bool {
	return checkMatchList(url, urlWhiteListGlobCache)
}
