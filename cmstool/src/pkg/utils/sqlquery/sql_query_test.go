package sqlquery_test

import (
	"subscribetool/src/pkg/utils/convert"
	"subscribetool/src/pkg/utils/sqlquery"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSQLQuery_GenerateQueryColumnNameInclude(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{})

		assert.Equal(t, "", result)
	})

	t.Run("should return empty when input struct and input parameter includefields is empty", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{})

		assert.Equal(t, "", result)
	})

	t.Run("should return single column when input struct single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{"Name"})

		assert.Equal(t, "[name]", result)
	})

	t.Run("should return multiple column when input struct multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{"Name", "LastName"})

		assert.Equal(t, "[name],[lastname]", result)
	})

	t.Run("should return specific column without sql ignore field when input struct sql ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"-"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{"Name", "LastName"})

		assert.Equal(t, "[lastname]", result)
	})

	t.Run("should return specific column without tag sql field when input struct tag without sql field", func(t *testing.T) {
		input := struct {
			Name     string
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{"Name", "LastName"})

		assert.Equal(t, "[lastname]", result)
	})

	t.Run("should return column include ref id when input include ref id", func(t *testing.T) {
		input := struct {
			ID    string `sql:"id"`
			RefID string `sql:"refId"`
		}{}

		result := sqlquery.GenerateQueryColumnNameInclude(input, []string{"RefID"})

		assert.Equal(t, "[refId]", result)
	})
}

func TestSQLQuery_GenerateQueryColumnNameIncludeAlias(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{}, "pr")

		assert.Equal(t, "", result)
	})

	t.Run("should return empty when input struct and input parameter includefields is empty", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{}, "pr")

		assert.Equal(t, "", result)
	})

	t.Run("should return single column when input struct single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{"Name"}, "pr")

		assert.Equal(t, "pr.[name]", result)
	})

	t.Run("should return multiple column when input struct multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{"Name", "LastName"}, "pr")

		assert.Equal(t, "pr.[name],pr.[lastname]", result)
	})

	t.Run("should return specific column without sql ignore field when input struct sql ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"-"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{"Name", "LastName"}, "pr")

		assert.Equal(t, "pr.[lastname]", result)
	})

	t.Run("should return specific column without tag sql field when input struct tag without sql field", func(t *testing.T) {
		input := struct {
			Name     string
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{"Name", "LastName"}, "pr")

		assert.Equal(t, "pr.[lastname]", result)
	})

	t.Run("should return column include ref id when input include ref id", func(t *testing.T) {
		input := struct {
			ID    string `sql:"id"`
			RefID string `sql:"refId"`
		}{}

		result := sqlquery.GenerateQueryColumnNameIncludeAlias(input, []string{"RefID"}, "pr")

		assert.Equal(t, "pr.[refId]", result)
	})
}

