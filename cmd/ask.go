package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask Gemini a question",
	Long:  `Ask Gemini a question and get a command as a response.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a question to ask Gemini.")
			return
		}
		query := args[0]

		config, err := loadConfig()
		if err != nil {
			fmt.Println("Please execute gencmd config first.")
		}
		apiKey := config.APIKey
		if apiKey == "" {
			log.Fatal("API key not set. Use 'config' command to set the API key.")
			return
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatal(err)
			return
		}
		defer client.Close()

		model := client.GenerativeModel(config.Model)

		resp, err := model.GenerateContent(ctx, genai.Text(buildQuery(query)))
		if err != nil {
			log.Fatal(err)
			return
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			command := resp.Candidates[0].Content.Parts[0]
			fmt.Printf("%s", command)
		} else {
			fmt.Println("No response received from Gemini.")
		}
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func buildQuery(query string) string {
	tmplStr := `You are a command-line tool. Convert the following request into a command.
The request is: {{.Query}}
Make sure to generate commands tailored to the shell you will use ({{.Shell}}) on ({{.OS}}).
When a file path is required, using /path/to/filename.
Provide only the command as a response, without code blocks or additional formatting, so it can be copied and used immediately.
Ensure the command is written as a one-liner.
`
	tmpl, err := template.New("query").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	var result strings.Builder
	err = tmpl.Execute(&result, map[string]interface{}{
		"Query": query,
		"Shell": getShell(),
		"OS":    getOperationSystem(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return result.String()
}

func getShell() string {
	defaultShell := os.Getenv("SHELL")
	if defaultShell == "" {
		return "bash"
	}
	return filepath.Base(defaultShell)
}

func getOperationSystem() string {
	os := runtime.GOOS

	switch os {
	case "windows":
		return "Windows"
	case "darwin":
		return "MacOS"
	case "linux":
		return "Linux"
	default:
		return "MacOS"
	}
}
