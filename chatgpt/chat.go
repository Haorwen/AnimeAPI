package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/google/uuid"
)

type ChatGPT struct {
	Auth string
}

func NewChatGPT(auth string) *ChatGPT {
	return &ChatGPT{Auth: auth}
}

func (c *ChatGPT) GetChatResponse(prompt string) (string, error) {
	type requestBody struct {
		Prompt  string   `json:"prompt"`
		MaxToke int      `json:"max_tokens"`
		Model   string   `json:"model"`
		Tokens  []string `json:"tokens"`
	}

	url := "https://api.openai.com/v1/completions"
	requestBody := requestBody{
		Prompt:  prompt,
		MaxToke: 2048,
		Model:   "text-davinci-002",
		Tokens:  []string{"|"},
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth)

	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	type responseBody struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}

	var response responseBody
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return response.Choices[0].Text, nil
}

func main() {
	chatGPT := NewChatGPT("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")

	response, err := chatGPT.GetChatResponse("What is the capital of France?")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
}
