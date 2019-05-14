package main

import (
	"Intercite/config"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func BroadCastMessage(msg string, msgID string, channelName string, originalServer string) {
	var sentChannels []string
	for i := range config.Sessions {
		sess := config.Sessions[i]
		for j := range sess.State.Guilds {
			if sess.State.Guilds[j].ID != originalServer {
				for k := range sess.State.Guilds[j].Channels {
					if sess.State.Guilds[j].Channels[k].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(sess.State.Guilds[j].Channels[k].Name) == config.CategoryName {
						InterID := sess.State.Guilds[j].Channels[k].ID
						for l := range sess.State.Guilds[j].Channels {
							if sess.State.Guilds[j].Channels[l].Type == discordgo.ChannelTypeGuildText && sess.State.Guilds[j].Channels[l].Name == channelName && sess.State.Guilds[j].Channels[l].ParentID == InterID {
								exist := false
								for alreadysent := range sentChannels {
									if sentChannels[alreadysent] == sess.State.Guilds[j].Channels[l].ID {
										exist = true
										break
									}
								}
								if !exist {
									Sent, err := sess.ChannelMessageSend(sess.State.Guilds[j].Channels[l].ID, msg)
									if err != nil {
										//
									} else {
										sentChannels = append(sentChannels, Sent.ChannelID)
										config.MessageCache[Sent.ID] = make(map[string]string)
										config.MessageCache[Sent.ID]["content"] = Sent.Content
										config.MessageCache[Sent.ID]["authorname"] = Sent.Author.Username
										config.MessageCache[Sent.ID]["authorid"] = config.MessageCache[msgID]["authorid"]
										config.MessageCache[Sent.ID]["guildOrigin"] = originalServer
									}
									break
								}
							}
						}
					}
				}
			}
		}
	}
}

func BroadCastDelete(msg string, msgID string, channelName string, originalServer string) {
	for i := range config.Sessions {
		sess := config.Sessions[i]
		for j := range sess.State.Guilds {
			if sess.State.Guilds[j].ID != originalServer {
				for k := range sess.State.Guilds[j].Channels {
					if sess.State.Guilds[j].Channels[k].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(sess.State.Guilds[j].Channels[k].Name) == config.CategoryName {
						InterID := sess.State.Guilds[j].Channels[k].ID
						for l := range sess.State.Guilds[j].Channels {
							if sess.State.Guilds[j].Channels[l].Type == discordgo.ChannelTypeGuildText && sess.State.Guilds[j].Channels[l].Name == channelName && sess.State.Guilds[j].Channels[l].ParentID == InterID {
								Messages, err := sess.ChannelMessages(sess.State.Guilds[j].Channels[l].ID, 10, "", "", msgID)
								if err != nil {
									//
								}
								for m := range Messages {
									if Messages[m].Content == msg {
										err = sess.ChannelMessageDelete(sess.State.Guilds[j].Channels[l].ID, Messages[m].ID)
										if err != nil {
											//
										}
										delete(config.MessageCache, Messages[m].ID)
										break
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func BroadCastEdit(msg string, msgID string, newMsg string, channelName string, originalServer string) {
	for i := range config.Sessions {
		sess := config.Sessions[i]
		for j := range sess.State.Guilds {
			if sess.State.Guilds[j].ID != originalServer {
				for k := range sess.State.Guilds[j].Channels {
					if sess.State.Guilds[j].Channels[k].Type == discordgo.ChannelTypeGuildCategory && strings.ToLower(sess.State.Guilds[j].Channels[k].Name) == config.CategoryName {
						InterID := sess.State.Guilds[j].Channels[k].ID
						for l := range sess.State.Guilds[j].Channels {
							if sess.State.Guilds[j].Channels[l].Type == discordgo.ChannelTypeGuildText && sess.State.Guilds[j].Channels[l].Name == channelName && sess.State.Guilds[j].Channels[l].ParentID == InterID {
								Messages, err := sess.ChannelMessages(sess.State.Guilds[j].Channels[l].ID, 10, "", "", msgID)
								if err != nil {
									//
								}
								for m := range Messages {
									if Messages[m].Content == msg {
										_, err = sess.ChannelMessageEdit(sess.State.Guilds[j].Channels[l].ID, Messages[m].ID, newMsg)
										if err != nil {
											//
										}
										break
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func RawToFormatted(RawMessage string, UserName string, ServerName string, Attachments []*discordgo.MessageAttachment) string {
	var re = regexp.MustCompile(`:`)
	ServerName = re.ReplaceAllString(ServerName, "")
	UserName = re.ReplaceAllString(UserName, "")
	for i := range Attachments {
		RawMessage += "￼" + "\n" + Attachments[i].URL
	}
	Msg := "[_" + ServerName + "_]\n"
	Msg += "**" + UserName + "** : " + RawMessage
	return Msg
}

func FormattedToRaw(FormattedMessage string) string {
	SplitMsg := strings.Split(FormattedMessage, ": ")
	if len(SplitMsg) > 1 {
		Msg := strings.Split(SplitMsg[1], "￼")[0]
		return Msg
	}
	return ""
}
