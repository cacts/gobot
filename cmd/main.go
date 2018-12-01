package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"../../gobot"
	// load handler
	_ "../handlers"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	gobot.Global.Open(Token)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	gobot.Global.Close()
}
