package zeitx_test

import (
	"testing"

	"github.com/gpollatos/zeitx-all/zeitx"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	const input = "testdata/simple.input"
	expected := &zeitx.Config{
		HTTPServer: zeitx.HTTPServerConfig{
			ListenAddr: ":8090",
		},
	}
	actual, err := zeitx.NewConfig(input)
	if err != nil {
		t.Errorf("the file should be properly loaded")
	}
	assert.Equal(t, expected, actual, "the ListenAddr whould be 8090")
}

func TestNewConfigWithExtraKeys(t *testing.T) {
	const input = "testdata/more_keys.input"
	expected := &zeitx.Config{
		HTTPServer: zeitx.HTTPServerConfig{
			ListenAddr: ":8090",
		},
	}
	actual, err := zeitx.NewConfig(input)
	if err != nil {
		t.Errorf("the file should be properly loaded")
	}
	assert.Equal(t, expected, actual, "the ListenAddr whould be 8090 and the extra keys should be ignored")
}

func TestNewConfigWithMissingKeys(t *testing.T) {
	const input = "testdata/missing_keys.input"
	_, err := zeitx.NewConfig(input)
	assert.NotNil(t, err, "there should be an error because of the missing listen address value")
}
