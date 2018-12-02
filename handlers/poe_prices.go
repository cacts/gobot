package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cactauz/gobot" 
	"github.com/bwmarrin/discordgo"
)

func init() {
	gobot.Global.AddMessageHandler(NewPrefixHandler("poepc ", poePriceHandler))
}
 
const (
	LEAGUE = "Delve"
)

var (
	baseUrl = fmt.Sprintf("https://www.poeprices.info/api?l=%s&i=", LEAGUE)
	client  = http.Client{Timeout: time.Second * 5}
)

type poePriceResponse struct {
	Status   int     `json:"status"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

func poePriceHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	item := strings.TrimPrefix(m.Content, "poepc ")
	fmt.Println(item)
	encoded := base64.URLEncoding.EncodeToString([]byte(item))
	req, err := http.NewRequest(http.MethodGet, baseUrl+encoded, nil)
	if err != nil { 
		fmt.Println("err on newrequest", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("err getting response", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err reading body", err)
		return 
	}
	price := poePriceResponse{}
	err = json.Unmarshal(body, &price)
	if err != nil {
		fmt.Println("err unmarshalling json", err, string(body))
		return
	}

	if price.Status != 200 {
		fmt.Println("err status", body)
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ML price estimate %.1f - %.1f %s", price.Min, price.Max, price.Currency))
	if err != nil {
		fmt.Println("err sending to channel", err)
	}
}
