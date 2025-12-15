package swagger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Parser struct {
	spec *Spec
}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Load(urlOrPath string) (*Spec, error) {
	var data []byte
	var err error

	if strings.HasPrefix(urlOrPath, "http://") || strings.HasPrefix(urlOrPath, "https://") {
		// Load from URL
		resp, err := http.Get(urlOrPath)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch swagger from URL: %w", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				err = closeErr
			}
		}()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch swagger: status %d", resp.StatusCode)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read swagger response: %w", err)
		}
	} else {
		// Load from a file
		data, err = os.ReadFile(urlOrPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read swagger file: %w", err)
		}
	}

	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse swagger JSON: %w", err)
	}

	// Also keep raw data for AI processing
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err == nil {
		spec.Raw = raw
	}

	return &spec, nil
}
