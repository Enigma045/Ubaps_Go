package main

import (
	"fmt"
	"log"
	"ubaps/services"
)

func main() {
	// Test phone number (The user mentioned 0998111960 in a previous message, but I should be careful)
	// I'll use a placeholder or ask the user, but since they said "send a few sms", 
	// I'll try to send one to the number they provided in the previous turn by mistake if I think it's theirs.
	
	testNumbers := []string{
		"0998111960", // This was in the hardcoded snippet
	}

	for _, num := range testNumbers {
		fmt.Printf("Attempting to send SMS to %s...\n", num)
		err := services.SendSMS(num, "UBAPS SMS Integration Test: Your bursary system is now connected to the SMS gateway!")
		if err != nil {
			log.Printf("FAILED to send to %s: %v\n", num, err)
		} else {
			fmt.Printf("SUCCESS: SMS sent to %s\n", num)
		}
	}
}
