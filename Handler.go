package main

import (
	"Intercite/config"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CommandHandle(session *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	BotUsr, err := session.User("@me")

	if user.ID == BotUsr.ID || user.Bot || err != nil {
		return
	}

	RealCmd := m.ContentWithMentionsReplaced()
	Prefix := "@" + BotUsr.Username

	if strings.HasPrefix(RealCmd, Prefix+" restart") {
		if m.Author.ID == "144472011924570113" {
			os.Exit(0)
		}
	}
	if strings.HasPrefix(RealCmd, Prefix+" help") {
		embed := &discordgo.MessageEmbed{
			Title:       "Help menu",
			Description: "These are the principals commands of this bot. (@CCTS command)",
			Color:       0xFFDD00,
			Fields: []*discordgo.MessageEmbedField{{
				Name: "__Commands :__",
				Value: "**info** : Gets statistics of the bot.\n" +
					"**officials** : Lists all the officials channels of the bot.\n" +
					"**help** : Summons this menu.",
			}, {
				Name:  "__Support :__",
				Value: "The bot is currently being updated to Golang, to follow the development, you can join this server : https://discord.gg/Q8NbBeQ",
			}},
		}
		session.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.HasPrefix(RealCmd, Prefix+" info") {
		users := 0
		for _, guild := range session.State.Ready.Guilds {
			users += len(guild.Members)
		}
		ServerAmount := -9
		for i := range config.Sessions {
			ServerAmount += len(config.Sessions[i].State.Guilds)
		}
		Uptime := time.Since(config.StartTime)

		embed := &discordgo.MessageEmbed{
			Title: "Bot Statistics :",
			Description: "**Servers** : " + strconv.Itoa(ServerAmount) + "\n" +
				"**Users** : " + strconv.Itoa(users) + "\n" +
				"**Tasks** : " + strconv.Itoa(runtime.NumGoroutine()) + "\n" +
				"**Uptime** : " + strconv.Itoa(int(Uptime.Hours())) + ":" + strconv.Itoa(int(Uptime.Minutes())%60) + ":" + strconv.Itoa(int(Uptime.Seconds())%60),
			Color: 0xFFDD00,
		}
		session.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.HasPrefix(RealCmd, Prefix+" officials") {
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
		Desc := ""
		for i := range MainGuild.Channels {
			if MainGuild.Channels[i].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(MainGuild.Channels[i].Name) == config.CategoryName {
				for k := range MainGuild.Channels {
					if MainGuild.Channels[k].ParentID == MainGuild.Channels[i].ID {
						Desc += MainGuild.Channels[k].Name + "\n"
					}
				}
			}
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Official channels :",
			Description: Desc,
			Color:       0xFFDD00,
		}
		session.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.HasPrefix(RealCmd, Prefix+" all") {
		if m.Author.ID == "144472011924570113" {
			StringChannels := ""
			for n := range config.Sessions {
				for i := range config.Sessions[n].State.Guilds {
					for k := range config.Sessions[n].State.Guilds[i].Channels {
						if config.Sessions[n].State.Guilds[i].Channels[k].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(config.Sessions[n].State.Guilds[i].Channels[k].Name) == config.CategoryName {
							for j := range config.Sessions[n].State.Guilds[i].Channels {
								if config.Sessions[n].State.Guilds[i].Channels[j].ParentID == config.Sessions[n].State.Guilds[i].Channels[k].ID && !strings.Contains(StringChannels, config.Sessions[n].State.Guilds[i].Channels[j].Name) {
									StringChannels += config.Sessions[n].State.Guilds[i].Channels[j].Name + "\n"
								}
							}
						}
					}
				}
			}
			r := bytes.NewReader([]byte(StringChannels))
			session.ChannelFileSendWithMessage(m.ChannelID, "Here is the file :", "Channels.txt", r)
		}
	}
	AdminCommands(session, m, Prefix)
}

func Edit(s *discordgo.Session, evt *discordgo.MessageUpdate) {
	if evt.Author != nil {
		if !evt.Author.Bot {
			RawMsg, err := evt.ContentWithMoreMentionsReplaced(s)
			if err != nil {
				//
			}
			g, err := s.State.Guild(evt.GuildID)
			if err != nil {
				//
			}
			c, err := s.State.Channel(evt.ChannelID)
			if err != nil {
				//
			}
			p, err := s.State.Channel(c.ParentID)
			if err != nil {
				//
			}
			if p.Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(p.Name) == config.CategoryName {
				Formatted := RawToFormatted(RawMsg, config.MessageCache[evt.ID]["authorname"], g.Name, evt.Attachments)
				go BroadCastEdit(Formatted, evt.ID, Formatted, c.Name, evt.GuildID)
				go BroadCastEdit(config.MessageCache[evt.ID]["formatted"], evt.ID, Formatted, c.Name, evt.GuildID)
				config.MessageCache[evt.ID]["content"] = RawMsg
				config.MessageCache[evt.ID]["formatted"] = Formatted
			}
		}
	}
}

func Send(s *discordgo.Session, evt *discordgo.MessageCreate) {
	if !evt.Author.Bot && !config.IsBanned(evt.GuildID, evt.Author.ID) {
		RawMsg, err := evt.ContentWithMoreMentionsReplaced(s)
		if err != nil {
			//
		}
		g, err := s.State.Guild(evt.GuildID)
		if err != nil {
			//
		}
		c, err := s.State.Channel(evt.ChannelID)
		if err != nil {
			//
		}
		p, err := s.State.Channel(c.ParentID)
		if err != nil {
			//
		}
		if p.Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(p.Name) == config.CategoryName && !evt.MentionEveryone {
			Formatted := RawToFormatted(RawMsg, evt.Author.Username, g.Name, evt.Attachments)
			config.MessageCache[evt.ID] = make(map[string]string)
			config.MessageCache[evt.ID]["content"] = RawMsg
			config.MessageCache[evt.ID]["formatted"] = Formatted
			config.MessageCache[evt.ID]["authorname"] = evt.Author.Username
			config.MessageCache[evt.ID]["authorid"] = evt.Author.ID
			BroadCastMessage(Formatted, evt.ID, c.Name, evt.GuildID)
		} else if p.Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(p.Name) == config.CategoryName && evt.MentionEveryone {
			s.ChannelMessageSend(evt.ChannelID, "You cannot tag `@everyone` or `@here` in this channel !")
		}
	}
}

func Delete(s *discordgo.Session, evt *discordgo.MessageDelete) {
	RawMsg := config.MessageCache[evt.ID]["content"]
	g, err := s.State.Guild(evt.GuildID)
	if err != nil {
		//
	}
	c, err := s.State.Channel(evt.ChannelID)
	if err != nil {
		//
	}
	p, err := s.State.Channel(c.ParentID)
	if err != nil {
		//
	}
	if p.Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(p.Name) == config.CategoryName {
		Formatted := RawToFormatted(RawMsg, config.MessageCache[evt.ID]["authorname"], g.Name, evt.Attachments)
		go BroadCastDelete(Formatted, evt.ID, c.Name, evt.GuildID)
		go BroadCastDelete(config.MessageCache[evt.ID]["formatted"], evt.ID, c.Name, evt.GuildID)
	}
	if g.ID == config.AdminGuild && config.MessageCache[evt.ID]["content"] != "" {
		Raw := FormattedToRaw(config.MessageCache[evt.ID]["content"])
		go BroadCastDelete(Raw, evt.ID, c.Name, evt.GuildID)
	}
	delete(config.MessageCache, evt.ID)
}
