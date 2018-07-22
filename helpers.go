package main

import "strings"

func containsKeywords(s string, kws []string) bool {
	for _, kw := range kws {
		if strings.Contains(strings.ToLower(s), strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func filterBuyers(ads []Ad, kws []string) []Ad {
	fb := []Ad{}
	for _, ad := range ads {
		msgContainsKws := containsKeywords(ad.Data.Message, kws)
		bankContainsKws := containsKeywords(ad.Data.BankName, kws)
		if msgContainsKws || bankContainsKws {
			fb = append(fb, ad)
		}
	}
	return fb
}
