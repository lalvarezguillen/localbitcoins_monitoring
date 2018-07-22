package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "currency",
			Usage: "The currency you want to sell your BTCs for.",
		},
	}
	app.Action = func(c *cli.Context) error {
		currency := c.String("currency")
		if currency == "" {
			panic("--currency is a required argument")
		}
		var keywords []string
		for _, kw := range c.Args() {
			keywords = append(keywords, strings.ToLower(kw))
		}
		fmt.Println("Currency: " + currency)
		fmt.Println(fmt.Sprintf("Keywords: %v", keywords))
		return gatherBuyers(currency, keywords)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type AdData struct {
	Message   string `json:"msg"`
	BankName  string `json:"bank_name"`
	Price     string `json:"temp_price"`
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
	Profile   struct {
		Username string `json:"username"`
	} `json:"profile"`
}

type BuyAd struct {
	Data    AdData `json:"data"`
	Actions struct {
		URL string `json:"public_view"`
	} `json:"actions"`
}

type Pagination struct {
	Next string `json:"next"`
}

type LBTCResponse struct {
	Data struct {
		Ads []BuyAd `json:"ad_list"`
	} `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func getPage(u string) (LBTCResponse, error) {
	resp, err := http.Get(u)
	if err != nil {
		return LBTCResponse{}, err
	}
	if resp.StatusCode != 200 {
		return LBTCResponse{}, fmt.Errorf(fmt.Sprintf("Status code %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LBTCResponse{}, err
	}

	var data LBTCResponse
	json.Unmarshal(body, &data)
	return data, nil
}

func containsKeywords(s string, kws []string) bool {
	for _, kw := range kws {
		if strings.Contains(strings.ToLower(s), strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func filterBuyers(ads []BuyAd, kws []string) []BuyAd {
	fb := []BuyAd{}
	for _, ad := range ads {
		msgContainsKws := containsKeywords(ad.Data.Message, kws)
		bankContainsKws := containsKeywords(ad.Data.BankName, kws)
		if msgContainsKws || bankContainsKws {
			fb = append(fb, ad)
		}
	}
	return fb
}

func gatherBuyers(c string, keywords []string) error {
	baseURL := fmt.Sprintf("https://localbitcoins.com/sell-bitcoins-online/%s/.json", c)
	buyers := []BuyAd{}
	var err error
	var page LBTCResponse
	page.Pagination.Next = baseURL
	for {
		page, err = getPage(page.Pagination.Next)
		if err != nil {
			return err
		}
		buyers = append(buyers, filterBuyers(page.Data.Ads, keywords)...)
		if page.Pagination.Next == "" {
			break
		}
	}

	jsonBuyers, err := json.MarshalIndent(buyers, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBuyers))
	return nil
}
