package handlers

import (
	"fmt"
	"log"

	"github.com/barelydead/rested/storage"
	"github.com/manifoldco/promptui"
)

func ShowCollection(db *storage.DB, index int) {
	collection := db.Collections[index]

	reqDisplay := []string{}
	for _, r := range collection.Requests {
		reqDisplay = append(reqDisplay, fmt.Sprintf("%s - %s", r.Method, r.RequestName))
	}

	if len(reqDisplay) == 0 {
		fmt.Println("⚠️ No requests in this collection.")
		return
	}

	prompt := promptui.Select{
		Label: collection.Title,
		Items: reqDisplay,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	NewRequest(db, collection.Requests[index])
}

func OpenCollection(db *storage.DB) {
	for {
		items := []string{}
		for _, c := range db.Collections {
			items = append(items, c.Title)
		}
		items = append(items, "+ New collection", "Back")

		prompt := promptui.Select{
			Label: "Choose collection",
			Items: items,
		}

		index, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		switch result {
		case "+ New collection":
			AddCollection(db)
		case "Back":
			Root(db)
		default:
			ShowCollection(db, index)
		}
	}
}

func AddCollection(db *storage.DB) {
	prompt := promptui.Prompt{
		Label: "Collection name",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("collection name can't be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Printf("Collection creation cancelled: %v", err)
		return
	}

	newCollection := storage.Collection{
		Title:    result,
		Requests: []storage.RestedRequest{},
	}

	db.Collections = append(db.Collections, newCollection)
	fmt.Printf("✅ Collection '%s' added!\n\n", result)
}
