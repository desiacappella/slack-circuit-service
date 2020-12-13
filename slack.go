package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

const slackToken = "slackUserToken"

const typePrivateChannel = "private_channel"

var api = slack.New(getToken(slackToken))

func slackGetChannelByName(channelName string) (slack.Channel, error) {
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}})
	if err != nil {
		epanic(err, "can't get user's channels")
	}

	for _, c := range channels {
		if c.Name == channelName {
			return c, nil
		}
	}

	return slack.Channel{}, fmt.Errorf("No channel found")
}

func slackGetDirectors() {
	channel, _ := slackGetChannelByName("comp-2021-directors")

	// Get all members
	userIds, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channel.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get channel's members")
	}

	// Use api.GetUsersInfo if needed
}

func slackGetOfficerEmails() {
	officerC, _ := slackGetChannelByName("circuit-officers")
	audioC, err := slackGetChannelByName("circuit-audio-production")
	if err != nil {
		epanic(err, "no find new chan")
	}

	newUsers, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: officerC.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get officers members")
	}

	oldUsers, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: audioC.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get officers members")
	}

	// Filter out from newUsers

	_, err = api.InviteUsersToConversation(newC.ID, newUsers...)
	if err != nil {
		epanic(err, "unable to invite users to conversation")
	}
}
