package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"Intercite/config"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}

func GuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	do := true
	Guild, err := s.Guild(event.Guild.ID)
	if err != nil {
		//
	}
	ids := []string{}
	for i := range config.Sessions {
		if config.Sessions[i].State.User.ID != s.State.User.ID {
			ids = append(ids, config.Sessions[i].State.User.ID)
		}
	}
	for i := range Guild.Members {
		for k := range ids {
			if Guild.Members[i].User.ID == ids[k] || config.IsBanned(event.Guild.ID, "") {
				do = false
				s.GuildLeave(event.Guild.ID)
				break
			}
		}
	}
	if do {
		rightSession := &discordgo.Session{}
		for i := range config.Sessions {
			if config.Sessions[i].State.User.ID == "529644145107533824" {
				rightSession = config.Sessions[i]
				break
			}
		}
		MainGuild, err := rightSession.Guild(config.AdminGuild)
		if err != nil {
			fmt.Println(err)
		}
		for i := range MainGuild.Channels {
			if MainGuild.Channels[i].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(MainGuild.Channels[i].Name) == config.CategoryName {
				for k := range MainGuild.Channels {
					if MainGuild.Channels[k].ParentID == MainGuild.Channels[i].ID {
						data := discordgo.GuildChannelCreateData{
							ParentID: MainGuild.Channels[i].ID,
							Name:     MainGuild.Channels[k].Name,
							Type:     discordgo.ChannelTypeGuildText,
						}
						s.GuildChannelCreateComplex(event.Guild.ID, data)
					}
				}
			}
		}
	}
}

func botReady(s *discordgo.Session, event *discordgo.Ready) {
	color.Green("Bot is ready ! (" + s.State.User.Username + ")")
	exist := false
	for i := range config.Sessions {
		if config.Sessions[i].State.User.ID == s.State.User.ID {
			exist = true
			break
		}
	}
	if !exist {
		config.Sessions = append(config.Sessions, s)
	}
}

func connect() {
	for i := range config.Tokens {
		dg, err := discordgo.New("Bot " + config.Tokens[i])
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}
		dg.AddHandler(botReady)
		dg.AddHandler(CommandHandle)
		dg.AddHandler(Send)
		dg.AddHandler(Delete)
		dg.AddHandler(Edit)
		err = dg.Open()
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func close() {
	for i := range config.Sessions {
		config.Sessions[i].Close()
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	config.LoadBans()
	config.DefineTokens()
	// Create a new Discord session using the provided bot token.
	go connect()
	config.StartTime = time.Now()
	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	close()
}
