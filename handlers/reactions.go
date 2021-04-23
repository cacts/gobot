package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

const ()

func init() {
	gobot.Global.AddMessageHandler(NewContainsHandler("lmao", lmaoHandler))
	gobot.Global.AddMessageHandler(NewContainsHandler(" b ", bHandler))
	gobot.Global.AddMessageHandler(NewPrefixHandler("F", fHandler))
}

func fHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != "F" {
		return
	}

	c := rune(127462 + int('f') - 97)

	err := s.MessageReactionAdd(m.ChannelID, m.ID, fmt.Sprintf("%c", c))
	if err != nil {
		fmt.Println(err)
	}
}

func lmaoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	id := m.ID

	prev, _ := s.ChannelMessages(m.ChannelID, 1, m.ID, "", "")

	// if the message is this exactly react to previous message
	if len(prev) > 0 && m.Content == "lmao" {
		id = prev[0].ID
	}

	if !gobot.GlobalTimeout(m.ChannelID) {
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
}

func bHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "🅱")
	if err != nil {
		fmt.Println(err)
	}
}
