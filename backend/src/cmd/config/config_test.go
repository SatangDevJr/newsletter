package config_test

import (
	"newsletter/src/cmd/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {

	t.Run("should return struct configuration when call new config", func(t *testing.T) {
		t.Setenv("PORT", "PORT")
		t.Setenv("APP_VERSION", "APP_VERSION")
		t.Setenv("DB_NAME", "DB_NAME")
		t.Setenv("DB_USERNAME", "DB_USERNAME")
		t.Setenv("DB_PASSWORD", "DB_PASSWORD")
		t.Setenv("DB_HOST", "DB_HOST")
		t.Setenv("DB_PORT", "DB_PORT")
		t.Setenv("ELS_URL", "ELS_URL")
		t.Setenv("ELS_USERNAME", "ELS_USERNAME")
		t.Setenv("ELS_PASSWORD", "ELS_PASSWORD")
		t.Setenv("ELS_INDEX", "ELS_INDEX")
		t.Setenv("HOST_WEB", "HOST_WEB")

		resNew := config.New()

		assert.Equal(t, "PORT", resNew.PORT)
		assert.Equal(t, "APP_VERSION", resNew.AppVersion)
		assert.Equal(t, "DB_NAME", resNew.DBName)
		assert.Equal(t, "DB_USERNAME", resNew.DBUserName)
		assert.Equal(t, "DB_PASSWORD", resNew.DBPassword)
		assert.Equal(t, "DB_HOST", resNew.DBHost)
		assert.Equal(t, "DB_PORT", resNew.DBPort)
		assert.Equal(t, "ELS_URL", resNew.ELSURL)
		assert.Equal(t, "ELS_USERNAME", resNew.ELSUsername)
		assert.Equal(t, "ELS_PASSWORD", resNew.ELSPassword)
		assert.Equal(t, "ELS_INDEX", resNew.ELSIndex)
		assert.Equal(t, "", resNew.Stage)
		assert.Equal(t, "HOST_WEB", resNew.Origin)
	})
}
