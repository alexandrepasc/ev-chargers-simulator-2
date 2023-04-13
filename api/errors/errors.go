package errors

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func InternalServerError(c *gin.Context, t translation.Translation) {
	var r = ErroMsg{
		Message: t.Get(text.InternalServerError),
	}

	err := validator.New().Struct(r)

	if err != nil {
		common.Log("internalServerError").Fatal(err)
	}

	c.JSON(http.StatusInternalServerError, r)
}
