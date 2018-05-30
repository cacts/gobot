package modules

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type race struct {
	*discordgo.Session
	channelID    string
	participants []string
	raceMessage  *discordgo.Message
}

var currentRace *race

func HandleRaceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.TrimPrefix(m.Content, "!race ")

	if currentRace == nil {
		if command != "start" {
			s.ChannelMessageSend("theres no race right now idiot", m.ChannelID)
		} else {
			currentRace = &race{
				s,
				m.ChannelID,
				[]string{},
				nil,
			}
		}
	} else {
		switch command {
		case "start":
			s.ChannelMessageSend("theres already a race in progress idiot", m.ChannelID)
		case "enter":

		}
	}

	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
}

func (r *race) startRace(channelID string) {
	msg, _ := r.ChannelMessageSend(r.channelID, "a race is about to begin! type `!race enter` to enter!")
	r.raceMessage = msg
}

func (r *race) enterRace(user discordgo.User) {
	r.participants = append(r.participants, user.Username)
	r.raceMessage, _ = r.ChannelMessageEdit(r.channelID, r.raceMessage.ID, fmt.Sprintf("a race is about to begin! type `!race enter` to enter!\nracers: %v", strings.Join(r.participants, ", ")))
}
