package yarn

import (
	"os"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
)

func TestValidateSupportedCommand(t *testing.T) {
	yarnCmd := NewYarnCommand()

	testCases := []struct {
		args  []string
		valid bool
	}{
		{[]string{}, true},
		{[]string{"--json"}, true},
		{[]string{"npm", "publish", "--json"}, false},
		{[]string{"npm", "--json", "publish"}, false},
		{[]string{"npm", "tag", "list"}, false},
		{[]string{"npm", "info", "package-name"}, true},
		{[]string{"npm", "whoami"}, true},
	}

	for _, testCase := range testCases {
		yarnCmd.yarnArgs = testCase.args
		err := yarnCmd.validateSupportedCommand()
		assert.Equal(t, testCase.valid, err == nil, "Test args:", testCase.args)
	}
}

func TestSetAndRestoreEnvironmentVariables(t *testing.T) {
	const jfrogCliTestingEnvVar = "JFROG_CLI_ENV_VAR_FOR_TESTING"
	yarnCmd := &YarnCommand{envVarsBackup: make(map[string]*string)}

	// Check backup and restore of an existing variable
	setEnvCallback := tests.SetEnvWithCallbackAndAssert(t, jfrogCliTestingEnvVar, "abc")
	err := yarnCmd.backupAndSetEnvironmentVariable(jfrogCliTestingEnvVar, "new-value")
	assert.NoError(t, err)
	assert.Equal(t, "new-value", os.Getenv(jfrogCliTestingEnvVar))
	err = yarnCmd.restoreEnvironmentVariables()
	assert.NoError(t, err)
	assert.Equal(t, "abc", os.Getenv(jfrogCliTestingEnvVar))

	// Check backup and restore of a variable that doesn't exist
	setEnvCallback()
	err = yarnCmd.backupAndSetEnvironmentVariable(jfrogCliTestingEnvVar, "another-value")
	assert.NoError(t, err)
	assert.Equal(t, "another-value", os.Getenv(jfrogCliTestingEnvVar))
	err = yarnCmd.restoreEnvironmentVariables()
	assert.NoError(t, err)
	_, exist := os.LookupEnv(jfrogCliTestingEnvVar)
	assert.False(t, exist)
}

func TestExtractAuthIdentFromNpmAuth(t *testing.T) {
	testCases := []struct {
		responseFromArtifactory string
		expectedExtractedAuth   string
	}{
		{"_auth = Z290Y2hhISB5b3UgcmVhbGx5IHRoaW5rIGkgd291bGQgcHV0IHJlYWwgY3JlZGVudGlhbHMgaGVyZT8=\nalways-auth = true\nemail = notexist@mail.com\n", "Z290Y2hhISB5b3UgcmVhbGx5IHRoaW5rIGkgd291bGQgcHV0IHJlYWwgY3JlZGVudGlhbHMgaGVyZT8="},
		{"always-auth=true\nemail=notexist@mail.com\n_auth=TGVhcCBhbmQgdGhlIHJlc3Qgd2lsbCBmb2xsb3c=\n", "TGVhcCBhbmQgdGhlIHJlc3Qgd2lsbCBmb2xsb3c="},
	}

	for _, testCase := range testCases {
		actualExtractedAuth, err := extractAuthIdentFromNpmAuth(testCase.responseFromArtifactory)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedExtractedAuth, actualExtractedAuth)
	}
}
