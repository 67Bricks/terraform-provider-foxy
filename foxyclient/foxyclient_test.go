package foxyclient

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func newFoxy() Foxy {
	c := readConfig()
	foxy, _ := New(c.BaseUrl, c.ClientID, c.ClientSecret, c.RefreshToken)
	return foxy
}

func TestReadConfig(t *testing.T) {
	c := readConfig()
	require.Equal(t, "client_1Q6iX3A1UjKNUZxEeV7P", c.ClientID)
	require.Len(t, c.ClientSecret, 40) // Not asserting what it is, but it must be set
}

func TestRetrieveToken(t *testing.T) {
	conf := readConfig()
	foxy := FoxyHttpClient{baseUrl: conf.BaseUrl}
	result, err := foxy.retrieveToken(conf.ClientID, conf.ClientSecret, conf.RefreshToken)
	require.Nil(t, err, "Should not have had error")
	require.NotEmpty(t, result.AccessToken)
}
