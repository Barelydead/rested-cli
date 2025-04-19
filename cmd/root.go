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

var rootCmd = &cobra.Command{
	Use:   "rested",
	Short: "An interactive REST CLI",
	Run: func(cmd *cobra.Command, args []string) {
		// Load DB
		db, err := storage.LoadFromFile(storage.DbFile)
		if err != nil {
			fmt.Printf("❌ Failed to load DB: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("starting?")

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rested.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
