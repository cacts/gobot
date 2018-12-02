package handlers

import (
	"encoding/base64"
	"io"
	"regexp"
	"fmt"
	"os/exec"
	"github.com/bwmarrin/discordgo"
	"github.com/cactauz/gobot"
)

func init() {
	gobot.Global.AddMessageHandler(NewPoeWikiHandler())
}

var (
	regex, _ = regexp.Compile("\\[\\[([^\\[\\]]*)\\]\\]")
	buf = make([]byte, 1024, 1024)
)

type PoeWikiHandler struct {
	nodeClient *nodeClient
}

func NewPoeWikiHandler() *PoeWikiHandler {
	cmd := exec.Command("node", "../poe_wiki_v2.js")
	stdin, err := cmd.StdinPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error:", err)
	}

	nodeClient := &nodeClient{
		stdIn:  stdin, 
		stdOut: stdout,
	}

	cmd.Start()

	go func() {
		err := cmd.Wait()
		fmt.Println("node process exited:", err)
	}()

	return &PoeWikiHandler{
		nodeClient: nodeClient,
	}
}

type nodeClient struct {
	stdIn  io.Writer
	stdOut io.Reader
}

func (wh *PoeWikiHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !regex.MatchString(m.Content) {
		return
	}

	item := regex.FindStringSubmatch(m.Content)[1]

	_, err := wh.nodeClient.stdIn.Write([]byte(item+"\n"))
	if err != nil {
		fmt.Println("err writing:", err)
		return
	} 

	out := decodeResponse(wh.nodeClient.stdOut)
	dec := base64.NewDecoder(base64.StdEncoding, out)

	file := &discordgo.File{
		Name: "item.png",
		ContentType: "image/png",
		Reader: dec,
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{file},
	})
	if err != nil {
		fmt.Println("err sending message:", err)
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		fmt.Println("err deleting message:", err)
	}
}

func decodeResponse(itemReader io.Reader) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		for {
			n, err := itemReader.Read(buf)

			if n > 0 {
				_, err := pw.Write(buf[:n])
				if err != nil {
					fmt.Println("err writing:", err) 
					return
				}

				if string(buf[:3]) == "failed" {
					fmt.Println("item lookup failed")
					return
				}

				if buf[n-1] == byte('\n') {
					pw.Close()
					break
				}
			}

			if err != nil {
				if err == io.EOF {
					fmt.Println("got EOF, this shouldn't happen")
					return
				}

				fmt.Println("error:", err)
				return
			}
		}
	}()

	return pr
}

func (wh *PoeWikiHandler) Name() string {
	return "poe_wiki_handler"
}
