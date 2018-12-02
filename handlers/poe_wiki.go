package handlers

import (
	"encoding/base64"
	"regexp"
	"fmt"
	"bytes"
	"os/exec"
	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

func init() {
	gobot.Global.AddMessageHandler(&PoeWikiHandler{})
}

var (
	regex, _ = regexp.Compile("\\[\\[([^\\[\\]]*)\\]\\]")
	command = "./node poe_wiki.js"
)

type PoeWikiHandler struct {
}

func (wh *PoeWikiHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !regex.MatchString(m.Content) {
		return
	}

	item := regex.FindStringSubmatch(m.Content)[1]

	cmd := exec.Command("node", "../poe_wiki.js", item)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	dec := base64.NewDecoder(base64.StdEncoding, &out)

	if err != nil {
		fmt.Println("err calling node:", err)
		return
	}

	file := &discordgo.File{
		Name: "item.png",
		ContentType: "image/png",
		Reader: dec,
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{file},
	})
}

func (wh *PoeWikiHandler) Name() string {
	return "poe_wiki_handler"
}
