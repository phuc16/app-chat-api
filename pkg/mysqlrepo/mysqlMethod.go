package mysqlrepo

import (
	"chat-api/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func RegisterNewUser(username, password string) error {
	// Prepare the SQL statement to insert a new user record
	query := `
		INSERT INTO users (username, password)
		VALUES (?, ?)
	`

	// Execute the SQL statement to insert the new user record
	_, err := dbClient.db.Exec(query, username, password)
	if err != nil {
		return err
	}

	return nil
}

func IsUserExist(username string) (bool, error) {
	// Prepare the SQL statement to count the number of rows with the given username
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE username = ?
	`

	// Execute the SQL query
	var count int
	err := dbClient.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	// If count is greater than 0, the user exists; otherwise, the user does not exist
	return count > 0, nil
}

func IsUserAuthentic(username, password string) error {
	// Prepare the SQL statement to retrieve the password for the given username
	query := `
		SELECT password
		FROM users
		WHERE username = ?
	`

	// Execute the SQL query
	var storedPassword string
	err := dbClient.db.QueryRow(query, username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid username")
		}
		return err
	}

	// Check if the stored password matches the provided password
	if storedPassword != password {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func UpdateContactList(username, contact string) error {
	// Prepare the SQL statement for inserting a new contact record
	query := `
		INSERT INTO contact_lists (username, contact_id, last_interaction_timestamp)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE last_interaction_timestamp = VALUES(last_interaction_timestamp)
	`

	// Execute the SQL statement to insert the new contact record
	_, err := dbClient.db.Exec(query, username, contact, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func CreateChat(c *model.Chat) (string, error) {
	// Prepare the SQL statement for inserting a new chat record
	query := `
		INSERT INTO chats (from_user, to_user, message, timestamp)
		VALUES (?, ?, ?, ?)
	`

	// Serialize chat data to JSON
	chatJSON, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	// Execute the SQL statement to insert the new chat record
	result, err := dbClient.db.Exec(query, c.From, c.To, string(chatJSON), c.Timestamp)
	if err != nil {
		return "", err
	}

	// Get the auto-generated chat ID from the result
	chatID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	// Convert the chat ID to string
	chatKey := fmt.Sprintf("chat#%d", chatID)

	// Update contact lists for both users
	err = UpdateContactList(c.From, c.To)
	if err != nil {
		return "", err
	}

	err = UpdateContactList(c.To, c.From)
	if err != nil {
		return "", err
	}

	return chatKey, nil
}

func CreateFetchChatBetweenIndex() error {
	// Perform MySQL query to create index
	_, err := dbClient.db.Exec("CREATE INDEX chat_index ON chats (from_user, to_user, timestamp)")
	if err != nil {
		return err
	}

	fmt.Println("Chat index created successfully in MySQL!")
	return nil
}

func FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	// Prepare the SQL query
	query := `
		SELECT *
		FROM chats
		WHERE (from_user = ? AND to_user = ? OR from_user = ? AND to_user = ?)
			AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
	`

	// Execute the SQL query
	rows, err := dbClient.db.Query(query, username1, username2, username2, username1, fromTS, toTS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and deserialize chats
	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		err := rows.Scan(&chat.ID, &chat.FromUser, &chat.ToUser, &chat.Message, &chat.Timestamp)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func FetchContactList(username string) ([]model.ContactList, error) {
	// Prepare the SQL query
	query := `
		SELECT contact_id, contact_name, last_interaction_timestamp
		FROM contact_lists
		WHERE username = ?
		ORDER BY last_interaction_timestamp DESC
	`

	// Execute the SQL query
	rows, err := dbClient.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and deserialize contact list
	var contactList []model.ContactList
	for rows.Next() {
		var contact model.ContactList
		err := rows.Scan(&contact.ContactID, &contact.ContactName, &contact.LastInteractionTimestamp)
		if err != nil {
			return nil, err
		}
		contactList = append(contactList, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contactList, nil
}
