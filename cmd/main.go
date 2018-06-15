package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

var timeoutChans = map[string]chan bool{}

func getTimeout(channelID string) bool {
	tc, ok := timeoutChans[channelID]
	if !ok {
		// TODO: this is almost certainly bad
		tc = make(chan bool, 1)
		timeoutChans[channelID] = tc
		timeoutChans[channelID] <- true
	}

	select {
	case t := <-tc:
		go func() {
			time.Sleep(2 * time.Minute)
			timeoutChans[channelID] <- true
		}()
		return t
	default:
		return false
	}
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	btcTicker := time.NewTicker(30 * time.Second)

	dg.UpdateStreamingStatus(0, btcPrice(), "")
	go func() {
		for range btcTicker.C {
			dg.UpdateStreamingStatus(0, btcPrice(), "")
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type messageHandler struct {
	matcher func(s string) bool
	exec    func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func getHandlers() []messageHandler {
	frog := messageHandler{
		func(s string) bool {
			return strings.Contains(s, "frogsiren")
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			// TODO: remove hardcoded emojis
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "a:frogsiren:396018906449313802")
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	b := messageHandler{
		func(s string) bool {
			return strings.Contains(s, " b ")
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "🅱")
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	lmao := messageHandler{
		func(s string) bool {
			return strings.Contains(s, "lmao")
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			id := m.ID

			prev, _ := s.ChannelMessages(m.ChannelID, 1, m.ID, "", "")

			// if the message is this exactly react to previous message
			if len(prev) > 0 && m.Content == "lmao" {
				id = prev[0].ID
			}

			if !getTimeout(m.ChannelID) {
				return
			}

			res := []rune{'l', 'm', 'a', 'o', '😂', '😹'}
			for _, c := range res {
				if c <= 'z' {
					c = rune(127462 + int(c) - 97)
				}
				err := s.MessageReactionAdd(m.ChannelID, id, fmt.Sprintf("%c", c))
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	}

	nice := messageHandler{func(s string) bool {
		return strings.Contains(s, "nice")
	},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if !getTimeout(m.ChannelID) {
				return
			}

			// TODO: remove hardcoded emojis
			err := s.MessageReactionAdd(m.ChannelID, m.ID, ":nice:395993706697719811")
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	raceHandler := messageHandler{func(s string) bool {
		return strings.Index(s, "!race") == 0
	},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
			modules.HandleRaceCommand(s, m)
		},
	}

	fuckboy := messageHandler{func(s string) bool {
		return strings.Index(s, "!mixy") == 0
	},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {

			res := []rune{'f', 'u', 'c', 'k', 'b', 'o', 'y'}
			for _, c := range res {
				if c <= 'z' {
					c = rune(127462 + int(c) - 97)
				}
				err := s.MessageReactionAdd(m.ChannelID, m.ID, fmt.Sprintf("%c", c))
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	}

	return []messageHandler{frog, b, lmao, nice, raceHandler, fuckboy}
}

var handlers []messageHandler

func init() {
	handlers = getHandlers()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "<@447429502390370315>") {
		split := strings.SplitN(strings.TrimSpace(m.Content), " ", 2)
		message := ""
		if len(split) > 1 {
			message = split[1]
		}
		println(message)
		if strings.HasPrefix(message, "<") && strings.HasSuffix(message, ">") {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%v> %v", m.Author.ID, message))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%v> <:fingerguns:342698356818182156>", m.Author.ID))
		}

		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			println(err)
		}
	}

	// debug messages
	fmt.Printf("%v: %v (embeds %v)\n", m.ChannelID, m.Content, m.Embeds)

	for _, handler := range handlers {
		if handler.matcher(m.Content) {
			fmt.Println("matched")
			handler.exec(s, m)
			return
		}
	}
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

func btcPrice() string {
	url := "https://api.coindesk.com/v1/bpi/currentprice.json"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		println(err)
	}
	println(res)
	response := coinbaseReponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println(err)
	}
	println(body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		println(err)
	}

	return fmt.Sprintf("BTC $%.2f", response.Bpi.USD.Price)
}
