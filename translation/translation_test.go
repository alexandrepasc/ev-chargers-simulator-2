package translation_test

import (
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/stretchr/testify/assert"
)

func TestGetTranslation(t *testing.T) {
	ts := translation.Translation{L: translation.PtPt}

	a := ts.Get(text.MethodNotAllowed)

	e := "Método não permitido."

	assert.Equal(t, e, a)
}

func TestGetList(t *testing.T) {
	a := translation.Translation{}.List()

	e := []string{"en-GB", "pt-PT"}

	assert.Equal(t, e, a)
}

func TestGetKey(t *testing.T) {
	a := translation.Translation{}.GetKey("pt-PT")

	e := translation.PtPt

	assert.Equal(t, e, a)
}

func TestGetKeyWrongString(t *testing.T) {
	a := translation.Translation{}.GetKey("asd")

	e := translation.EnGb

	assert.Equal(t, e, a)
}