func TestSQLQuery_GenerateQueryColumnValuesInclude(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{})

		assert.Equal(t, "", result)
	})

	t.Run("should return single column and value when input struct string single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Name",
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"Name"})

		assert.Equal(t, "N'Name'", result)
	})

	t.Run("should return single column and value when input struct string single field include '", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Test 'Name'",
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"Name"})

		assert.Equal(t, "N'Test ''Name'''", result)
	})

	t.Run("should return value from pointer string fields when input struct pointer string fields", func(t *testing.T) {
		input := struct {
			Name *string `sql:"name"`
		}{
			Name: convert.ValueToStringPointer("Test 'Name'"),
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{
			"Name",
			"Total",
			"IsShow",
			"DocumentDate",
		})

		assert.Equal(t, "N'Test ''Name'''", result)
	})

	t.Run("should return multiple column and value when input struct string multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"Name", "LastName"})

		assert.Equal(t, "N'Name',N'LastName'", result)
	})

	t.Run("should return query field type of bool when input bool", func(t *testing.T) {
		input := struct {
			IsShow bool `sql:"isShow"`
		}{
			IsShow: true,
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"IsShow"})

		assert.Equal(t, "N'true'", result)
	})

	t.Run("should return query field type of time when input time", func(t *testing.T) {
		dateTime, _ := time.Parse(time.RFC822, "01 Jan 21 10:00 UTC")
		input := struct {
			DocumentDate time.Time `sql:"documentDate"`
		}{
			DocumentDate: dateTime,
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"DocumentDate"})

		assert.Equal(t, "N'"+input.DocumentDate.Format(time.RFC3339)+"'", result)
	})

	t.Run("should return query field specific type when input without string, time and bool", func(t *testing.T) {
		input := struct {
			Total int32 `sql:"total"`
		}{
			Total: 10,
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"Total"})

		assert.Equal(t, "10", result)
	})

	t.Run("should return query field to be null when input nil", func(t *testing.T) {
		input := struct {
			FileID *int32 `sql:"fileId"`
		}{}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"FileID"})

		assert.Equal(t, "null", result)
	})

	t.Run("should return query without sql ignore field when input struct sql ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"-"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"Name", "LastName"})

		assert.Equal(t, "N'LastName'", result)
	})

	t.Run("should return value from pointer fields when input struct pointer fields", func(t *testing.T) {
		dateTime, _ := time.Parse(time.RFC822, "01 Jan 21 10:00 UTC")
		input := struct {
			Name         *string    `sql:"name"`
			Total        *int       `sql:"total"`
			IsShow       *bool      `sql:"isShow"`
			DocumentDate *time.Time `sql:"documentDate"`
		}{
			Name:         convert.ValueToStringPointer("Name"),
			Total:        convert.ValueToIntPointer(10),
			IsShow:       convert.ValueToBoolPointer(true),
			DocumentDate: convert.ValueToTimePointer(dateTime),
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{
			"Name",
			"Total",
			"IsShow",
			"DocumentDate",
		})

		assert.Equal(t, "N'Name',10,N'true'"+",N'"+input.DocumentDate.Format(time.RFC3339)+"'", result)
	})

	t.Run("should return value ref id when input include ref id", func(t *testing.T) {
		input := struct {
			ID    *int32 `sql:"id"`
			RefID *int32 `sql:"refId"`
		}{
			ID:    convert.ValueToInt32Pointer(1),
			RefID: convert.ValueToInt32Pointer(2),
		}

		result := sqlquery.GenerateQueryColumnValuesInclude(input, []string{"RefID"})

		assert.Equal(t, "2", result)
	})
}

func TestSQLQuery_GenerateQueryColumnNamesWithAlias(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "")

		assert.Equal(t, "", result)
	})

	t.Run("should return single column when input struct single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "")

		assert.Equal(t, "[name]", result)
	})

	t.Run("should return multiple column when input struct multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "")

		assert.Equal(t, "[name],[lastname]", result)
	})

	t.Run("should return specific column without sql ignore field when input struct sql ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"-"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "")

		assert.Equal(t, "[lastname]", result)
	})

	t.Run("should return specific column without tag sql field when input struct tag without sql field", func(t *testing.T) {
		input := struct {
			Name     string
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "")

		assert.Equal(t, "[lastname]", result)
	})

	t.Run("should return query without ignore field when input ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{"Name"}, "")

		assert.Equal(t, "[lastname]", result)
	})

	t.Run("should return query with alias table name when input alias name", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			Lastname string `sql:"lastname"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "t")

		assert.Equal(t, "t.[name],t.[lastname]", result)
	})

	t.Run("should return query with alias table ref id when input alias name", func(t *testing.T) {
		input := struct {
			ID        *int32 `sql:"id"`
			CompanyID *int32 `sql:"companyId"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{}, "t")

		assert.Equal(t, "t.[id],t.[companyId]", result)
	})

	t.Run("should return query with alias table ref id ignore when input alias name", func(t *testing.T) {
		input := struct {
			ID        *int32 `sql:"id"`
			CompanyID *int32 `sql:"companyId"`
		}{}

		result := sqlquery.GenerateQueryColumnNamesWithAlias(input, []string{"CompanyID"}, "t")

		assert.Equal(t, "t.[id]", result)
	})
}

