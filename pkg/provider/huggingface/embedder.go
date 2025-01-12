package huggingface

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/adrianliechti/llama/pkg/provider"
)

var _ provider.Embedder = (*Embedder)(nil)

type Embedder struct {
	*Config
}

func NewEmbedder(url string, options ...Option) (*Embedder, error) {
	cfg := &Config{
		url: url,

		token: "-",
		model: "tei",

		client: http.DefaultClient,
	}

	for _, option := range options {
		option(cfg)
	}

	return &Embedder{
		Config: cfg,
	}, nil
}

func (e *Embedder) Embed(ctx context.Context, content string) ([]float32, error) {
	body := map[string]any{
		"inputs": strings.TrimSpace(content),
	}

	resp, err := e.client.Post(e.url, "application/json", jsonReader(body))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to encode input")
	}

	defer resp.Body.Close()

	var result []float32

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
