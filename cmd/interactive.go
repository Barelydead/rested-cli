/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/barelydead/rested/handlers"
	"github.com/barelydead/rested/storage"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactively build or execute REST calls",
	Run: func(cmd *cobra.Command, args []string) {
		// Load DB
		db, err := storage.LoadFromFile(storage.DbFile)
		if err != nil {
			fmt.Printf("❌ Failed to load DB: %v\n", err)
			os.Exit(1)
		}

		// Auto-save on Ctrl+C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			db.SaveToFile(storage.DbFile)
			fmt.Println("\n💾 Data saved. Goodbye!")
			os.Exit(0)
		}()

		// Start the main app
		handlers.Root(db)

		// Save on exit
		err = db.SaveToFile(storage.DbFile)
		if err != nil {
			fmt.Printf("❌ Failed to save DB: %v\n", err)
		} else {
			fmt.Println("💾 DB saved.")
		}
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// interactiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// interactiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
