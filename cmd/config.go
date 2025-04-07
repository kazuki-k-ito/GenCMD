package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Gemini API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			fmt.Println("API key cannot be empty.")
			return
		}

		fmt.Print("Enter gemini model. default:(gemini-2.0-flash): ")
		model, _ := reader.ReadString('\n')
		model = strings.TrimSpace(model)
		if model == "" {
			model = "gemini-2.0-flash"
		}

		err := saveConfig(apiKey, model)
		if err != nil {
			fmt.Println("Failed to save config:", err)
			return
		}

		fmt.Println("API key and model set successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
