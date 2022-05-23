package main

import (
	"context"
	"fmt"
	"os"

	onesignal "github.com/namnm1991/onesignal-go-api"
)

const (
	ONESIGNAL_APP_ID  = "d1c7ef44-9172-474e-a40c-98d0881da725"
	ONESIGNAL_API_KEY = "N2I5MmYwYjctM2Q5Ni00MzI1LTk0ZTQtYTBjMjJhYWJlODlh"
)

func sendEmail(emails []string, subject string, content string) {
	notification := *onesignal.NewNotification(ONESIGNAL_APP_ID)

	// ===========================================================================
	// Config email field
	notification.SetIncludeEmailTokens(emails)
	notification.SetEmailSubject(subject)
	notification.SetEmailBody(content)

	sendNoti(notification)
}

func sendNoti(notification onesignal.Notification) {
	configuration := onesignal.NewConfiguration()
	apiClient := onesignal.NewAPIClient(configuration)

	appAuth := context.WithValue(context.Background(), onesignal.AppAuth, ONESIGNAL_API_KEY)

	resp, r, err := apiClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateNotification``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	// response from `CreateNotification`: InlineResponse200
	fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateNotification`: %v\n", resp)
}
