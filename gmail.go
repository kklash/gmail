// Package gmail provides a lightweight API to send emails using the google SMTP server.
package gmail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

const (
	// https://support.google.com/a/answer/176600?hl=en
	address = "smtp.gmail.com"
	port    = 465 // not sure why google says 465 is SSL only and 587 is for TLS... seems to be the reverse.
)

// Type Login represents a set of gmail account credentials
// which can be used to log in and send email via SMTP.
type Login struct {
	sender   string // fooexample@gmail.com
	password string // app password
}

// NewLogin returns a pointer to a login struct
// containing the given email address and password.
func NewLogin(sender, password string) *Login {
	return &Login{sender, password}
}

// Send sends a plaintext email message to the given recipients.
func (login *Login) Send(recipients []string, subject, body string) error {
	if recipients == nil || len(recipients) == 0 {
		return fmt.Errorf("cannot supply nil or empty recipients slice")
	}

	tlsConfig := &tls.Config{
		ServerName: address,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", address, port), tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, address)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", login.sender, login.password, address)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(login.sender); err != nil {
		return err
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("From: %s\n", login.sender)
	message += fmt.Sprintf("To: %s\n", strings.Join(recipients, ";"))
	// message += fmt.Sprintf("CC: %s\n", strings.Join(ccRecipients, ";"))
	message += fmt.Sprintf("Subject: %s\n\n", subject)
	message += body
	if _, err := wc.Write([]byte(message)); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}
