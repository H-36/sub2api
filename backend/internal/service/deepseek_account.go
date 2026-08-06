package service

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
)

const (
	DeepSeekProvider     = "deepseek"
	PlatformDeepSeek     = DeepSeekProvider
	DeepSeekV4FlashModel = "deepseek-v4-flash"
)

// GroupExecutionPlatform keeps DeepSeek as a distinct group brand while
// reusing the existing OpenAI-compatible execution path.
func GroupExecutionPlatform(platform string) string {
	if platform == PlatformDeepSeek {
		return PlatformOpenAI
	}
	return platform
}

func (a *Account) IsDeepSeekAPIKeyAccount() bool {
	return a != nil &&
		a.IsOpenAIApiKey() &&
		strings.EqualFold(strings.TrimSpace(a.GetExtraString("provider")), DeepSeekProvider)
}

func (a *Account) DefaultOpenAITestModel() string {
	if a.IsDeepSeekAPIKeyAccount() {
		return DeepSeekV4FlashModel
	}
	return openai.DefaultTestModel
}
