package main

import (
	"fmt"
	"github.com/aurumbot/lib/dat"
	f "github.com/aurumbot/lib/foundation"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"regexp"
	"time"
)

type Config struct {
	Guild map[string]struct {
		Epoch          time.Time      `json:"epoch"`
		OofCount       map[string]int `json:"oofcount"`
		BlChans        []string       `json:"blacklist"`
		TotalOofs      int            `json:"totaloofs"`
		ReplyFrequency int            `json:"replyfrequency"`
	} `json:"guild"`
}

var myConfig Config

func init() {
	dat.Load("oof/config.json", &myConfig)
	f.Session.AddHandler(OofCount)
}

func OofCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message.Author.Bot {
		return
	}
	check, err := regexp.MatchString(".*oof.*", m.Message.Content)
	if err != nil {
		dat.Log.Println(err)
		return
	}
	if !check {
		return
	}
	guild, err := f.GetGuild(s, m.Message)
	if err != nil {
		dat.Log.Println(err)
		return
	}
	if myConfig.Guild == nil {
		myConfig.Guild = make(map[string]struct {
			Epoch          time.Time      `json:"epoch"`
			OofCount       map[string]int `json:"oofcount"`
			BlChans        []string       `json:"blacklist"`
			TotalOofs      int            `json:"totaloofs"`
			ReplyFrequency int            `json:"replyfrequency"`
		})
	}
	/*if myConfig.Guild[guild.ID].OofCount == nil {
		myConfig.Guild[guild.ID].OofCount = make(map[string]int)
	}*/
	myConfig.Guild[guild.ID].OofCount[m.Message.Author.ID]++
	myConfig.Guild[guild.ID].TotalOofs++
	dat.Save("oof/myConfig.json", &myConfig)
	for _, channel := range myConfig.Guild[guild.ID].BlChans {
		if channel == m.Message.ChannelID {
			return
		}
	}
	check, err = regexp.MatchString("[A-z]+oof|oof[A-z]+", m.Message.Content)
	if err != nil {
		dat.Log.Println(err)
		return
	}
	if !check {
		if myConfig.Guild[guild.ID].ReplyFrequency == 0 {
			myConfig.Guild[guild.ID].ReplyFrequency = 100
		}
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(1000)
		if number <= myConfig.Guild[guild.ID].ReplyFrequency {
			s.ChannelMessageSend(m.Message.ChannelID, fmt.Sprintf("**oof** indeed! You've oof'd %d times! Thats %f% of all oofs in the server (%d) since I started counting at Epoch %v",
				myConfig.Guild[guild.ID].OofCount[m.Message.Author.ID],
				(100*(myConfig.Guild[guild.ID].OofCount[m.Message.Author.ID]/myConfig.Guild[guild.ID].TotalOofs)),
				myConfig.Guild[guild.ID].Epoch.Format("Mon, 2 Jan 2006 at 15:04.")))
		}
	}
}
