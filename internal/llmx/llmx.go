package llmx

import (
	"encoding/json"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkoukk/tiktoken-go"

	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
)

type LLMx struct {
	tokens   int
	messages []api.Message
}

func NewLLMx() *LLMx {
	return &LLMx{tokens: 0, messages: []api.Message{}}
}

func (l *LLMx) HandleUserMessage(msg, apiURL, model, apiKey string) (string, error) {
	logger.Log(logger.INFO, "handling user message (length: %d chars)", len(msg))
	// Create system prompt
	sysPrompt, err := createSystemPrompt()
	if err != nil {
		return "", err
	}
	if len(l.messages) > 0 {
		l.messages[0].Content = sysPrompt
	} else {
		l.messages = append(l.messages, api.Message{Role: "system", Content: sysPrompt})
	}
	// Add user message to messages
	userMsg := api.Message{
		Role:    "user",
		Content: msg,
	}
	l.messages = append(l.messages, userMsg)
	// Initialize and start the spinner for visual feedback
	// Use a dot spinner style
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// Start spinning in a goroutine
	s.Start()
	s.Suffix = " Sending codebase and querying LLM..."
	resp, err := api.SendMessage(apiURL, model, apiKey, l.messages)
	logger.Log(logger.INFO, "%s", "sending codebase and querying LLM")
	// Stop the spinner after the call completes
	s.Stop()
	if err != nil {
		logger.Log(logger.ERROR, "API call failed for user message: %v", err)
		// Remove last message from messages (user message)
		if len(l.messages) >= 1 {
			l.messages = l.messages[:len(l.messages)-1]
		}
		return "", err
	}
	logger.Log(logger.INFO, "received response from API for user message")
	// Add AI message to messages
	l.messages = append(l.messages, api.Message{
		Role:    "assistant",
		Content: resp,
	})
	// Update tokens amount
	l.updateTokens(model)
	return resp, nil
}

func (l *LLMx) GetTokens() int {
	return l.tokens
}

func (l *LLMx) ClearMessages() {
	l.messages = []api.Message{}
	l.updateTokens("")
}

// updateTokens calculates the current token count in the conversation context.
func (l *LLMx) updateTokens(model string) {
	encoding, err := tiktoken.EncodingForModel(model)
	if err != nil {
		// Fallback to cl100k_base (GPT-4 encoding)
		encoding, _ = tiktoken.GetEncoding("cl100k_base")
	}
	count := 0
	for _, v := range l.messages {
		count += len(encoding.Encode(v.Content, nil, nil))
	}
	l.tokens = count
}

// createSystemPrompt builds the system prompt by combining HZMIND.md and codebase data.
func createSystemPrompt() (string, error) {
	// Get code base
	files, err := codebase.GetCodeBase(".")
	if err != nil {
		return "", err
	}
	jsonCodeBase, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	// Load HZMIND.md
	data, err := os.ReadFile(common.PATH_FILE_README)
	if err != nil {
		output.Println()
		output.PrintfWarning("no %s file\n", common.FILE_IGNORE)
		logger.Log(logger.ERROR, "%v", err)
		data = []byte{}
	}
	// Create System Prompt message
	return string(data) + "\n\n## Codebase\n\n" + string(jsonCodeBase), nil
}
