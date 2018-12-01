package gobot

import (
	"sync"
	"time"
)

const (
	GLOBAL_TIMEOUT = 2 * time.Minute
)

var mutex sync.Mutex
var timeoutChans = map[string]chan bool{}

var GlobalTimeout = func(channelID string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	return getTimeout(channelID)
}

func getTimeout(channelID string) bool {
	tc, ok := timeoutChans[channelID]
	if !ok {
		tc = make(chan bool, 1)
		timeoutChans[channelID] = tc
		timeoutChans[channelID] <- true
	}

	select {
	case t := <-tc:
		go func() {
			time.Sleep(GLOBAL_TIMEOUT)
			timeoutChans[channelID] <- true
		}()
		return t
	default:
		return false
	}
}
