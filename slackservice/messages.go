package slackservice

import (
	"fmt"
	"strings"
)

var farewellMsg = strings.Join([]string{
	"Hi there! Our records show that you graduated from your team in 2020.",
	"If this information is incorrect and you are still on your team, please reply with your updated graduation year in the *next 72 hours*.",
	"If you have already graduated, we will move you to the #circuit-alumni channel and deactivate your account after 4 months of inactivity.",
	"Thank you!",
}, " ")

func mirchiMsg(link string) string {
	return strings.Join([]string{
		"<!channel> :rotating_light:$$$ ALERT:rotating_light: WHAT'S UP ASA TEAMS!! I'm sure by now you have heard of Mirchi :hot_pepper:, a dating/friendship app for South Asians! Mirchi has agreed to sponsor us for this season, and we're sooo excited to get y'all involved!",
		"Starting today until February 28th, we want you guys to get all of your teammates, friends, and fam to create profiles for Mirchi. Each team that gets at or above 210 unique profile creations will be entered into a raffle to win...drumroll please...$1000!! :money_mouth_face::gungho:",
		fmt.Sprintf("Your team’s unique link to share and download the Mirchi app is: %s", link),
		"Please do lmk if you have any questions and goooood luck :))",
	}, "\n\n")
}

var extensionMsg = strings.Join([]string{
	"Hi <!channel> - hope you are doing well and hard at work on your second arrangement's performance video! I wanted to share some updates regarding the submission for Series 2.",
	"As I'm sure you all know, Texas has been experiencing an ice and snowstorm causing extended power outages, lack of water, connectivity, and more. As a result, some of the Texas teams have been severely impacted and been unable to work on their performance videos at all over the last week and more.",
	"In order to give them a fair amount of time to work on their videos, they have been granted an extension to submit their Series 2 videos by *March 5, 11:59PM PST*. We believe this is a reasonable time extension for teams who have been negatively impacted by circumstances outside their control.",
	"If you have any questions or concerns, feel free to send them here and I would be happy to chat further about this situation. Thank you for being understanding, and good luck!",
}, " ")

var extension2Msg = "Just to clarify, the deadline for you is still *March 1, 11:59PM PST* - the extended deadline is only for the teams in Texas who were impacted and have reached out to us! Thank you for understanding! :)"
