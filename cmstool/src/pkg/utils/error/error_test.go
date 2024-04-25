package error_test

import (
	subscribetool "subscribetool/src/pkg/utils/error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {

	t.Run("should return error struct when new error", func(t *testing.T) {
		err := subscribetool.NewError(subscribetool.InternalServerError, subscribetool.InternalServerErrorMessage)

		expectedError := &subscribetool.Error{Code: subscribetool.InternalServerError, Message: subscribetool.InternalServerErrorMessage}
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error with data when error has data", func(t *testing.T) {
		err := subscribetool.NewError(subscribetool.InternalServerError, subscribetool.InternalServerErrorMessage, "some data")

		expectedError := &subscribetool.Error{
			Code:    subscribetool.InternalServerError,
			Message: subscribetool.InternalServerErrorMessage,
			Data:    "some data",
		}
		assert.Equal(t, expectedError, err)
	})

}
