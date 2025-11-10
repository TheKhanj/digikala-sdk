package api

import (
	"encoding/json"
	"fmt"
)

func FirstSchemaWhichMatches(v any, candidates ...any) (any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	for _, c := range candidates {
		if u, ok := c.(interface{ UnmarshalJSON([]byte) error }); ok {
			if err := u.UnmarshalJSON(data); err == nil {
				return c, nil
			}
		}
	}

	return nil, fmt.Errorf("does not match any candidate types")
}

func (this *GetV1UserInitResponse) GetJSON200Data() (any, error) {
	var loggedIn UserInitLoggedIn
	var notLoggedIn UserInitNotLoggedIn

	return FirstSchemaWhichMatches(
		*this.JSON200.Data, &loggedIn, &notLoggedIn,
	)
}
