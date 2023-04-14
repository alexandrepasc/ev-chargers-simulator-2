package errors

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StructValidate(e ErroMsg, c *gin.Context, t translation.Translation) bool {
	err := validator.New().Struct(e)

	if err != nil {
		common.Log("StructValidate").Error(err)

		InternalServerError(c, t)

		return false
	}

	return true
}

func InternalServerError(c *gin.Context, t translation.Translation) {
	var r = ErroMsg{
		Message: t.Get(text.InternalServerError),
	}

	err := validator.New().Struct(r)

	if err != nil {
		common.Log("internalServerError").Fatal(err)
	}

	c.IndentedJSON(http.StatusInternalServerError, r)
}
