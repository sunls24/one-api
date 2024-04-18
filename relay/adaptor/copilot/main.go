package copilot

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/relaymode"
)

func GetRequestURL(meta *meta.Meta) (string, error) {
	if meta.Mode == relaymode.ChatCompletions {
		return fmt.Sprintf("%s/chat/completions", meta.BaseURL), nil
	}
	return "", fmt.Errorf("unsupported relay mode %d for copilot", meta.Mode)
}

func SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	token, err := GetToken(meta.APIKey)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Host", "api.githubcopilot.com")
	req.Header.Set("X-Request-Id", uuid.New().String())
	req.Header.Set("Vscode-Sessionid", GetSessionID(meta.APIKey))
	req.Header.Set("Vscode-Machineid", GetMachineID(meta.APIKey))
	req.Header.Set("X-Github-Api-Version", "2023-07-07")
	req.Header.Set("Editor-Version", "vscode/1.89.1")
	req.Header.Set("Editor-Plugin-Version", "copilot-chat/0.15.1")
	req.Header.Set("User-Agent", "GitHubCopilotChat/0.15.1")
	req.Header.Set("Openai-Organization", "github-copilot")
	req.Header.Set("Copilot-Integration-Id", "vscode-chat")
	req.Header.Set("Openai-Intent", "conversation-panel")
	req.Header.Set("Connection", "close")
	return nil
}
