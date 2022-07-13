package gobot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	BOT_ID = "<@447429502390370315>"

	DEFAULT_CHANNEL_ID = "921044490360348736"
)

var Global = &Gobot{}

type Gobot struct {
	handlers []MessageHandler
	tickers  []TickHandler
	ticker   *time.Ticker
	stopCh   chan struct{}
	discord  *discordgo.Session
	sync.Mutex
}

func (g *Gobot) Open(token string) error {
	g.stopCh = make(chan struct{}, 1)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return err
	}

	dg.AddHandler(g.onMessageCreate)
	g.discord = dg

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	g.startTicking()

	return nil
}

func (g *Gobot) SendMessage(message string) {
	g.discord.ChannelMessageSend(DEFAULT_CHANNEL_ID, message)
}

func (g *Gobot) Session() *discordgo.Session {
	return g.discord
}

func (g *Gobot) startTicking() {
	g.ticker = time.NewTicker(1 * time.Second)

	// do all initial ticks
	for _, t := range g.tickers {
		t.Tick(g.discord)
	}

	go func() {
		defer close(g.stopCh)
		defer g.ticker.Stop()

		for {
			select {
			case <-g.stopCh:
				return
			case <-g.ticker.C:
				g.Lock()
				tickers := g.tickers
				g.Unlock()

				for _, t := range tickers {
					if time.Now().Second()%t.Interval() == 0 {
						go t.Tick(g.discord)
					}
				}
			}
		}
	}()
}

func (g *Gobot) Close() {
	err := g.discord.Close()
	if err != nil {
		fmt.Println(err)
	}

	g.stopCh <- struct{}{}
}

type MessageHandler interface {
	Name() string
	Handle(s *discordgo.Session, m *discordgo.MessageCreate)
}

type TickHandler interface {
	Name() string
	Interval() int // seconds
	Tick(*discordgo.Session)
}

func (g *Gobot) AddMessageHandler(h MessageHandler) {
	g.Lock()
	defer g.Unlock()

	g.handlers = append(g.handlers, h)
}

func (g *Gobot) RemoveMessageHandler(name string) {
	g.Lock()
	defer g.Unlock()

	i := -1
	for idx, h := range g.handlers {
		if h.Name() == name {
			i = idx
			break
		}
	}

	if i >= 0 {
		l := len(g.handlers)
		g.handlers[i] = g.handlers[l-1]
		g.handlers[l-1] = nil
		g.handlers = g.handlers[:l-1]
	}
}

func (g *Gobot) AddTickHandler(t TickHandler) {
	g.Lock()
	defer g.Unlock()

	g.tickers = append(g.tickers, t)
}

func (g *Gobot) RemoveTickHandler(name string) {
	g.Lock()
	defer g.Unlock()

	i := -1
	for idx, t := range g.tickers {
		if t.Name() == name {
			i = idx
			break
		}
	}

	if i >= 0 {
		l := len(g.tickers)
		g.tickers[i] = g.tickers[l-1]
		g.tickers[l-1] = nil
		g.tickers = g.tickers[:l-1]
	}
}

func (g *Gobot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// debug messages
	fmt.Printf("%s|%s:%s#%s> %s (embeds %v)\n", m.ChannelID, m.Author.ID, m.Author.Username, m.Author.Discriminator, m.Content, m.Embeds)

	for _, handler := range g.handlers {
		go handler.Handle(s, m)
	}
}
