package mail

import (
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendMailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping email test in short mode.")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "Hello, World!"
	content := `
	<h1>Hello, World!</h1>
	<p>This is a test email from <b>Golang Bank App</b>.</p>
	`
	to := []string{"vukyanh001@gmail.com"}
	attachFiles := []string{"../sample/mail.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
