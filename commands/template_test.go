package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestTemplating(t *testing.T) {
	jsonText := `
	{
		"some-template": [
			{ "Repo": "bingo-local" },
			{ "Repo": "bango-local" },
			{ "Repo": "bongo-local" }
		]
	}`

	someTemplateText := `
		"files": [{
			"aql": {
				"items.find": {
					"repo": "{{.Repo}}",
				}
			}
		}]
	`

	var subscriptions map[string][]interface{}
	jsonErr := json.Unmarshal([]byte(jsonText), &subscriptions)
	require.NoError(t, jsonErr)

	for subscription, entries := range subscriptions {
		fmt.Println("Subscription:", subscription)

		someTemplate, parseErr := template.New(subscription).Parse(someTemplateText)
		require.NoError(t, parseErr)
		for _, entry := range entries {
			templatingErr := someTemplate.Execute(os.Stdout, entry)
			require.NoError(t, templatingErr)
		}
	}
}
