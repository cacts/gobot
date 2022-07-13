package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

func init() {
	gobot.Global.AddMessageHandler(NewPrefixHandler("!color", colorHandler))
}

var colorRegex = regexp.MustCompile("^#[0-9A-Fa-f]{6}$")

func colorHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	colorHex := strings.TrimPrefix(m.Content, "!color ")
	fmt.Println(colorHex)
	if !colorRegex.MatchString(colorHex) {
		s.ChannelMessageSend(m.ChannelID, "pls give valid color like `!color #6495ed`")
		return
	}

	roles, err := getGuildRolesByID(s, m.GuildID)
	if err != nil {
		fmt.Println(err)
		return
	}

	colorRole, err := getAuthorColorRole(s, m, roles)
	if err != nil {
		fmt.Println(err)
		return
	}

	if colorRole == nil {
		err = createColorRole(s, m.GuildID, m.Author.ID, colorHex)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	roleUserIDs, err := getRoleUserIDs(s, m.GuildID, colorRole.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(roleUserIDs) == 1 && roleUserIDs[0] == m.Author.ID {
		colorRole.Name = colorHex
		colorRole.Color = hexColorToInt(colorHex[1:])
		_, err = s.GuildRoleEdit(m.GuildID, colorRole.ID, colorRole.Name, colorRole.Color, colorRole.Hoist, colorRole.Permissions, false)
	} else {
		err = createColorRole(s, m.GuildID, m.Author.ID, colorHex)

	}

	if err != nil {
		fmt.Println(err)
	}
}

func createColorRole(s *discordgo.Session, guildID, userID string, color string) error {
	role, err := s.GuildRoleCreate(guildID)
	if err != nil {
		return fmt.Errorf("creating role: %w", err)
	}

	role.Name = color
	role.Color = hexColorToInt(color[1:])
	_, err = s.GuildRoleEdit(guildID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, false)
	if err != nil {
		return fmt.Errorf("editing role: %w", err)
	}

	err = s.GuildMemberRoleAdd(guildID, userID, role.ID)
	if err != nil {
		return fmt.Errorf("adding guild member role: %w", err)
	}

	return nil
}

func getRoleUserIDs(s *discordgo.Session, guildID, roleID string) ([]string, error) {
	users, err := s.GuildMembers(guildID, "", 250)
	if err != nil {
		return nil, fmt.Errorf("getting role user ids: %w", err)
	}

	ids := make([]string, 0, len(users))
	for _, u := range users {
		for _, rid := range u.Roles {
			if rid == roleID {
				ids = append(ids, u.User.ID)
				break
			}
		}
	}

	return ids, nil
}

func getGuildRolesByID(s *discordgo.Session, guildID string) (map[string]*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("getting guild roles: %w", err)
	}

	rm := make(map[string]*discordgo.Role, len(roles))
	for _, r := range roles {
		rm[r.ID] = r
	}

	return rm, nil
}

func getAuthorColorRole(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	roles map[string]*discordgo.Role,
) (*discordgo.Role, error) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("getting guildmember: %w", err)
	}

	for _, r := range member.Roles {
		role := roles[r]

		if colorRegex.MatchString(role.Name) {
			return role, nil
		}
	}

	return nil, nil
}

func hexColorToInt(hex string) int {
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return int((r << 16) | (g << 8) | (b))
}
