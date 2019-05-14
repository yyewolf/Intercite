package main

import (
	"Intercite/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AdminCommands(session *discordgo.Session, m *discordgo.MessageCreate, Prefix string) {
	RealCmd := m.ContentWithMentionsReplaced()
	if strings.HasPrefix(RealCmd, Prefix+" userban ") {
		if m.Author.ID == "144472011924570113" {
			msgID := strings.Replace(RealCmd, Prefix+" userban ", "", 1)
			toBanID := config.MessageCache[msgID]["authorid"]
			if toBanID != "" {
				config.Userbans.Banned = append(config.Userbans.Banned, toBanID)
				config.SaveBans()
				session.ChannelMessageSend(m.ChannelID, "Banned user : "+toBanID)
			}
		}
	}

	if strings.HasPrefix(RealCmd, Prefix+" ban ") {
		if m.Author.ID == "144472011924570113" {
			msgID := strings.Replace(RealCmd, Prefix+" ban ", "", 1)
			toBanID := config.MessageCache[msgID]["guildOrigin"]
			if toBanID != "" {
				config.BansS.Banned = append(config.BansS.Banned, toBanID)
				for i := range config.Sessions {
					_ = config.Sessions[i].GuildLeave(toBanID)
				}
				config.SaveBans()
				session.ChannelMessageSend(m.ChannelID, "Banned guild : "+toBanID)
			}
		}
	}

	if strings.HasPrefix(RealCmd, Prefix+" softban ") {
		if m.Author.ID == "144472011924570113" {
			msgID := strings.Replace(RealCmd, Prefix+" softban ", "", 1)
			toBanID := config.MessageCache[msgID]["guildOrigin"]
			if toBanID != "" {
				config.Softbans.Banned = append(config.Softbans.Banned, toBanID)
				config.SaveBans()
				session.ChannelMessageSend(m.ChannelID, "Banned guild : "+toBanID)
			}
		}
	}

	if strings.HasPrefix(RealCmd, Prefix+" pardon ") {
		if m.Author.ID == "144472011924570113" {
			msgID := strings.Replace(RealCmd, Prefix+" softban ", "", 1)
			config.RemoveBan(config.MessageCache[msgID]["guildOrigin"])
			config.RemoveBan(config.MessageCache[msgID]["authorid"])
			config.RemoveBan(msgID)
		}
	}
}
