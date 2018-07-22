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
		m := NewMonitor()
		_, err := m.GatherBuyers(currency, keywords)
		return err
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Monitor contains methods to filter Localbitcoins' ads based on Currency and Keywords
type Monitor struct {
	HTTPClient *http.Client
}

// getPage fetches a page of ads from localbitcoin's API.
func (m *Monitor) getPage(u string) (LBTCResponse, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return LBTCResponse{}, err
	}
	resp, err := m.HTTPClient.Do(req)
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

// GatherBuyers filters the available buyers, leaving the ones that:
//    a) Match the currency we're interested in
//    b) Contain some of the keywords we;re looking for, either in their description or bank name
func (m *Monitor) GatherBuyers(c string, keywords []string) ([]Ad, error) {
	baseURL := fmt.Sprintf("https://localbitcoins.com/sell-bitcoins-online/%s/.json", c)
	buyers := []Ad{}
	var err error
	var page LBTCResponse
	page.Pagination.Next = baseURL
	for {
		page, err = m.getPage(page.Pagination.Next)
		if err != nil {
			return []Ad{}, err
		}
		buyers = append(buyers, filterBuyers(page.Data.Ads, keywords)...)
		if page.Pagination.Next == "" {
			break
		}
	}

	jsonBuyers, err := json.MarshalIndent(buyers, "", "    ")
	if err != nil {
		return []Ad{}, err
	}
	fmt.Println(fmt.Sprintf("Buyers Found: %d", len(buyers)))
	fmt.Println(string(jsonBuyers))
	return buyers, err
}

// NewMonitor creates a new Monitor
func NewMonitor() Monitor {
	return Monitor{HTTPClient: &http.Client{}}
}
