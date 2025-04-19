package handlers

import (
	"fmt"
	"os"

	"github.com/barelydead/rested/storage"
	"github.com/manifoldco/promptui"
)

func Root(db *storage.DB) {
	prompt := promptui.Select{
		Label: "Select Action",
		Items: []string{"New request", "Open collection", "Exit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "New request":
		NewRequest(db, storage.RestedRequest{})
	case "Open collection":
		OpenCollection(db)
	case "Exit":
		db.SaveToFile(storage.DbFile)
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}
}
