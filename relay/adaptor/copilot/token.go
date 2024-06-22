package copilot

import (
	"fmt"
	"github.com/songquanpeng/one-api/common/client"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/songquanpeng/one-api/common/cache"
)

type auth struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	m         sync.Mutex
}

func createAuth(key string) *auth {
	return &auth{}
}

var cacheAuth = cache.NewCache[*auth](createAuth)

func GetToken(gt string) (string, error) {
	cached, ok := cacheAuth.Get(gt)
	cached.m.Lock()
	defer cached.m.Unlock()
	if ok && cached.ExpiresAt > time.Now().Unix()+int64(rand.Intn(600)+300) {
		return cached.Token, nil
	}

	const tokenURL = "https://api.github.com/copilot_internal/v2/token"
	req, _ := http.NewRequest(http.MethodGet, tokenURL, nil)
	req.Header.Set("Host", "api.github.com")
	req.Header.Set("Authorization", "token "+gt)
	req.Header.Set("Editor-Version", "vscode/1.89.1")
	req.Header.Set("Editor-Plugin-Version", "copilot-chat/0.15.1")
	req.Header.Set("User-Agent", "GitHubCopilotChat/0.15.1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "close")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "request copilot token failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request copilot token code != 200: %s", string(msg))
	}
	var res auth
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "request copilot token read response failed")
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return "", errors.Wrap(err, "request copilot token unmarshal failed")
	}
	cacheAuth.Set(gt, &res)
	return res.Token, nil
}
