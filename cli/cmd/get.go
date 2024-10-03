package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var getQuizCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch the quiz questions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8080/quiz")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(getQuizCmd)
}
