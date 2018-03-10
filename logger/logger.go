package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"

	"github.com/qor/mailer"
)

type Sender struct {
	*Config
}

type Config struct {
	Output io.Writer
}

func New(config *Config) *Sender {
	if config == nil {
		config = &Config{}
	}

	if config.Output == nil {
		config.Output = os.Stderr
	}

	return &Sender{Config: config}
}

func (sender *Sender) Send(email mailer.Email) error {
	var result = new(bytes.Buffer)

	formatAddress := func(key string, addresses []mail.Address) {
		var emails []string

		if len(addresses) > 0 {
			fmt.Fprintf(result, "%v: ", key)

			for _, address := range addresses {
				emails = append(emails, address.String())
			}

			fmt.Fprintf(result, "%s\n", strings.Join(emails, ", "))
		}
	}

	formatAddress("TO", email.TO)
	formatAddress("CC", email.CC)
	formatAddress("BCC", email.BCC)

	if email.From != nil {
		formatAddress("From", []mail.Address{*email.From})
	}

	if email.ReplyTo != nil {
		formatAddress("ReplyTO", []mail.Address{*email.ReplyTo})
	}

	if email.Subject != "" {
		fmt.Fprintf(result, "Subject: %v\n", email.Subject)
	}

	if email.Headers != nil {
		for key, value := range email.Headers {
			fmt.Fprintf(result, "%v: %v\n", key, value)
		}
	}

	for _, attachment := range email.Attachments {
		if attachment.Inline {
			fmt.Fprintf(result, "\nContent-Disposition: inline; filename=\"%v\"\n\n", attachment.FileName)
		} else {
			fmt.Fprintf(result, "\nContent-Disposition: attachment; filename=\"%v\"\n\n", attachment.FileName)
		}
	}

	if email.Text != "" {
		fmt.Fprintf(result, "\nContent-Type: text/plain; charset=UTF-8\n%v\n", email.Text)
	}

	if email.HTML != "" {
		fmt.Fprintf(result, "\nContent-Type: text/html; charset=UTF-8\n%v\n", email.HTML)
	}


	_, err := io.Copy(sender.Output, result)

	return err
}
