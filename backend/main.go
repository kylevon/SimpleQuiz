package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Quiz struct {
	Questions []string   `json:"questions"`
	Answers   [][]string `json:"answers"`
}

type PostQuiz struct {
	Answers []int `json:"answers"`
}

var quiz Quiz
var correctAnswers []int
var mu sync.Mutex

var (
	totalSubmissions int
	scoreHistogram   []*int // Slice to store counts for each possible score
)

func loadQuizFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var quizData struct {
		Quiz struct {
			Questions []string   `json:"questions"`
			Answers   [][]string `json:"answers"`
		} `json:"quiz"`
		CorrectAnswers []int `json:"correctAnswers"`
	}

	if err := json.Unmarshal(data, &quizData); err != nil {
		return err
	}

	quiz.Questions = quizData.Quiz.Questions
	quiz.Answers = quizData.Quiz.Answers
	correctAnswers = quizData.CorrectAnswers

	return nil
}

func init() {
	if err := loadQuizFromFile("quiz.json"); err != nil {
		fmt.Printf("Failed to load quiz: %v\n", err)
		os.Exit(1)
	}

	scoreHistogram = make([]*int, len(correctAnswers)+1)
	for i := range scoreHistogram {
		scoreHistogram[i] = new(int)
	}
}

func getQuiz(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(quiz)
}

// parseAndValidateAnswers checks that:
// 1. Input is valid JSON
// 2. Number of answers matches questions
// 3. Each answer is within valid range
func parseAndValidateAnswers(r *http.Request) ([]int, error) {
	var submitted PostQuiz
	if err := json.NewDecoder(r.Body).Decode(&submitted); err != nil {
		return nil, fmt.Errorf("invalid JSON input: %v", err)
	}

	if len(submitted.Answers) != len(correctAnswers) {
		return nil, fmt.Errorf("incorrect number of answers: expected %d, got %d", len(correctAnswers), len(submitted.Answers))
	}

	validatedAnswers := make([]int, len(submitted.Answers))
	for i, answer := range submitted.Answers {
		if answer < 0 || answer >= len(quiz.Answers[i]) {
			return nil, fmt.Errorf("invalid answer for question %d: must be between 0 and %d", i+1, len(quiz.Answers[i])-1)
		}
		validatedAnswers[i] = answer
	}

	return validatedAnswers, nil
}

func postQuiz(w http.ResponseWriter, r *http.Request) {
	validatedAnswers, err := parseAndValidateAnswers(r)
	if err != nil {
		// Remove the newline character from the error message
		errMsg := strings.TrimSpace(err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	score := 0
	for i, ans := range validatedAnswers {
		if ans == correctAnswers[i] {
			score++
		}
	}

	mu.Lock()
	totalSubmissions++
	*scoreHistogram[score]++
	betterThan := 0
	for i := 0; i < score; i++ {
		betterThan += *scoreHistogram[i]
	}

	var ranking float64
	if totalSubmissions == 1 {
		// If this is the first submission, set ranking based on score
		if score == len(correctAnswers) {
			ranking = 100.0
		} else {
			ranking = 0.0
		}
	} else {
		// Calculate ranking normally for subsequent submissions
		ranking = float64(betterThan) / float64(totalSubmissions-1) * 100
	}
	mu.Unlock()

	response := fmt.Sprintf("You got %d correct answers. You performed better than %.2f%% of users.", score, ranking)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/quiz", getQuiz)
	http.HandleFunc("/quiz/submit", postQuiz)
	fmt.Println("Quiz API running on :8080")
	http.ListenAndServe(":8080", nil)
}
