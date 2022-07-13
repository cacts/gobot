package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

func init() {
	gobot.Global.AddMessageHandler(NewPrefixHandler("!race", raceHandler))
}

type participant struct {
	*discordgo.User
	progress float64
	dead     bool
}

func (p *participant) String() string {
	return p.Username
}

type race struct {
	*discordgo.Session
	channelID         string
	participants      []*participant
	raceMessage       *discordgo.Message
	secondsUntilStart int
}

var currentRace *race

func setupNewRace(s *discordgo.Session, channelID string) {
	currentRace = &race{
		s,
		channelID,
		[]*participant{},
		nil,
		60,
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			currentRace.secondsUntilStart--

			if currentRace.secondsUntilStart == 5 {
				currentRace.raceMessage, _ = currentRace.ChannelMessageSend(currentRace.channelID, "the race is about to begin!")
			}

			if currentRace.secondsUntilStart <= 0 {
				ticker.Stop()
				currentRace.startRace()
			}
		}
	}()

	currentRace.prepare()
}

func raceHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.TrimPrefix(m.Content, "!race ")
	println(command)
	if currentRace == nil {
		if command != "start" {
			s.ChannelMessageSend(m.ChannelID, "theres no race right now type !race start to start one")
		} else {
			setupNewRace(s, m.ChannelID)
		}
	} else {
		switch command {
		case "start":
			s.ChannelMessageSend(m.ChannelID, "theres already a race in progress~")
		case "enter":
			if currentRace.secondsUntilStart > 0 {
				currentRace.enterUser(m.Author)
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v the race already started", m.Author.Username))
			}
		case "detonate":
			if m.Author.ID == "126363515438104576" {
				name := strings.TrimPrefix(m.Content, "!race detonate ")
				for _, p := range currentRace.participants {
					if p.User.Username == name {
						p.dead = true
					}
				}
			}
		}
	}

	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
}

func (r *race) prepare() {
	msg, _ := r.ChannelMessageSend(r.channelID, "a race will begin in 60 seconds! type `!race enter` to enter!")
	r.raceMessage = msg
}

func (r *race) hasParticipant(p *participant) bool {
	for _, px := range r.participants {
		if px.ID == p.ID {
			return true
		}
	}
	return false
}

func (r *race) enterUser(user *discordgo.User) {
	newParticipant := &participant{user, 0.0, false}

	if r.hasParticipant(newParticipant) {
		r.ChannelMessageSend(r.channelID, fmt.Sprintf("%v is already in the race!", newParticipant))
		return
	}

	r.participants = append(r.participants, newParticipant)

	strs := []string{}
	for _, p := range r.participants {
		strs = append(strs, p.String())
	}

	r.raceMessage, _ = r.ChannelMessageEdit(r.channelID, r.raceMessage.ID, fmt.Sprintf("a race is about to begin! type `!race enter` to enter!\nracers: %v", strings.Join(strs, ", ")))
}

func (r *race) updateRaceInProgress() {
	message := "and they're off!\n```"
	for _, p := range r.participants {
		progress := math.Min(p.progress*50, 50)
		var icon string
		if p.dead {
			icon = "💥"
		} else {
			icon = "🚗"
		}

		message += fmt.Sprintf("%-20s 🏁%50s🚦\n", p, icon+strings.Repeat("~", int(math.Max(progress-1, 0))))
	}
	message += "```"
	r.raceMessage, _ = r.ChannelMessageEdit(r.channelID, r.raceMessage.ID, message)
}

func (r *race) startRace() {
	r.updateRaceInProgress()

	raceTicker := time.NewTicker(1500 * time.Millisecond)
	for range raceTicker.C {
		for _, p := range r.participants {
			if !p.dead {
				p.progress += math.Abs(.8-rand.Float64()) * .08
			}

			if p.progress > 1 {
				r.updateRaceInProgress()
				raceTicker.Stop()
				r.endRace(p)
				return
			}
		}
		r.updateRaceInProgress()
	}
}

func (r *race) endRace(winner *participant) {
	r.ChannelMessageSend(r.channelID, fmt.Sprintf("the race is over!!! %v is the winner", winner))
	currentRace = nil
}
