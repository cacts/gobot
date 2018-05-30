package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"github.com/cactauz/modules"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

var timeoutChan chan bool

func init() {
	timeoutChan = make(chan bool, 1)
	timeoutChan <- true
}

func getTimeout() bool {
	select {
	case t := <-timeoutChan:
		go func() {
			time.Sleep(2 * time.Minute)
			timeoutChan <- true
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
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type messageHandler struct {
	matcher func(s string) bool
	// return true to trigger global spam timeout
	exec func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func getHandlers() []messageHandler {
	frog := messageHandler{
		func(s string) bool {
			return strings.Contains(s, "frogsiren")
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

			if !getTimeout() {
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
			if !getTimeout() {
				return
			}

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

	return []messageHandler{frog, b, lmao, nice}
}

var handlers []messageHandler

func init() {
	handlers = getHandlers()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || m.Author.ID == "124571897391349760" {
		return
	}

	fmt.Printf("%v (embeds %v)\n", m.Content, m.Embeds)

	for _, handler := range handlers {
		if handler.matcher(m.Content) {
			fmt.Println("matched")
			handler.exec(s, m)
			return
		}
	}

	if strings.Contains(m.Content, "<@126363515438104576>") {
		s.ChannelMessageSend(m.ChannelID, "no paging <@126363515438104576>")
	}

	if m.Content == "bot emergency mode" && m.Author.ID == "126363515438104576" {
		msgs, _ := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
		for _, m := range msgs {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "a:canyounot:397537060828872715")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
