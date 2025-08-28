package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	Name       string `json:"name"`
	Activities []struct {
		Name          string `json:"name"`
		StartAmount   int    `json:"startAmount"`
		CurrentAmount int    `json:"currentAmount"`
	} `json:"activities"`
}

type Participant struct {
	ID       string    `json:"id"`
	Points   int       `json:"points"`
	Accounts []Account `json:"accounts"`
}

func main() {
	// Open the SQLite database directly
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Read the JSON file
	jsonData, err := os.ReadFile("participants.json")
	if err != nil {
		log.Fatalf("Failed to read participants.json: %v", err)
	}

	var participants []Participant
	if err := json.Unmarshal(jsonData, &participants); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Found %d participants to import\n", len(participants))

	// Get the first server ID from the database
	var serverID, serverName string
	err = db.QueryRow("SELECT id, name FROM server LIMIT 1").Scan(&serverID, &serverName)
	if err != nil {
		log.Fatalf("No servers found in database: %v", err)
	}
	fmt.Printf("Importing to server: %s (%s)\n", serverName, serverID)

	// Begin transaction for bulk insert
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Prepare statements for better performance
	participantStmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO participant (discord_id, server_id, botm_points, kots_points)
		VALUES (?, ?, ?, 0)`)
	if err != nil {
		log.Fatalf("Failed to prepare participant statement: %v", err)
	}
	defer participantStmt.Close()

	accountStmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO account (participant_id, username, failed_fetch_count)
		VALUES (?, ?, 0)`)
	if err != nil {
		log.Fatalf("Failed to prepare account statement: %v", err)
	}
	defer accountStmt.Close()

	// Bulk import participants
	successCount := 0
	errorCount := 0

	for i, participant := range participants {
		if i%100 == 0 {
			fmt.Printf("Processing participant %d/%d...\n", i+1, len(participants))
		}

		if err := importParticipantDirect(tx, participantStmt, accountStmt, serverID, participant); err != nil {
			fmt.Printf("❌ Failed to import participant %s: %v\n", participant.ID, err)
			errorCount++
		} else {
			successCount++
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("\n✅ Import completed!\n")
	fmt.Printf("Successfully imported: %d participants\n", successCount)
	fmt.Printf("Failed to import: %d participants\n", errorCount)
}

func importParticipantDirect(tx *sql.Tx, participantStmt, accountStmt *sql.Stmt, serverID string, participant Participant) error {
	// Insert participant with BOTM points
	_, err := participantStmt.Exec(participant.ID, serverID, participant.Points)
	if err != nil {
		return fmt.Errorf("failed to insert participant: %w", err)
	}

	// Insert each account for the participant
	for _, account := range participant.Accounts {
		// Clean the account name (remove any whitespace)
		accountName := strings.TrimSpace(account.Name)
		if accountName == "" {
			continue
		}

		// Insert the account
		_, err := accountStmt.Exec(participant.ID, accountName)
		if err != nil {
			// Log error but continue with other accounts
			fmt.Printf("    ⚠️  Failed to add account %s: %v\n", accountName, err)
			continue
		}
	}

	return nil
}
