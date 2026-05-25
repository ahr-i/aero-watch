package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/ahr-i/aero-watch/ai-analysis/setting"
)

type openAIRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"`
}

type openAIMessage struct {
	Role    string          `json:"role"`
	Content []openAIContent `json:"content"`
}

type openAIContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type openAIResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (h *Handler) questionHandler(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid json body"})
		return
	}

	prompt := strings.TrimSpace(body.Prompt)
	user := strings.TrimSpace(body.User)
	if prompt == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "prompt is required"})
		return
	}
	if user == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "user is required"})
		return
	}

	answer, err := requestOpenAI(openAIRequest{
		Model: setting.Setting.OpenAI.Model,
		Input: joinPrompt(
			setting.Setting.OpenAI.TextPrePrompt,
			prompt,
			setting.Setting.OpenAI.TextPostPrompt,
		),
	})
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: err.Error()})
		return
	}

	if err := saveQuestionAnswer(user, prompt, answer); err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: err.Error()})
		return
	}

	rend.JSON(w, http.StatusOK, answerResponseBody{Answer: answer})
}

func (h *Handler) questionWithImageHandler(w http.ResponseWriter, r *http.Request) {
	var body imageRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid json body"})
		return
	}

	prompt := strings.TrimSpace(body.Prompt)
	user := strings.TrimSpace(body.User)
	image := strings.TrimSpace(body.ImageBase64)
	if prompt == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "prompt is required"})
		return
	}
	if user == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "user is required"})
		return
	}
	if image == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "imageBase64 is required"})
		return
	}

	imageURL, err := makeImageDataURL(image, body.ImageType)
	if err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: err.Error()})
		return
	}

	answer, err := requestOpenAI(openAIRequest{
		Model: setting.Setting.OpenAI.Model,
		Input: []openAIMessage{
			{
				Role: "user",
				Content: []openAIContent{
					{
						Type: "input_text",
						Text: joinPrompt(
							setting.Setting.OpenAI.ImagePrePrompt,
							prompt,
							setting.Setting.OpenAI.ImagePostPrompt,
						),
					},
					{
						Type:     "input_image",
						ImageURL: imageURL,
					},
				},
			},
		},
	})
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: err.Error()})
		return
	}

	if err := saveQuestionAnswer(user, prompt, answer); err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: err.Error()})
		return
	}

	rend.JSON(w, http.StatusOK, answerResponseBody{Answer: answer})
}

func requestOpenAI(requestBody openAIRequest) (string, error) {
	if setting.Setting.OpenAI.APIKey == "" {
		return "", errors.New("OPENAI_API_KEY is not set")
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(setting.Setting.OpenAI.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	req, err := http.NewRequest(http.MethodPost, setting.Setting.OpenAI.APIURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+setting.Setting.OpenAI.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: timeout}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var openAIRes openAIResponse
	if err := json.Unmarshal(responseBody, &openAIRes); err != nil {
		return "", err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		if openAIRes.Error != nil && openAIRes.Error.Message != "" {
			return "", fmt.Errorf("openai api error: %s", openAIRes.Error.Message)
		}
		return "", fmt.Errorf("openai api error: status %d", res.StatusCode)
	}

	answer := extractOpenAIAnswer(openAIRes)
	if answer == "" {
		return "", errors.New("openai api returned empty answer")
	}

	return answer, nil
}

func extractOpenAIAnswer(response openAIResponse) string {
	if strings.TrimSpace(response.OutputText) != "" {
		return strings.TrimSpace(response.OutputText)
	}

	var builder strings.Builder
	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				if builder.Len() > 0 {
					builder.WriteString("\n")
				}
				builder.WriteString(strings.TrimSpace(content.Text))
			}
		}
	}

	return builder.String()
}

func joinPrompt(prePrompt, prompt, postPrompt string) string {
	parts := make([]string, 0, 3)
	for _, part := range []string{prePrompt, prompt, postPrompt} {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, "\n\n")
}

func makeImageDataURL(imageBase64, mimeType string) (string, error) {
	if strings.HasPrefix(imageBase64, "data:") {
		return imageBase64, nil
	}

	if _, err := base64.StdEncoding.DecodeString(imageBase64); err != nil {
		return "", errors.New("imageBase64 must be a valid base64 string")
	}

	mimeType = strings.TrimSpace(mimeType)
	if mimeType == "" {
		mimeType = "image/png"
	}

	return "data:" + mimeType + ";base64," + imageBase64, nil
}

func saveQuestionAnswer(user, prompt, answer string) error {
	historyDir := strings.TrimSpace(setting.Setting.HistoryDir)
	if historyDir == "" {
		historyDir = "./questionHistory"
	}

	now := time.Now()
	userDir := filepath.Join(historyDir, safePathName(user))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return err
	}

	fileName := now.Format("20060102_150405.000000000") + ".txt"
	filePath := filepath.Join(userDir, fileName)
	content := fmt.Sprintf(
		"User: %s\nCreatedAt: %s\n\nQuestion:\n%s\n\nAnswer:\n%s\n",
		user,
		now.Format(time.RFC3339),
		prompt,
		answer,
	)

	return os.WriteFile(filePath, []byte(content), 0644)
}

func safePathName(value string) string {
	var builder strings.Builder
	for _, r := range strings.TrimSpace(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			builder.WriteRune(r)
			continue
		}
		builder.WriteRune('_')
	}

	if builder.Len() == 0 {
		return "unknown"
	}

	return builder.String()
}
