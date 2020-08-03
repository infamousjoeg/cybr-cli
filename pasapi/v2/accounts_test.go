// +build integration

package pasapi

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccounts(t *testing.T) {
	c := NewClient(os.Getenv("PAS_INTEGRATION_TOKEN"))

	ctx := context.Background()
	res, err := c.GetAccounts(nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting not-nil error")

	assert.Equal(t, 1, res.Count, "expecting 1 account found")

	assert.Equal(t, "97_3", res.Accounts[0].ID, "expecting correct account ID")
	assert.NotEmpty(t, res.Accounts[0].ID, "expecting non-empty account ID")
	assert.Greater(t, len(res.Accounts[0].SecretManagement), 0, "expecting non-empty secret management")
}
