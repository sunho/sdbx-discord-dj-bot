package djbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const GoodToken = ""
const GoodId = ""
const GoodPw = ""

func TestConstants(t *testing.T) {
	assert.NotEmpty(t, GoodToken, "goodtoken is empty")
	assert.NotEmpty(t, GoodPw, "goodid is empty")
	assert.NotEmpty(t, GoodId, "goodpw is empty")
}
func TestBotBase(t *testing.T) {
	botbase, err := BotBase.NewBot("token")
	assert.Nil(t, err, "botbase made with wrong token should return error")
	assert.NotNil(t, botbase, "wrongly initiallized botbase should be nil")
	botbase, err = BotBase.NewBot(GoodToken)
	assert.NotNil(t, err, "botbase made with good token shold not return error")
	botbase.Close()
}

func TestBotBaseAddCommand(t *testing.T) {
	botbase, err := BotBase.NewBot(GoodToken)
	cmd := Command.New("test")
	botbase.AddCommand()
}
