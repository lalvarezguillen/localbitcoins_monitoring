package main

// Ad represents an Localbitcoin's Ad buying BTC
type Ad struct {
	Data struct {
		Message   string `json:"msg"`
		BankName  string `json:"bank_name"`
		Price     string `json:"temp_price"`
		MinAmount string `json:"min_amount"`
		MaxAmount string `json:"max_amount"`
		Profile   struct {
			Username string `json:"username"`
		} `json:"profile"`
	} `json:"data"`
	Actions struct {
		URL string `json:"public_view"`
	} `json:"actions"`
}

// LBTCResponse represents a response from Localbitcoin's API
type LBTCResponse struct {
	Data struct {
		Ads []Ad `json:"ad_list"`
	} `json:"data"`
	Pagination struct {
		Next string `json:"next"`
	} `json:"pagination"`
}
