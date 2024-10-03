package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up a test quiz.json file
	testQuiz := []byte(`{
		"quiz": {
			"questions": ["Test Question 1?", "Test Question 2?"],
			"answers": [
				["A", "B", "C", "D"],
				["W", "X", "Y", "Z"]
			]
		},
		"correctAnswers": [2, 3]
	}`)

	err := os.WriteFile("test_quiz.json", testQuiz, 0644)
	if err != nil {
		panic("Failed to create test quiz file")
	}

	// Load the test quiz
	if err := loadQuizFromFile("test_quiz.json"); err != nil {
		panic("Failed to load test quiz")
	}

	// Run tests
	code := m.Run()

	// Clean up
	os.Remove("test_quiz.json")

	os.Exit(code)
}

func TestParseAndValidateAnswers(t *testing.T) {
	tests := []struct {
		name    string
		input   PostQuiz
		want    []int
		wantErr bool
	}{
		{
			name:    "Valid input",
			input:   PostQuiz{Answers: []int{2, 3}},
			want:    []int{2, 3},
			wantErr: false,
		},
		{
			name:    "Invalid number of answers",
			input:   PostQuiz{Answers: []int{2}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid answer range",
			input:   PostQuiz{Answers: []int{4, 3}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/quiz/submit", bytes.NewReader(body))
			got, err := parseAndValidateAnswers(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAndValidateAnswers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalSlices(got, tt.want) {
				t.Errorf("parseAndValidateAnswers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostQuiz(t *testing.T) {
	// Reset global variables before testing
	totalSubmissions = 0
	scoreHistogram = make([]*int, len(correctAnswers)+1)
	for i := range scoreHistogram {
		scoreHistogram[i] = new(int)
	}

	tests := []struct {
		name           string
		input          PostQuiz
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "All correct answers",
			input:          PostQuiz{Answers: []int{2, 3}},
			expectedStatus: http.StatusOK,
			expectedBody:   "You got 2 correct answers. You performed better than 100.00% of users.",
		},
		{
			name:           "Some correct answers",
			input:          PostQuiz{Answers: []int{2, 0}},
			expectedStatus: http.StatusOK,
			expectedBody:   "You got 1 correct answers. You performed better than 0.00% of users.",
		},
		{
			name:           "Invalid input",
			input:          PostQuiz{Answers: []int{5, 5}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid answer for question 1: must be between 0 and 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/quiz/submit", bytes.NewReader(body))
			w := httptest.NewRecorder()
			postQuiz(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("postQuiz() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Trim whitespace from both actual and expected responses
			actualBody := strings.TrimSpace(w.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)

			if actualBody != expectedBody {
				t.Errorf("postQuiz() body = %v, want %v", actualBody, expectedBody)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
