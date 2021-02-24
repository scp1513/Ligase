package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/finogeeks/ligase/sdk/go-ligase"
	"github.com/finogeeks/ligase/sdk/go-ligase/logger"
)

var (
	botNum        = flag.Int("botNum", 1, "bot number")
	createRoomNum = flag.Int("createRoomNum", 1, "create room number")
	sendInterval  = flag.Int("sendInterval", 100, "send interval")
	botIdx        = flag.Int("botIdx", 0, "bot userID index suffix")
	botPrefix     = flag.String("botPrefix", "bot0_", "bot userID prefix")
	botDevPrefix  = flag.String("botDevPrefix", "botdev0_", "bot device prefix")
	domain        = flag.String("domain", "", "ligase domain")
	url           = flag.String("url", "", "ligase url")
	concurrent    = flag.Int("concurrent", 100, "concurrent limit")

	help = flag.Bool("h", false, "show this message")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	if domain == nil || *domain == "" {
		flag.PrintDefaults()
		return
	}

	if url == nil || *url == "" {
		flag.PrintDefaults()
		return
	}

	ligase.Init(ligase.ClientOpts{
		URL:                 *url,
		MaxConnsPerHost:     *concurrent,
		MaxIdleConnsPerHost: *concurrent,

		LogTime:   true,
		LogFile:   true,
		LogColors: true,
		LogDebug:  false,
		LogTrace:  false,
	})

	bots := make(map[string]*Bot, *botNum)
	for i := 0; i < *botNum; i++ {
		userID := "@" + *botPrefix + strconv.Itoa(*botIdx+i) + ":" + *domain
		bot := NewBot(*botDevPrefix, userID)
		bots[userID] = bot
	}

	b := bots["@"+*botPrefix+strconv.Itoa(*botIdx)+":"+*domain]
	err := b.Login()
	if err != nil {
		logger.GetLogger().Errorf("%s login failed %s", b.GetUserID(), err.Error())
		panic(err)
	}

	ts0 := time.Now()
	var wgLogin sync.WaitGroup
	wgLogin.Add(*botNum - 1)
	for _, bot := range bots {
		if b.GetUserID() == bot.GetUserID() {
			continue
		}
		go func(bot *Bot) {
			defer wgLogin.Done()
			err := bot.Login()
			if err != nil {
				logger.GetLogger().Errorf("%s login failed %s", bot.GetUserID(), err.Error())
			}
		}(bot)
	}

	var wgCreateRoom sync.WaitGroup
	wgCreateRoom.Add(*createRoomNum)
	roomIDs := make(map[string]struct{}, *createRoomNum)
	for i := 0; i < *createRoomNum; i++ {
		go func() {
			defer wgCreateRoom.Done()
			roomID, err := b.CreateRoom()
			if err != nil {
				logger.GetLogger().Errorf("%s createRoom failed %s", b.GetUserID(), err.Error())
				return
			}
			roomIDs[roomID] = struct{}{}
		}()
	}

	wgLogin.Wait()
	wgCreateRoom.Wait()
	ts1 := time.Now()
	logger.GetLogger().Noticef("login & createRoom finished use time %s", ts1.Sub(ts0))

	ts0 = time.Now()
	var wgJoin sync.WaitGroup
	wgJoin.Add(*botNum - 1)
	for _, bot := range bots {
		if b.GetUserID() == bot.GetUserID() {
			continue
		}
		go func(bot *Bot) {
			defer wgJoin.Done()
			for k := range roomIDs {
				err := bot.JoinRoom(k)
				if err != nil {
					logger.GetLogger().Errorf("%s login failed %s", bot.GetUserID(), err.Error())
				}
			}
		}(bot)
	}

	for _, bot := range bots {
		go func(bot *Bot) {
			for {
				_, err := bot.Sync()
				if err != nil {
					logger.GetLogger().Errorf("%s login failed %s", bot.GetUserID(), err.Error())
				}
			}
		}(bot)
	}

	wgJoin.Wait()
	ts1 = time.Now()
	logger.GetLogger().Noticef("joinRoom & startSync finished use time %s", ts1.Sub(ts0))

	for _, bot := range bots {
		go func(bot *Bot) {
			sendCount := 0
			ts0 := time.Now()
			for {
				for roomID := range roomIDs {
					_, err := bot.SendMessage(roomID, fmt.Sprintf("%s send to %s at %s", bot.GetUserID(), roomID, time.Now()))
					if err != nil {
						logger.GetLogger().Errorf("%s sendMessage failed %s", bot.GetUserID(), err.Error())
					}
					sendCount++
					if sendCount%100 == 0 {
						ts1 := time.Now()
						logger.GetLogger().Noticef("%s sendMessage count %d use time %s", bot.GetUserID(), sendCount, ts1.Sub(ts0))
					}
					time.Sleep(time.Millisecond * time.Duration(*sendInterval))
				}
			}
		}(bot)
	}

	select {}
}
