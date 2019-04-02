package log

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextFormatter_Format(t *testing.T) {
	assert := assert.New(t)

	var entry = &logrus.Entry{
		Message: "hello world",
	}

	message, err := TextFormatter.Format(entry)
	assert.NoError(err)
	assert.Contains(string(message), `msg="hello world"`)

	TextFormatter.Prefix = "prefix"

	message, err = TextFormatter.Format(entry)
	assert.NoError(err)
	assert.Contains(string(message), `msg="[prefix] hello world"`)
}
