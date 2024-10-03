package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var postQuizCmd = &cobra.Command{
	Use:   "post [answers...]",
	Short: "Submit quiz answers",
	Long:  "Submit quiz answers. Provide answers as space-separated integers.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse arguments into integers
		answers := make([]int, len(args))
		for i, arg := range args {
			answer, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Error: Invalid answer '%s'. Please provide integers.\n", arg)
				return
			}
			answers[i] = answer
		}

		// Prepare the request body
		requestBody := map[string]interface{}{
			"answers": answers,
		}
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println("Error preparing request:", err)
			return
		}

		// Send the request
		resp, err := http.Post("http://localhost:8080/quiz/submit", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read and print the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println("Response:", string(body))

		// Print the submitted answers
		fmt.Printf("Submitted answers: %v\n", answers)
	},
}

func init() {
	rootCmd.AddCommand(postQuizCmd)
}
