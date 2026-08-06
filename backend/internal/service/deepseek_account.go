package service

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
)

const (
	DeepSeekProvider     = "deepseek"
	DeepSeekV4FlashModel = "deepseek-v4-flash"
)

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
