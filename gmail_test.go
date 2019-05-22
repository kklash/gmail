package gmail

import (
	"os"
	"testing"
)

// To run the test, set the following environment variables:
//  GMAIL_ADDRESS=senderexample@gmail.com
//  GMAIL_PASSWORD=mysecretpassword
//  GMAIL_RECIPIENT=someotherguy@hotmail.com
func TestSend(t *testing.T) {
	sender := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_PASSWORD")
	recipient := os.Getenv("GMAIL_RECIPIENT")

	if len(sender) == 0 || len(password) == 0 || len(recipient) == 0 {
		t.Errorf(
			"Testing environment variables are not defined.\n" +
				"Set GMAIL_ADDRESS (sender), GMAIL_PASSWORD, and GMAIL_RECIPIENT environment variables.",
		)
		return
	}

	login := NewLogin(sender, password)
	err := login.Send([]string{recipient}, "My test email", "this is a body of text. are you reading this?")
	if err != nil {
		t.Errorf("failed to send email: %s", err)
	}
}
