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
	channels, _, err := userAPI.GetConversations(&slack.GetConversationsParameters{Types: []string{typePrivateChannel}})
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

// Returns true if newly added, false if already existing
func addEmailToChannel(channel slack.Channel, email string, emailToIds map[string]string) string {
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
	if !found {
		_, err := userAPI.InviteUsersToConversation(channel.ID, liaisonID)
		if err != nil {
			epanic(err, "Can't invite liaison")
		}
	}

	if found {
		return "already there"
	}
	return "newly added"
}

func slackTeamChannels(teams []Team) {
	// First check in unusual
	// Random heuristics:
	// - stop once you see a column
	// - remove A Cappella

	allChannels := slackCollectChannels()
	emailToIds := make(map[string]string)

	fmt.Printf("%-25v | %-25v | %-25v | %-25v | %-25v\n", "Team Name", "Channel Name", "Status", "Liaison", "Captain")
	for _, team := range teams {
		channelName, ok := unusualMapping[team.Name]
		if !ok {
			channelName = strings.ReplaceAll(strings.TrimSpace(strings.Split(strings.ReplaceAll(strings.ToLower(team.Name), "a cappella", ""), ":")[0]), " ", "-")
		}
		channelName = "team-" + channelName

		var channel slack.Channel
		for _, channel = range allChannels {
			if channel.Name == channelName {
				break
			}
		}

		var channelStatus string
		if len(channel.ID) == 0 {
			// Create channels
			channel, err := userAPI.CreateConversation(channelName, true)
			time.Sleep(time.Second)
			if err != nil {
				channelStatus = err.Error()
			} else {
				channelStatus = "created successfully " + channel.ID
			}
		} else {
			channelStatus = "already created"
		}

		// Add liaison
		addedLiai := false // addEmailToChannel(channel, team.Liaison.Email, emailToIds)

		// Add captain
		addCapResult := addEmailToChannel(channel, team.Captain.Email, emailToIds)

		fmt.Printf("%-25v | %-25v | %-25v | %-25v | %-25v\n", team.Name, channelName, channelStatus,
			fmt.Sprintf("added %s (new: %v)", team.Liaison.Name, addedLiai), fmt.Sprintf("added %s (res: %s)", team.Captain.Name, addCapResult))
	}
}

var farewellMsg = strings.Join([]string{
	"Hi there! Our records show that you graduated from your team in 2020.",
	"If this information is incorrect and you are still on your team, please reply with your updated graduation year in the *next 72 hours*.",
	"If you have already graduated, we will move you to the #circuit-alumni channel and deactivate your account after 4 months of inactivity.",
	"Thank you!",
}, " ")

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

// Remove the given officers from #circuit-officers and add to #circuit-alumni
func slackRemoveOldOfficers(subms []jotformSubm) {

}
