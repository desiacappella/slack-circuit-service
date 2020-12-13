package main

import (
	"github.com/slack-go/slack"
)

const slackToken = "slackUserToken"

const typePrivateChannel = "private_channel"

var api = slack.New(getToken(slackToken))

func slackGetChannelMembers(channelName string) []string {
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

	// TODO check if channel was never found

	// 2. Get all members
	userIds, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channel.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get channel's members")
	}

	return userIds
}

func slackGetDirectors() {
	channelName := "comp-2021-directors"

	slackGetChannelMembers(channelName)
	// Use api.GetUsersInfo if needed
}

func slackGetOfficerEmails() {
	channelName := "circuit-officers"

	slackGetChannelMembers(channelName)

	// Get channel ID of audio-production

	// api.InviteUserToChannel()
}
