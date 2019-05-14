package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type UserBans struct {
	Banned []string
}

type SoftBans struct {
	Banned []string
}

type Bans struct {
	Banned []string
}

func RemoveBan(ID string) {
	found := false
	for i := range BansS.Banned {
		if BansS.Banned[i] == ID {
			found = true
			BansS.Banned[i] = BansS.Banned[len(BansS.Banned)-1]
			BansS.Banned[len(BansS.Banned)-1] = ""
			BansS.Banned = BansS.Banned[:len(BansS.Banned)-1]
			break
		}
	}
	if !found {
		for i := range Softbans.Banned {
			if Softbans.Banned[i] == ID {
				found = true
				Softbans.Banned[i] = Softbans.Banned[len(Softbans.Banned)-1]
				Softbans.Banned[len(Softbans.Banned)-1] = ""
				Softbans.Banned = Softbans.Banned[:len(Softbans.Banned)-1]
				break
			}
		}
	}
	if !found {
		for i := range Userbans.Banned {
			if Userbans.Banned[i] == ID {
				Userbans.Banned[i] = Userbans.Banned[len(Userbans.Banned)-1]
				Userbans.Banned[len(Userbans.Banned)-1] = ""
				Userbans.Banned = Userbans.Banned[:len(Userbans.Banned)-1]
				break
			}
		}
	}
	SaveBans()
}

func IsBanned(GuildID string, UserID string) bool {
	banned := false
	for i := range BansS.Banned {
		if BansS.Banned[i] == GuildID {
			banned = true
			break
		}
	}
	if !banned {
		for i := range Softbans.Banned {
			if Softbans.Banned[i] == GuildID {
				banned = true
				break
			}
		}
	}
	if !banned {
		for i := range Userbans.Banned {
			if Userbans.Banned[i] == UserID {
				banned = true
				break
			}
		}
	}
	return banned
}

func LoadBans() {
	path, _ := filepath.Abs("./bans/bans.json")
	file, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &BansS)
	path, _ = filepath.Abs("./bans/softbans.json")
	file, _ = ioutil.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &Softbans)
	path, _ = filepath.Abs("./bans/userbans.json")
	file, _ = ioutil.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &Userbans)
}

func SaveBans() {
	path, _ := filepath.Abs("./bans/bans.json")
	jsonText, _ := json.Marshal(BansS)
	_ = ioutil.WriteFile(path, jsonText, 0644)
	path, _ = filepath.Abs("./bans/softbans.json")
	jsonText, _ = json.Marshal(Softbans)
	_ = ioutil.WriteFile(path, jsonText, 0644)
	path, _ = filepath.Abs("./bans/userbans.json")
	jsonText, _ = json.Marshal(Userbans)
	_ = ioutil.WriteFile(path, jsonText, 0644)
}
