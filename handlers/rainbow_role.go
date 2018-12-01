package handlers

import (
	"math"

	"../../gobot"
	"github.com/bwmarrin/discordgo"
)

const (
	GUILD_ID = "124572142485504002"
	ROLE_ID  = "438028062303453205"
)

func init() {
	gobot.Global.AddTickHandler(gobot.NewSimpleTickHandler("rainbow", 1, rainbowTickHandler))
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

func rainbowTickHandler(s *discordgo.Session) {
	roles, _ := s.GuildRoles(GUILD_ID)
	var rRole *discordgo.Role
	for _, r := range roles {
		if r.ID == ROLE_ID {
			rRole = r
		}
	}

	if rRole != nil {
		s.GuildRoleEdit(GUILD_ID, ROLE_ID, rRole.Name, rainbowGen(), rRole.Hoist, rRole.Permissions, rRole.Mentionable)
	}
}
