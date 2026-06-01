package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiURL = "https://api.groq.com/openai/v1/chat/completions"
	model  = "openai/gpt-oss-20b"
)

// Client wraps the Groq REST API (OpenAI-compatible).
type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`

	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (c *Client) Complete(ctx context.Context, system, userPrompt string) (string, error) {
	messages := []Message{}
	if system != "" {
		messages = append(messages, Message{Role: "system", Content: system})
	}
	messages = append(messages, Message{Role: "user", Content: userPrompt})

	reqBody := request{
		Model:       model,
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.3,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq api returned status %d: %s", resp.StatusCode, string(rawBody))
	}

	var groqResp response
	if err := json.Unmarshal(rawBody, &groqResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("groq api error [%s]: %s", groqResp.Error.Type, groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty choices in groq response")
	}

	return groqResp.Choices[0].Message.Content, nil
}

func (c *Client) AnalyzeResume(ctx context.Context, resumeText string) (string, error) {
	system := `You are an expert resume reviewer and career coach.
Analyze the provided resume text and return ONLY valid JSON — no markdown fences, no explanation, no extra text.

JSON format:
{
  "overall_score": <integer 0-100>,
  "summary_feedback": "<overall feedback>",
  "sections": [
    { "name": "<section>", "score": <0-100>, "feedback": "<specific feedback>" }
  ],
  "suggestions": ["<actionable suggestion>"],
  "keywords": ["<relevant keyword>"]
}

Sections to evaluate (if present): Experience, Education, Skills, Summary, Certifications.
Be specific, honest, and constructive. Use strict but fair scoring.`

	prompt := fmt.Sprintf("Analyze this resume:\n\n%s", resumeText)

	return c.Complete(ctx, system, prompt)
}
