package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/slack-go/slack"
)

// This token needs to be acquired after authenticating with the app, on the app's home page
func getToken() string {
	b, err := ioutil.ReadFile("userToken")
	if err != nil {
		panic("Unable to find user token")
	}
	return string(b)
}

const channelName = "comp-2021-directors"
const typePrivateChannel = "private_channel"

func epanic(err error, msg string) {
	panic(fmt.Sprintf("%v, %s", err, msg))
}

func main() {
	token := strings.TrimSpace(getToken())
	api := slack.New(token)

	// 1. Find comp-2021-directors
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}})
	if err != nil {
		epanic(err, "can't get user's channels")
	}

	var channel slack.Channel
	for _, c := range channels {
		if c.Name == channelName {
			channel = c
			break
		}
	}

	// 2. Get all members
	userIds, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channel.ID})
	if err != nil {
		epanic(err, "can't get channel's members")
	}

	// 3. Evaluate which ones are not full members
	users, err := api.GetUsersInfo(userIds...)
	if err != nil {
		epanic(err, "can't get users info")
	}

	for _, user := range *users {
		if user.IsRestricted || user.IsUltraRestricted {
			fmt.Printf("user [%s] %s email %s is a guest\n", user.ID, user.Name, user.Profile.Email)
		}
	}
}
