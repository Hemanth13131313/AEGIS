package adapters

import "fmt"

var registry = map[string]func(backendURL, model string) Adapter{
	"openai":    func(url, model string) Adapter { return NewOpenAIAdapter(url, model) },
	"anthropic": func(url, model string) Adapter { return NewAnthropicAdapter(url, model) },
	"gemini":    func(url, model string) Adapter { return NewGeminiAdapter(url, model) },
	"bedrock":   func(url, model string) Adapter { return NewBedrockAdapter(url, model) },
}

func Get(provider, backendURL, model string) (Adapter, error) {
	factory, ok := registry[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q: supported providers are openai, anthropic, gemini, bedrock", provider)
	}
	return factory(backendURL, model), nil
}

func RegisterAdapter(name string, factory func(backendURL, model string) Adapter) {
	registry[name] = factory
}
