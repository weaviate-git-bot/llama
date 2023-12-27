package dispatcher

import (
	"context"
	"errors"

	"github.com/adrianliechti/llama/pkg/provider"
)

var (
	_ provider.Provider = &Provider{}
)

type Provider struct {
	models    map[string]provider.Model
	providers map[string]provider.Provider
}

func New(providers ...provider.Provider) (*Provider, error) {
	p := &Provider{
		models:    map[string]provider.Model{},
		providers: map[string]provider.Provider{},
	}

	for _, provider := range providers {
		models, err := provider.Models(context.Background())

		if err != nil {
			return nil, err
		}

		for _, m := range models {
			p.models[m.ID] = m
			p.providers[m.ID] = provider
		}
	}

	return p, nil
}

func (p *Provider) Models(ctx context.Context) ([]provider.Model, error) {
	result := make([]provider.Model, 0)

	for _, m := range p.models {
		result = append(result, m)
	}

	return result, nil
}

func (p *Provider) Embed(ctx context.Context, model, content string) (*provider.Embedding, error) {
	provider, ok := p.providers[model]

	if !ok {
		return nil, errors.New("no provider configured for model")
	}

	return provider.Embed(ctx, model, content)
}

func (p *Provider) Complete(ctx context.Context, model string, messages []provider.CompletionMessage) (*provider.Completion, error) {
	provider, ok := p.providers[model]

	if !ok {
		return nil, errors.New("no provider configured for model")
	}

	return provider.Complete(ctx, model, messages)
}

func (p *Provider) CompleteStream(ctx context.Context, model string, messages []provider.CompletionMessage, stream chan<- provider.Completion) error {
	provider, ok := p.providers[model]

	if !ok {
		return errors.New("no provider configured for model")
	}

	return provider.CompleteStream(ctx, model, messages, stream)
}