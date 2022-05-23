package main

import (
	"fmt"
	"math/rand"
)

func main() {

	// ======================================================
	// send an email
	emails := []string{"nam@krystal.app"}
	subject := "Welcome to Krystal SmartAlert"
	content := fmt.Sprintf("S.O.S %d", rand.Intn(100))
	sendEmail(emails, subject, content)
}
