package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	//gobot.Global.AddTickHandler(gobot.NewSimpleTickHandler("btc_price", 30, btcPriceTicker))
}

type bpi struct {
	USD USD `json:"USD"`
}

type USD struct {
	Price float32 `json:"rate_float"`
}

type coinbaseReponse struct {
	Bpi bpi `json:"bpi"`
}

// blah blah Powered by Coinbase https://www.coindesk.com/price/
func btcPriceTicker(s *discordgo.Session) {
	url := "https://api.coindesk.com/v1/bpi/currentprice.json"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := coinbaseReponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	s.UpdateGameStatus(0, fmt.Sprintf("BTC $%.2f", response.Bpi.USD.Price))
}
