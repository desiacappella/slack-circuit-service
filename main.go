package main

import "fmt"

func main() {
	// slackRemoveOldOfficers()
	subms := jotformGetInactiveCaptains()

	var emails []string
	for _, s := range subms {
		emails = append(emails, s.Email)
		fmt.Println(s.Email)
	}

	// slackSendMessage(emails, farewellMsg)
}
