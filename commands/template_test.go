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
	subscriptionsFile, readErr := os.ReadFile("../examples/subscriptions.json")
	require.NoError(t, readErr)

	var subscriptions map[string][]interface{}
	jsonErr := json.Unmarshal(subscriptionsFile, &subscriptions)
	require.NoError(t, jsonErr)

	for subscription, entries := range subscriptions {
		fmt.Println("Subscription:", subscription)
		templateText, readErr := os.ReadFile("../examples/templates/" + subscription + ".json")
		require.NoError(t, readErr)

		someTemplate, parseErr := template.New(subscription).Parse(string(templateText))
		require.NoError(t, parseErr)
		for _, entry := range entries {
			templatingErr := someTemplate.Execute(os.Stdout, entry)
			require.NoError(t, templatingErr)
		}
	}
}
