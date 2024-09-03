package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	journal "github.com/49pctber/journal/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var config journal.Config
		if config.LoadFromFile() != nil {
			config.Configure()
		}

		date_string := time.Now().Format("2006-01-02")
		entryfname := filepath.Join(config.EntryDirectory, date_string)

		// Open the file in Vim
		c := exec.Command(config.PreferredEditor, entryfname)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr

		// Run the command and wait for it to finish
		err := c.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Perhaps you need to update your configuration file.")
			return
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.journal.yaml)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
