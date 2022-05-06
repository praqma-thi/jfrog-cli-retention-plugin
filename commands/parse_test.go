package commands

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jfrog/jfrog-cli-core/v2/common/spec"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/require"
)

func TestParseSpec(t *testing.T) {
	var specFiles spec.SpecFiles
	somefile := "../example/filespecs/example.json"

	reader, fileErr := os.Open(somefile)
	require.NoError(t, fileErr)

	decodeErr := json.NewDecoder(reader).Decode(&specFiles)
	require.NoError(t, decodeErr)

	for _, file := range specFiles.Files {
		var (
			deleteParams services.DeleteParams
			castErr      error
		)
		deleteParams.CommonParams, castErr = file.ToCommonParams()
		require.NoError(t, castErr)

		spew.Dump(deleteParams)
	}
}
