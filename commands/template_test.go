package commands

import (
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestTemplating(t *testing.T) {

	subscriptions := make([]map[string]string, 0, 3)
	subscriptions = append(subscriptions, map[string]string{"Repo": "bingo-local"})
	subscriptions = append(subscriptions, map[string]string{"Repo": "bango-local"})
	subscriptions = append(subscriptions, map[string]string{"Repo": "bongo-local"})

	templateText := `
		"files": [{
			"aql": {
				"items.find": {
					"repo": "{{.Repo}}",
				}
			}
		}]
	`

	templ, err := template.New("the-retention").Parse(templateText)
	require.NoError(t, err)

	err = templ.Execute(os.Stdout, subscriptions[0])
	require.NoError(t, err)
}
