package config

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var Tokens []string
var Sessions []*discordgo.Session
var MessageCache = make(map[string]map[string]string)
var StartTime time.Time
var CategoryName = "intercité"
var AdminGuild = "486564991915393045"
var BansS Bans
var Softbans SoftBans
var Userbans UserBans
