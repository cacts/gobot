package handlers

import (
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

const (
	GUILD_ID = "396018642527059973"
	ROLE_ID  = "835145640852848640"
)

func init() {
	gobot.Global.AddTickHandler(NewSimpleTickHandler("rainbow", 1, rainbowTickHandler))
}

var steps = 128

func sinToHex(i int, phase float64) int {
	sin := math.Sin(math.Pi/float64(steps)*2.0*float64(i) + phase)
	asInt := math.Floor(sin*127) + 128
	return int(asInt)
}

func rainbow() func() int {
	idx := 0

	return func() int {
		idx = (idx + 1) % steps

		r := sinToHex(idx, 0*math.Pi*2.0/3.0)
		g := sinToHex(idx, 1*math.Pi*2.0/3.0)
		b := sinToHex(idx, 2*math.Pi*2.0/3.0)

		return (r << 16) | (g << 8) | (b)
	}
}

var rainbowGen = rainbow()
var rainbowRole *discordgo.Role

func rainbowTickHandler(s *discordgo.Session) {
	if rainbowRole != nil {
		var err error
		rainbowRole, err = s.GuildRoleEdit(GUILD_ID, ROLE_ID, rainbowRole.Name, rainbowGen(), rainbowRole.Hoist, rainbowRole.Permissions, rainbowRole.Mentionable)
		if err != nil {
			fmt.Printf("error changing role color: %v\n", err)
		}
		return
	}

	roles, err := s.GuildRoles(GUILD_ID)
	if err != nil {
		fmt.Printf("error retrieving rainbow role: %v", err)
		gobot.Global.RemoveTickHandler("rainbow")
	}

	for _, r := range roles {
		if r.ID == ROLE_ID {
			rainbowRole = r
			break
		}
	}

	if rainbowRole == nil {
		fmt.Printf("role %s doesnt exist in guild %s\n", ROLE_ID, GUILD_ID)
		gobot.Global.RemoveTickHandler("rainbow")
	}
}
