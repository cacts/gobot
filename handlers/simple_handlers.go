package handlers

import ( 
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
