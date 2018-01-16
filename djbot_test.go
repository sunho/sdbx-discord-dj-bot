package djbot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const GoodToken = "NDAyNDkwNTM0OTkzNzIzMzky.DT_JQw.nZSQxUhUBDi6W-orIuqN2JvvJwI"

func TestConstants(t *testing.T) {
	assert.NotEmpty(t, GoodToken, "goodtoken is empty")
}
func TestCreation(t *testing.T) {
	botbase, err := NewFromToken("token", "!!", os.Stdout)
	assert.NotNil(t, err, "DJ made with wrong token should return error")
	assert.Nil(t, botbase, "wrongly initiallized DJ should be nil")
	botbase, err = NewFromToken(GoodToken, "!!", os.Stdout)
	assert.Nil(t, err, "DJ made with good token shold not return error")
	botbase.Close()
	assert.Nil(t, err, "closed DJ must be nil")
}

func TestHandler(t *testing.T) {
	botbase, err := NewFromToken(GoodToken, "!!", os.Stdout)
	assert.Nil(t, err, "DJ made with good token shold not return error")
	botbase.HandleNewMessage(nil, nil)

}
