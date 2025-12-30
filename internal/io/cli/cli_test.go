package cli_test

import (
	_ "md/internal/infr"
	"md/internal/io/cli"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	app := cli.Initialize("/tmp/data/")
	assert.Nil(t, app.Run([]string{"", "user:create", "test_id", "test_password"}))
}