func TestSQLQuery_GenerateQueryColumnValues(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "", result)
	})

	t.Run("should return single column and value when input struct string single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Name",
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'Name'", result)
	})

	t.Run("should return multiple column and value when input struct string multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'Name',N'LastName'", result)
	})

	t.Run("should return query field type of bool when input bool", func(t *testing.T) {
		input := struct {
			IsShow bool `sql:"isShow"`
		}{
			IsShow: true,
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'true'", result)
	})

	t.Run("should return query field type of time when input time", func(t *testing.T) {
		dateTime, _ := time.Parse(time.RFC822, "01 Jan 21 10:00 UTC")
		input := struct {
			DocumentDate time.Time `sql:"documentDate"`
		}{
			DocumentDate: dateTime,
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'"+input.DocumentDate.Format(time.RFC3339)+"'", result)
	})

	t.Run("should return query field specific type when input without string, time and bool", func(t *testing.T) {
		input := struct {
			Total int32 `sql:"total"`
		}{
			Total: 10,
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "10", result)
	})

	t.Run("should return query field to be null when input nil", func(t *testing.T) {
		input := struct {
			FileID *int32 `sql:"fileId"`
		}{}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "null", result)
	})

	t.Run("should return query without ignore field when input ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{"Name"})

		assert.Equal(t, "N'LastName'", result)
	})

	t.Run("should return query without sql ignore field when input struct sql ignore field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"-"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'LastName'", result)
	})

	t.Run("should return specific column without tag sql field when input struct tag without sql field", func(t *testing.T) {
		input := struct {
			Name     string
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'LastName'", result)
	})

	t.Run("should return value from pointer fields when input struct pointer fields", func(t *testing.T) {
		dateTime, _ := time.Parse(time.RFC822, "01 Jan 21 10:00 UTC")
		input := struct {
			Name         *string    `sql:"name"`
			Total        *int       `sql:"total"`
			IsShow       *bool      `sql:"isShow"`
			DocumentDate *time.Time `sql:"documentDate"`
		}{
			Name:         convert.ValueToStringPointer("Name"),
			Total:        convert.ValueToIntPointer(10),
			IsShow:       convert.ValueToBoolPointer(true),
			DocumentDate: convert.ValueToTimePointer(dateTime),
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'Name',10,N'true'"+",N'"+input.DocumentDate.Format(time.RFC3339)+"'", result)
	})

	t.Run("should return value not ' when input value include '", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Test 'Name'",
		}
		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'Test ''Name'''", result)
	})

	t.Run("should return value from pointer fields when input struct pointer fields", func(t *testing.T) {
		input := struct {
			Name *string `sql:"name"`
		}{
			Name: convert.ValueToStringPointer("Test 'Name'"),
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{})

		assert.Equal(t, "N'Test ''Name'''", result)
	})

	t.Run("should return value without ref id when input ignore value ref id", func(t *testing.T) {
		input := struct {
			ID    *int32 `sql:"id"`
			RefID *int32 `sql:"refId"`
		}{
			ID:    convert.ValueToInt32Pointer(1),
			RefID: convert.ValueToInt32Pointer(2),
		}

		result := sqlquery.GenerateQueryColumnValues(input, []string{"RefID"})

		assert.Equal(t, "1", result)
	})
}
func TestSQLQuery_GenerateQueryUpdateFields(t *testing.T) {
	t.Run("should return empty query when input empty struct", func(t *testing.T) {
		input := struct{}{}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "", result)
	})

	t.Run("should return single column and value when input struct string single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Name",
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[name] = N'Name'", result)
	})

	t.Run("should return single column and value when input struct string single field", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Test 'Name'",
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[name] = N'Test ''Name'''", result)
	})

	t.Run("should return multiple column and value when input struct string multiple field", func(t *testing.T) {
		input := struct {
			Name     string `sql:"name"`
			LastName string `sql:"lastname"`
		}{
			Name:     "Name",
			LastName: "LastName",
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[name] = N'Name',[lastname] = N'LastName'", result)
	})

	t.Run("should return query field type of bool when input bool", func(t *testing.T) {
		input := struct {
			IsShow bool `sql:"isShow"`
		}{
			IsShow: true,
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[isShow] = N'true'", result)
	})

	t.Run("should return query field type of time when input time", func(t *testing.T) {
		dateTime, _ := time.Parse(time.RFC822, "01 Jan 21 10:00 UTC")
		input := struct {
			DocumentDate time.Time `sql:"documentDate"`
		}{
			DocumentDate: dateTime,
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[documentDate] = N'"+input.DocumentDate.Format(time.RFC3339)+"'", result)
	})

	t.Run("should return query field specific type when input without string, time and bool", func(t *testing.T) {
		input := struct {
			Total int32 `sql:"total"`
		}{
			Total: 10,
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[total] = 10", result)
	})

	t.Run("should return query field to be null when input nil", func(t *testing.T) {
		input := struct {
			FileID *int32 `sql:"fileId"`
		}{}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[fileId] = null", result)
	})

	t.Run("should return query without id when input id", func(t *testing.T) {
		input := struct {
			ID int64 `sql:"id"`
		}{
			ID: 1,
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "", result)
	})

	t.Run("should return query that value include comma", func(t *testing.T) {
		input := struct {
			Name string `sql:"name"`
		}{
			Name: "Test, Comma",
		}

		result := sqlquery.GenerateQueryUpdateFields(input, []string{"ID", "CreatedDate"})

		assert.Equal(t, "[name] = N'Test, Comma'", result)
	})
}

func TestSQLQuery_GenerateQueryUpdateFieldInclude(t *testing.T) {

	t.Run("should return empty when input empty struct", func(t *testing.T) {
		input := struct{}{}

		actual := sqlquery.GenerateQueryUpdateFieldInclude(input, []string{})

		expected := ""
		assert.Equal(t, expected, actual)
	})

	t.Run("should return two field when have two annotations", func(t *testing.T) {
		input := struct {
			ID   int64  `sql:"id"`
			Name string `sql:"name"`
		}{
			ID:   1,
			Name: "Name",
		}

		actual := sqlquery.GenerateQueryUpdateFieldInclude(input, []string{"ID", "Name"})

		expected := "[id] = 1,[name] = N'Name'"
		assert.Equal(t, expected, actual)
	})

	t.Run("should return name when include", func(t *testing.T) {
		input := struct {
			ID   int64  `sql:"id"`
			Name string `sql:"name"`
		}{
			ID:   1,
			Name: "Name",
		}

		actual := sqlquery.GenerateQueryUpdateFieldInclude(input, []string{"Name"})

		expected := "[name] = N'Name'"
		assert.Equal(t, expected, actual)
	})

	t.Run("should return empty when include empty", func(t *testing.T) {
		input := struct {
			ID   int64  `sql:"id"`
			Name string `sql:"name"`
		}{
			ID:   1,
			Name: "Name",
		}

		actual := sqlquery.GenerateQueryUpdateFieldInclude(input, []string{})

		expected := ""
		assert.Equal(t, expected, actual)
	})
}

func TestSQLQuery_GenerateTransactionAndRollback(t *testing.T) {
	t.Run("should return transaction query", func(t *testing.T) {

		result := sqlquery.GenerateTransactionAndRollback([]string{"TRUNCATE TABLE TABLE_A"})

		expectedQuery := `
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		TRUNCATE TABLE TABLE_A

	COMMIT TRANSACTION;
	END TRY

	BEGIN CATCH
	ROLLBACK TRANSACTION;
	DECLARE
	    @ErrorMessage nvarchar(4000) = ERROR_MESSAGE(),
	    @ErrorNumber int = ERROR_NUMBER(),
	    @ErrorSeverity int = ERROR_SEVERITY(),
	    @ErrorState int = ERROR_STATE(),
	    @ErrorLine int = ERROR_LINE(),
	    @ErrorProcedure nvarchar(200) = ISNULL(ERROR_PROCEDURE(), '-');
	SELECT @ErrorMessage = N'Error %d, Level %d, State %d, Procedure %s, Line %d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`
		assert.Equal(t, expectedQuery, result)
	})

	t.Run("should return transaction query when input multiple queries", func(t *testing.T) {

		result := sqlquery.GenerateTransactionAndRollback([]string{"SELECT * FROM TABLE_A", "TRUNCATE TABLE TABLE_A"})

		expectedQuery := `
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		SELECT * FROM TABLE_A
	
	TRUNCATE TABLE TABLE_A

	COMMIT TRANSACTION;
	END TRY

	BEGIN CATCH
	ROLLBACK TRANSACTION;
	DECLARE
	    @ErrorMessage nvarchar(4000) = ERROR_MESSAGE(),
	    @ErrorNumber int = ERROR_NUMBER(),
	    @ErrorSeverity int = ERROR_SEVERITY(),
	    @ErrorState int = ERROR_STATE(),
	    @ErrorLine int = ERROR_LINE(),
	    @ErrorProcedure nvarchar(200) = ISNULL(ERROR_PROCEDURE(), '-');
	SELECT @ErrorMessage = N'Error %d, Level %d, State %d, Procedure %s, Line %d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`
		assert.Equal(t, expectedQuery, result)
	})
}

func TestSQLQuery_GenerateTransactionAndRollbackWithRepeatableRead(t *testing.T) {
	t.Run("should return transaction query", func(t *testing.T) {

		result := sqlquery.GenerateTransactionAndRollbackWithRepeatableRead([]string{"TRUNCATE TABLE TABLE_A"})

		expectedQuery := `
	SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		TRUNCATE TABLE TABLE_A

	COMMIT TRANSACTION;
	END TRY

	BEGIN CATCH
	ROLLBACK TRANSACTION;
	DECLARE
	    @ErrorMessage nvarchar(4000) = ERROR_MESSAGE(),
	    @ErrorNumber int = ERROR_NUMBER(),
	    @ErrorSeverity int = ERROR_SEVERITY(),
	    @ErrorState int = ERROR_STATE(),
	    @ErrorLine int = ERROR_LINE(),
	    @ErrorProcedure nvarchar(200) = ISNULL(ERROR_PROCEDURE(), '-');
	SELECT @ErrorMessage = N'Error %d, Level %d, State %d, Procedure %s, Line %d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`
		assert.Equal(t, expectedQuery, result)
	})

	t.Run("should return transaction query when input multiple queries", func(t *testing.T) {

		result := sqlquery.GenerateTransactionAndRollbackWithRepeatableRead([]string{"SELECT * FROM TABLE_A", "TRUNCATE TABLE TABLE_A"})

		expectedQuery := `
	SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		SELECT * FROM TABLE_A
	
	TRUNCATE TABLE TABLE_A

	COMMIT TRANSACTION;
	END TRY

	BEGIN CATCH
	ROLLBACK TRANSACTION;
	DECLARE
	    @ErrorMessage nvarchar(4000) = ERROR_MESSAGE(),
	    @ErrorNumber int = ERROR_NUMBER(),
	    @ErrorSeverity int = ERROR_SEVERITY(),
	    @ErrorState int = ERROR_STATE(),
	    @ErrorLine int = ERROR_LINE(),
	    @ErrorProcedure nvarchar(200) = ISNULL(ERROR_PROCEDURE(), '-');
	SELECT @ErrorMessage = N'Error %d, Level %d, State %d, Procedure %s, Line %d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`
		assert.Equal(t, expectedQuery, result)
	})
}
