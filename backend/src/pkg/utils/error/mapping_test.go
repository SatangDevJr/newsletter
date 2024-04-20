package error_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	newsletterError "newsletter/src/pkg/utils/error"
)

func TestError_MappingErrorMessage(t *testing.T) {

	t.Run("should return Something when not found in mapping and input language is th", func(t *testing.T) {
		statusCode, err := newsletterError.MapMessageError("NO-HAVE-CODE-IN-MAP", "th")

		expectedError := newsletterError.NewError(newsletterError.InternalServerError, "เกิดข้อผิดพลาดบางอย่าง")
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, *expectedError, err)
	})

	t.Run("should return Something when not found in mapping and input language is en", func(t *testing.T) {
		statusCode, err := newsletterError.MapMessageError("NO-HAVE-CODE-IN-MAP", "en")

		expectedError := newsletterError.NewError(newsletterError.InternalServerError, "Something was wrong.")
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, *expectedError, err)
	})

	t.Run("should return error in mapping when found error code", func(t *testing.T) {
		statusCode, err := newsletterError.MapMessageError(newsletterError.Forbidden, "en")

		expectedError := newsletterError.NewError(newsletterError.Forbidden, newsletterError.ForbiddenMessage)
		assert.Equal(t, http.StatusForbidden, statusCode)
		assert.Equal(t, *expectedError, err)
	})

	t.Run("should return error message th when input language th", func(t *testing.T) {
		statusCode, err := newsletterError.MapMessageError(newsletterError.Forbidden, "th")

		expectedError := newsletterError.NewError(newsletterError.Forbidden, "ไม่มีสิทธิ์เข้าถึง")
		assert.Equal(t, http.StatusForbidden, statusCode)
		assert.Equal(t, *expectedError, err)
	})
}
