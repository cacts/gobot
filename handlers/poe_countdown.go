package handlers

import (
	"fmt"
	"time"

	"github.com/cactauz/gobot" 
	"github.com/bwmarrin/discordgo"
)

var (
	target = time.Date(2018, 12, 7, 19, 0, 0, 0, time.UTC)
)

func init() {
	if time.Now().Sub(target) < 0 {
	 gobot.Global.AddTickHandler(NewSimpleTickHandler("poe_countdown", 5, poeTicker))
	}
}

func poeTicker(sess *discordgo.Session) {
	timeleft := target.Sub(time.Now()).Round(time.Second)

	d := timeleft / (24 * time.Hour)
	timeleft -= 24 * time.Hour * d
	h := timeleft / time.Hour
	timeleft -= h * time.Hour
	m := timeleft / time.Minute
	timeleft -= m * time.Minute
	s := timeleft / time.Second

	var str string
	if timeleft < 0 {
		sess.UpdateStatus(0, "GL WITH UR MEPS!!")
		gobot.Global.RemoveTickHandler("poe_countdown")
	} else {
		if d > 0 {
			str += fmt.Sprintf("%dd ", d)
		}
		if d > 0 || h > 0 { 
			str += fmt.Sprintf("%dh ", h)
		}
		if d > 0 || h > 0 || m > 0 {
			str += fmt.Sprintf("%dm ", m)
		}
		str += fmt.Sprintf("%ds", s)
	}

	sess.UpdateStatus(0, fmt.Sprintf("%s TIL POE", str))
}
