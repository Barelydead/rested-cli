package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type DB struct {
	Collections []Collection `json:"collections"`
}

type RestedRequest struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	RequestName string            `json:"requestName"`
}

type Collection struct {
	Title    string          `json:"title"`
	Requests []RestedRequest `json:"requests"`
}

const DbFile = "rested_data.json"

// SaveToFile serializes DB to a JSON file
func (db *DB) SaveToFile(path string) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal DB: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write DB to file: %w", err)
	}

	return nil
}

// LoadFromFile reads DB state from a JSON file
func LoadFromFile(path string) (*DB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty DB if file doesn't exist yet
			return &DB{Collections: []Collection{}}, nil
		}
		return nil, fmt.Errorf("failed to read DB file: %w", err)
	}

	var db DB
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DB: %w", err)
	}

	return &db, nil
}
