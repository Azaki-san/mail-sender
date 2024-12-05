package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	from := "secret.anry@gmail.com"
	smtpPort := "587"
	smtpHost := "smtp.gmail.com"
	password := "yezb ykbj uemy rlfk"
	to := []string{"leonandrey221@gmail.com"}
	auth := smtp.PlainAuth("", from, password, smtpHost)
	body := fmt.Sprintf("From user: %s\nMessage: %s", "kudasov", "message message allo allo asfaslfas")
	// xpaubddcdptvqage - yandex
	subject := "Subject: SA Course\n"
	fromHeader := "From: " + from + "\n"
	toHeader := "To: " + strings.Join(to, ", ") + "\n"

	// Construct the full email message
	msg := []byte(subject + fromHeader + toHeader + "\n" + body)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return
	}
	log.Println("Email sent successfully")
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
