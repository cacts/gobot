package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cactauz/gobot"
	// load handlers
	_ "github.com/cactauz/gobot/handlers"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	if Token == "" {
		Token = os.Getenv("GOBOT_TOKEN")
	}

	gobot.Global.Open(Token)
	defer gobot.Global.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Type exit to quit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	inCh := make(chan string)
	go func() {
		defer close(inCh)

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			text, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("error:", err)
			}

			if text == "exit\n" || text == "exit\r\n" { // windows :)
				sc <- os.Interrupt
				return
			}

			inCh <- text
		}
	}()

	for {
		select {
		case <-sc:
			return
		case text := <-inCh:
			if text != "" {
				gobot.Global.SendMessage(text)
			}
		}
	}
}
