package error_test

import (
	newsletter "newsletter/src/pkg/utils/error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewError(t *testing.T) {

	t.Run("should return error struct when new error", func(t *testing.T) {
		err := newsletter.NewError(newsletter.InternalServerError, newsletter.InternalServerErrorMessage)

		expectedError := &newsletter.Error{Code: newsletter.InternalServerError, Message: newsletter.InternalServerErrorMessage}
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error with data when error has data", func(t *testing.T) {
		err := newsletter.NewError(newsletter.InternalServerError, newsletter.InternalServerErrorMessage, "some data")

		expectedError := &newsletter.Error{
			Code:    newsletter.InternalServerError,
			Message: newsletter.InternalServerErrorMessage,
			Data:    "some data",
		}
		assert.Equal(t, expectedError, err)
	})

}
