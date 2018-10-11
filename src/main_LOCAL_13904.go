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
var Commands = make(map[string]*f.Command)

func init() {
	err := dat.Load("oofcounter/config.json", &myConfig)
	if err != nil {
		dat.Log.Println(err)
	}
	defer f.Session.AddHandler(OofCount)
	Commands["oof"] = &f.Command{
		Name:    "oof",
		Help:    "oof",
		Perms:   -1,
		Version: "oof",
		Action:  oof,
	}
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
	server, err := f.GetGuild(s, m.Message)
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
	guildOofs := myConfig.Guild[server.ID]
	guildOofs.Epoch = time.Now()
	if guildOofs.OofCount == nil {
		guildOofs.OofCount = make(map[string]int)
	}
	guildOofs.OofCount[m.Message.Author.ID]++
	guildOofs.TotalOofs++
	myConfig.Guild[server.ID] = guildOofs
	dat.Save("oofcounter/config.json", &myConfig)
	for _, channel := range guildOofs.BlChans {
		if channel == m.Message.ChannelID {
			return
		}
	}
	check, err = regexp.MatchString("[A-z]+oof|oof[A-z]+|"+f.Config.Prefix+"oof", m.Message.Content)
	if err != nil {
		dat.Log.Println(err)
		return
	}
	if !check {
		if guildOofs.ReplyFrequency == 0 {
			guildOofs.ReplyFrequency = 100
		}
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(1000)
		if number <= guildOofs.ReplyFrequency {
			oofPercent := float64(100) * (float64(guildOofs.OofCount[m.Message.Author.ID]) / float64(guildOofs.TotalOofs))
			s.ChannelMessageSend(m.Message.ChannelID, fmt.Sprintf("**oof** indeed! You've oof'd %d times! Thats %f%% of all oofs in the server (%d) since I started counting at Epoch %v",
				guildOofs.OofCount[m.Message.Author.ID],
				oofPercent,
				guildOofs.TotalOofs,
				guildOofs.Epoch.Format("Mon, 2 Jan 2006 at 15:04.")))
		}
	}
}

func oof(session *discordgo.Session, message *discordgo.Message) {
	server, err := f.GetGuild(session, message)
	if err != nil {
		dat.AlertDiscord(session, message, err)
		dat.Log.Println(err)
		return
	}
	if myConfig.Guild == nil {
		session.ChannelMessageSend(message.ChannelID, "Nobody has oof'd anywhere! Try oofing?")
		return
	}
	guildOofs := myConfig.Guild[server.ID]
	oofPercent := float64(100) * (float64(guildOofs.OofCount[message.Author.ID]) / float64(guildOofs.TotalOofs))
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("**oof** indeed! You've oof'd %d times! Thats %f%% of all oofs in the server (%d) since I started counting at Epoch %v",
		guildOofs.OofCount[message.Author.ID],
		oofPercent,
		guildOofs.TotalOofs,
		guildOofs.Epoch.Format("Mon, 2 Jan 2006 at 15:04.")))
}

//func blacklistOofs(session *discordgo.Session, message *discordgo.Message) {

//}
