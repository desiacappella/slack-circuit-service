package main

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
