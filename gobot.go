package gobot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	BOT_ID = "<@447429502390370315>"
	// TODO: look these up from the api
	FINGER_GUNS = ":fingerguns:342698356818182156"
	NICE        = ":nice:395993706697719811"
	FROGSIREN   = "a:frogsiren:396018906449313802"
)

var Global = &Gobot{}

type Gobot struct {
	handlers []MessageHandler
	tickers  []TickHandler
	ticker   *time.Ticker
	stopCh   chan struct{}
	discord  *discordgo.Session
}

func (g *Gobot) Open(token string) error {
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

func (g *Gobot) startTicking() {
	g.stopCh = make(chan struct{}, 1)
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
				for _, t := range g.tickers {
					if time.Now().Second()%t.Interval() == 0 {
						t.Tick(g.discord)
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

type SimpleHandler struct {
	matchStr string
	matchFn  func(string) bool
	handler  func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func NewPrefixHandler(prefix string, handler func(s *discordgo.Session, m *discordgo.MessageCreate)) *SimpleHandler {
	return &SimpleHandler{
		matchStr: prefix,
		matchFn:  func(s string) bool { return strings.HasPrefix(s, prefix) },
		handler:  handler,
	}
}

func NewContainsHandler(contains string, handler func(s *discordgo.Session, m *discordgo.MessageCreate)) *SimpleHandler {
	return &SimpleHandler{
		matchStr: contains,
		matchFn:  func(s string) bool { return strings.Contains(s, contains) },
		handler:  handler,
	}
}

func (sh *SimpleHandler) Name() string {
	return sh.matchStr
}

func (sh *SimpleHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if sh.matchFn(m.Content) {
		str := m.Content
		if len([]rune(str)) > 40 {
			str = string([]rune(str)[:40])
		}
		fmt.Println(sh.matchStr, "matched", str)
		sh.handler(s, m)
	}
}

type TickHandler interface {
	Name() string
	Interval() int // seconds
	Tick(*discordgo.Session)
}

type SimpleTickHandler struct {
	name     string
	interval int
	tickFn   func(*discordgo.Session)
}

func NewSimpleTickHandler(name string, interval int, tickFn func(*discordgo.Session)) *SimpleTickHandler {
	return &SimpleTickHandler{
		name:     name,
		interval: interval,
		tickFn:   tickFn,
	}
}

func (th *SimpleTickHandler) Name() string {
	return th.name
}

func (th *SimpleTickHandler) Interval() int {
	return th.interval
}

func (th *SimpleTickHandler) Tick(s *discordgo.Session) {
	th.tickFn(s)
}

func (g *Gobot) AddMessageHandler(h MessageHandler) {
	g.handlers = append(g.handlers, h)
}

func (g *Gobot) AddTickHandler(t TickHandler) {
	g.tickers = append(g.tickers, t)
}

func (g *Gobot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// debug messages
	fmt.Printf("%v: %v (embeds %v)\n", m.ChannelID, m.Content, m.Embeds)

	for _, handler := range g.handlers {
		go handler.Handle(s, m)
	}
}
