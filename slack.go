package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

const typePrivateChannel = "private_channel"

// Operate as a user
var userAPI = slack.New(getToken("slackUserToken"))

// Operate as a bot
var botAPI = slack.New(getToken("slackBotToken"))

func slackGetChannelByName(channelName string) (slack.Channel, error) {
	channels, _, err := userAPI.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}, Limit: 250})
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
	userIds, _, err := userAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channel.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get channel's members")
	}

	// Use api.GetUsersInfo if needed
	fmt.Println(userIds)
}

func slackGetOfficerEmails() {
	officerC, _ := slackGetChannelByName("circuit-officers")
	newUsers, _, err := userAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: officerC.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get officers members")
	}

	audioC, _ := slackGetChannelByName("circuit-audio-production")
	oldUsers, _, err := userAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: audioC.ID, Limit: 200})
	if err != nil {
		epanic(err, "can't get audio members")
	}

	// Filter out from newUsers
	usersToAdd := []string{}
outer:
	for _, u := range newUsers {
		for _, u2 := range oldUsers {
			if u == u2 {
				continue outer
			}
		}
		usersToAdd = append(usersToAdd, u)
	}

	fmt.Println(strings.Join(usersToAdd, ","))

	// _, err = api.InviteUsersToConversation(audioC.ID, usersToAdd...)
	// if err != nil {
	// 	epanic(err, "unable to invite users to conversation")
	// }
}

func slackCollectChannels() (allChannels []slack.Channel) {
	var cursor string
	for {
		var channels []slack.Channel
		var err error
		channels, cursor, err = userAPI.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}, Cursor: cursor})
		if err != nil {
			epanic(err, "Can't get channels for user")
		}
		allChannels = append(allChannels, channels...)

		if len(cursor) == 0 {
			break
		}
	}
	return allChannels
}

func slackTeamChannels(teams []Team, create bool) {
	// First check in unusual
	// Random heuristics:
	// - stop once you see a column
	// - remove A Cappella

	allChannels := slackCollectChannels()
	emailToIds := make(map[string]string)

	// Returns true if newly added, false if already existing
	addEmailToChannel := func(channel slack.Channel, email string) string {
		if len(email) == 0 {
			return "invalid email"
		}

		liaisonID, ok := emailToIds[email]
		if !ok {
			user, err := userAPI.GetUserByEmail(email)
			if err != nil {
				if err.Error() == "users_not_found" {
					return "NEED TO ADD " + email
				}
				epanic(err, "Can't get user by email "+email)
			}

			emailToIds[email] = user.ID
			liaisonID = user.ID
		}

		// Is the liaison already in?
		members, _, err := userAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channel.ID})
		if err != nil {
			epanic(err, "can't list members")
		}
		found := false
		for _, mem := range members {
			if mem == liaisonID {
				found = true
				break
			}
		}
		if found {
			return "already there"
		}

		_, err = userAPI.InviteUsersToConversation(channel.ID, liaisonID)
		if err != nil {
			epanic(err, "Can't invite member")
		}
		return "newly added"
	}

	numCreated := 0
	numNew := 0

	for _, team := range teams {
		channelName := "team-" + teamToID(team.Name)

		var channel slack.Channel
		for _, _channel := range allChannels {
			if _channel.Name == channelName {
				channel = _channel
				break
			}
		}

		var _ string
		if len(channel.ID) == 0 {
			if create {
				// Create channels
				channel, err := userAPI.CreateConversation(channelName, true)
				time.Sleep(time.Second)
				if err != nil {
					_ = err.Error()
				} else {
					numNew++
					_ = "created successfully " + channel.ID
				}
			} else {
				numNew++
				_ = "doesn't exist"
			}
		} else {
			numCreated++
			_ = "already created"
		}

		// Add liaison
		addEmailToChannel(channel, team.Liaison.Email)

		// Add captain
		addEmailToChannel(channel, team.Captain.Email)

		// Add officers
		for _, o := range team.Officers {
			addEmailToChannel(channel, o.Email)
			// fmt.Printf("[%s] Added %s: %s\n", channelName, o.Name, status)
		}
		if len(team.Officers) == 0 {
			fmt.Println(team.Name)
		}

		// fmt.Printf("[%s] %s\n", channelName, channelStatus)
		// fmt.Printf("%-25v | %-25v | %-25v | %-25v | %-25v\n", team.Name, channelName, channelStatus,
		// 	fmt.Sprintf("added %s (new: %v)", team.Liaison.Name, addedLiai), fmt.Sprintf("added %s (res: %s)", team.Captain.Name, addCapResult))
	}

	fmt.Println("NEW CHANNELS", numNew, "ALREADY CREATED CHANNELS", numCreated)
}

func slackSendMessage(emails []string, message string) {
	for _, email := range emails {
		user, err := userAPI.GetUserByEmail(email)
		if err != nil {
			if err.Error() == "users_not_found" {
				fmt.Printf("Unable to find user %s. Moving on...\n", email)
				continue
			} else {
				epanic(err, "WHO AM I????")
			}
		}

		msg := slack.MsgOptionText(message, true)

		_, _, err = userAPI.PostMessage(user.ID, msg)
		if err != nil {
			epanic(err, "can't send message")
		}

		fmt.Println("Sent message to", email)
	}
}

func slackSendToChannel(name string, message string) {
	channel, err := slackGetChannelByName(name)
	if err != nil {
		epanic(err, "Can't find channel")
	}

	block := slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, message, false, false), nil, nil)
	msg := slack.MsgOptionBlocks(block)

	// msg := slack.MsgOptionText(message, true)

	_, _, err = userAPI.PostMessage(channel.ID, msg)
	if err != nil {
		epanic(err, "can't send message")
	}

	fmt.Println("Sent message to", name)
}

// TODO: Remove the given officers from #circuit-officers and add to #circuit-alumni
func slackRemoveOldOfficers(subms []jotformSubm) {

}

func slackAddToAllChannels(prefix string, email string) {
	channels, _, err := userAPI.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}, Limit: 250})
	if err != nil {
		epanic(err, "can't get user's channels")
	}

	user, err := userAPI.GetUserByEmail(email)
	if err != nil {
		epanic(err, "can't find user")
	}

	for _, c := range channels {
		if strings.HasPrefix(c.Name, prefix) {
			_, err = userAPI.InviteUsersToConversation(c.ID, user.ID)
			if err != nil {
				if err.Error() == "already_in_channel" {
					fmt.Println(c.Name, "already added")
					continue
				} else {
					epanic(err, "can't invite")
				}
			}

			fmt.Println(c.Name, "added")
		}
	}
}
