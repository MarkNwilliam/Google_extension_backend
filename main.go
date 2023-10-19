package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	go_openai "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/chains"
	langchaingo_openai "github.com/tmc/langchaingo/llms/openai"
)

func main() {

	//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

	apiToken := os.Getenv("OPENAI_API_KEY")
	if apiToken == "" {
		panic("Missing OPENAI_API_KEY environment variable")
	}

	http.HandleFunc("/api/gpt", gptHandler)
	http.HandleFunc("/api/connect", connectHandler)
	http.HandleFunc("/api/langchaingo", langchaingoHandler)
	http.HandleFunc("/api/math-agent", mathAgentHandler)
	http.ListenAndServe(":8080", nil)
}

type RequestBody struct {
	Prompt string `json:"prompt"`
}

type ResponseData struct {
	Message string `json:"message"`
}

func gptHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	// Decode the request body (expects JSON with a "prompt" field)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	client := go_openai.NewClient("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") // replace with your OpenAI API key

	resp, err := client.CreateChatCompletion(
		context.Background(),
		go_openai.ChatCompletionRequest{
			Model: go_openai.GPT3Dot5Turbo,
			Messages: []go_openai.ChatCompletionMessage{
				{
					Role:    go_openai.ChatMessageRoleUser,
					Content: requestBody.Prompt,
				},
			},
		},
	)

	if err != nil {
		http.Error(w, "Error processing the request", http.StatusInternalServerError)
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	responseData := ResponseData{
		Message: resp.Choices[0].Message.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	responseData := ResponseData{
		Message: "connected",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func langchaingoHandler(w http.ResponseWriter, r *http.Request) {
	apiToken := os.Getenv("OPENAI_API_KEY")
	var requestBody RequestBody

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Use langchaingo to interact with OpenAI

	llm, err := langchaingo_openai.New(langchaingo_openai.WithToken(apiToken))
	if err != nil {
		responseData := ResponseData{
			Message: "Error initializing langchaingo",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseData)
		return
	}

	completion, err := llm.Call(context.Background(), requestBody.Prompt)
	if err != nil {
		responseData := ResponseData{
			Message: "Error processing the request with langchaingo",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseData)
		return
	}

	responseData := ResponseData{
		Message: completion,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func mathAgentHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	llm, err := langchaingo_openai.New()
	if err != nil {
		http.Error(w, "Error initializing OpenAI for math", http.StatusInternalServerError)
		return
	}

	llmMathChain := chains.NewLLMMathChain(llm)
	ctx := context.Background()
	result, err := chains.Run(ctx, llmMathChain, requestBody.Prompt)

	if err != nil {
		http.Error(w, "Error processing math query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := ResponseData{
		Message: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
