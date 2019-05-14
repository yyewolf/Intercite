package Permissions

import "github.com/bwmarrin/discordgo"

var CREATE_INSTANT_INVITE = 0x00000001
var KICK_MEMBERS = 0x00000002
var BAN_MEMBERS = 0x00000004
var ADMINISTRATOR = 0x00000008
var MANAGE_CHANNELS = 0x00000010
var MANAGE_GUILD = 0x00000020
var ADD_REACTIONS = 0x00000040
var VIEW_AUDIT_LOG = 0x00000080
var VIEW_CHANNEL = 0x00000400
var SEND_MESSAGES = 0x00000800
var SEND_TTS_MESSAGES = 0x00001000
var MANAGE_MESSAGES = 0x00002000
var EMBED_LINKS = 0x00004000
var ATTACH_FILES = 0x00008000
var READ_MESSAGE_HISTORY = 0x00010000
var MENTION_EVERYONE = 0x00020000
var USE_EXTERNAL_EMOJIS = 0x00040000
var CONNECT = 0x00100000
var SPEAK = 0x00200000
var MUTE_MEMBERS = 0x00400000
var DEAFEN_MEMBERS = 0x00800000
var MOVE_MEMBERS = 0x01000000
var USE_VAD = 0x02000000
var PRIORITY_SPEAKER = 0x00000100
var CHANGE_NICKNAME = 0x04000000
var MANAGE_NICKNAMES = 0x08000000
var MANAGE_ROLES = 0x10000000
var MANAGE_WEBHOOKS = 0x20000000
var MANAGE_EMOJIS = 0x40000000

func HasPermission(User *discordgo.Member, s *discordgo.Session, GuildID string, Perm int) (isHe bool) {
	g, err := s.Guild(GuildID)
	if err != nil {
		return false
	}
	for _, roleID := range User.Roles {
		for i := 0; i < len(g.Roles); i++ {
			if g.Roles[i].ID == roleID {
				if g.Roles[i].Permissions&Perm == Perm {
					return true
				}
			}
		}
	}
	return false
}
