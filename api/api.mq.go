package api

import (
	"encoding/json"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"smartrb.com/auth"
	"smartrb.com/db"
	"smartrb.com/mq"
	"smartrb.com/users"
)

func getTopicKeys() []string {
	// Get topic keys from all domains
	auth_topic_keys := auth.GetAuthTopicKeys()

	// Pre-allocate slice with exact capacity needed
	totalLen := len(auth_topic_keys)
	all_topic_keys := make([]string, 0, totalLen)

	// Append all slices
	all_topic_keys = append(all_topic_keys, auth_topic_keys...)

	return all_topic_keys
}

func createMQSubscriber() {
	go func() {
		all_topic_keys := getTopicKeys()

		// handleMessages is imported from api.handlers.go file
		mq.Subscribe("events_topic", all_topic_keys, handleMessages) // This should block forever
	}()
}

func handleMessages(topic string, message string) {
	if strings.Contains(topic, "auth.") {
		handleAuthMessages(topic, message)
	}
	// if strings.Contains(topic, "platform_users.") {
	// 	handlePlatformUsersMessages(topic, message)
	// }
	// if strings.Contains(topic, "organizations.") {
	// 	handleOrganizationsMessages(topic, message)
	// }
	// if strings.Contains(topic, "schools.") {
	// 	handleSchoolsMessages(topic, message)
	// }
	// if strings.Contains(topic, "courses.") {
	// 	handleCoursesMessages(topic, message)
	// }
}

func handleAuthMessages(topic string, message string) {
	if topic == "auth.signup" {
		var auth_signup_message struct {
			Email  string `json:"email"`
			UserId string `json:"user_id"`
		}

		err := json.Unmarshal([]byte(message), &auth_signup_message)
		if err != nil {
			log.Printf("Error parsing JSON message: %v", err)
			return
		}

		users.CreateUser(auth_signup_message.Email, auth_signup_message.UserId, db.Default(), mq.Publish)

	}
}
