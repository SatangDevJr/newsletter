package sqlquery

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/thoas/go-funk"
)

const (
	TimeType   = "time.Time"
	StringType = "string"
	BoolType   = "bool"

	PointerTimeType   = "*time.Time"
	PointerStringType = "*string"
	PointerBoolType   = "*bool"
)

func GenerateQueryColumnNameInclude(s interface{}, includeFields []string) string {
	t := reflect.TypeOf(s)
	strFields := []string{}
	if t.NumField() == 0 {
		return ""
	}

	if len(includeFields) == 0 {
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Tag.Get("sql")
		if funk.ContainsString(includeFields, t.Field(i).Name) {
			if field != "-" && field != "" {
				strFields = append(strFields, fmt.Sprintf("[%s]", field))
			}
		}
	}

	return strings.Join(strFields, ",")
}

func GenerateQueryColumnNameIncludeAlias(s interface{}, includeFields []string, aliasName string) string {
	t := reflect.TypeOf(s)
	strFields := []string{}
	if t.NumField() == 0 {
		return ""
	}

	if len(includeFields) == 0 {
		return ""
	}

	if aliasName != "" {
		aliasName = fmt.Sprintf("%s.", aliasName)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Tag.Get("sql")
		if funk.ContainsString(includeFields, t.Field(i).Name) {
			if field != "-" && field != "" {
				strFields = append(strFields, fmt.Sprintf("%s[%s]", aliasName, field))
			}
		}
	}

	return strings.Join(strFields, ",")
}

func GenerateQueryColumnValuesInclude(s interface{}, includeFields []string) string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	strValues := []string{}

	if t.NumField() == 0 {
		return ""
	}

	for i := 0; i < v.NumField(); i++ {
		if !funk.ContainsString(includeFields, t.Field(i).Name) {
			continue
		}

		field := t.Field(i).Tag.Get("sql")
		if field == "-" || field == "" {
			continue
		}

		value := v.Field(i).Interface()

		if fmt.Sprintf("%v", value) == "<nil>" {
			value = "null"
		}

		if fmt.Sprintf("%s", reflect.TypeOf(value)) == TimeType {
			value = value.(time.Time).Format(time.RFC3339)
		} else if fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerTimeType {
			value = value.(*time.Time).Format(time.RFC3339)
		}

		if value != "null" && (fmt.Sprintf("%s", reflect.TypeOf(value)) == TimeType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == StringType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == BoolType) {
			if fmt.Sprintf("%s", reflect.TypeOf(value)) == StringType {
				strValues = append(strValues, fmt.Sprintf("N'%v'", strings.ReplaceAll(value.(string), "'", "''")))
			} else {
				strValues = append(strValues, fmt.Sprintf("N'%v'", value))
			}
		} else if value != "null" && (fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerTimeType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerStringType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerBoolType) {
			if fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerStringType {
				strValue := reflect.ValueOf(value).Elem().String()
				strValues = append(strValues, fmt.Sprintf("N'%v'", strings.ReplaceAll(strValue, "'", "''")))
			} else {
				strValues = append(strValues, fmt.Sprintf("N'%v'", reflect.ValueOf(value).Elem()))
			}
		} else {
			if strings.Contains(fmt.Sprintf("%s", reflect.ValueOf(value)), "*") {
				strValues = append(strValues, fmt.Sprintf("%v", reflect.ValueOf(value).Elem()))
			} else {
				strValues = append(strValues, fmt.Sprintf("%v", value))
			}
		}

	}

	return strings.Join(strValues, ",")
}

func GenerateQueryColumnNames(s interface{}, ignoreFields []string) string {
	return GenerateQueryColumnNamesWithAlias(s, ignoreFields, "")
}

func GenerateQueryColumnNamesWithAlias(s interface{}, ignoreFields []string, aliasName string) string {
	t := reflect.TypeOf(s)
	strs := []string{}

	if aliasName != "" {
		aliasName = fmt.Sprintf("%s.", aliasName)
	}

	for i := 0; i < t.NumField(); i++ {
		if funk.ContainsString(ignoreFields, t.Field(i).Name) {
			continue
		}

		field := t.Field(i).Tag.Get("sql")
		if field != "-" && field != "" {
			strs = append(strs, fmt.Sprintf("%s[%s]", aliasName, field))
		}

	}
	return strings.Join(strs, ",")
}

func ManipulateValues(s interface{}, ignoreFields []string) []string {
	v := reflect.ValueOf(s)
	typeOfS := v.Type()
	strs := []string{}

	for i := 0; i < v.NumField(); i++ {
		if funk.ContainsString(ignoreFields, typeOfS.Field(i).Name) {
			continue
		}

		field := typeOfS.Field(i).Tag.Get("sql")
		if field == "-" || field == "" {
			continue
		}

		value := v.Field(i).Interface()

		if fmt.Sprintf("%v", value) == "<nil>" {
			value = "null"
		}

		if fmt.Sprintf("%s", reflect.TypeOf(value)) == TimeType {
			value = value.(time.Time).Format(time.RFC3339)
		} else if fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerTimeType {
			value = value.(*time.Time).Format(time.RFC3339)
		}

		if value != "null" && (fmt.Sprintf("%s", reflect.TypeOf(value)) == TimeType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == StringType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == BoolType) {
			if fmt.Sprintf("%s", reflect.TypeOf(value)) == StringType {
				strs = append(strs, fmt.Sprintf("N'%v'", strings.ReplaceAll(value.(string), "'", "''")))
			} else {
				strs = append(strs, fmt.Sprintf("N'%v'", value))
			}
		} else if value != "null" && (fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerTimeType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerStringType ||
			fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerBoolType) {
			if fmt.Sprintf("%s", reflect.TypeOf(value)) == PointerStringType {
				strValue := reflect.ValueOf(value).Elem().String()
				strs = append(strs, fmt.Sprintf("N'%v'", strings.ReplaceAll(strValue, "'", "''")))
			} else {
				strs = append(strs, fmt.Sprintf("N'%v'", reflect.ValueOf(value).Elem()))
			}
		} else {
			if strings.Contains(fmt.Sprintf("%s", reflect.TypeOf(value)), "*") {
				strs = append(strs, fmt.Sprintf("%v", reflect.ValueOf(value).Elem()))
			} else {
				strs = append(strs, fmt.Sprintf("%v", value))
			}
		}
	}

	return strs
}

func GenerateQueryColumnValues(s interface{}, ignoreFields []string) string {
	strs := ManipulateValues(s, ignoreFields)
	return strings.Join(strs, ",")
}

func GenerateQueryUpdateFields(s interface{}, ignoreFields []string) string {
	strs := []string{}

	columns := strings.Split(GenerateQueryColumnNames(s, ignoreFields), ",")

	values := ManipulateValues(s, ignoreFields)

	if columns[0] != "" {
		for i := 0; i < len(columns); i++ {
			strs = append(strs, fmt.Sprintf("%s = %v", columns[i], values[i]))
		}
	}

	return strings.Join(strs, ",")
}

func GenerateQueryUpdateFieldInclude(s interface{}, includeFields []string) string {
	t := reflect.TypeOf(s)
	ignoreFields := []string{}
	fields := []string{}

	if len(includeFields) == 0 {
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		if !funk.ContainsString(includeFields, t.Field(i).Name) {
			ignoreFields = append(ignoreFields, t.Field(i).Name)
		}
	}

	columns := strings.Split(GenerateQueryColumnNames(s, ignoreFields), ",")

	values := ManipulateValues(s, ignoreFields)

	if columns[0] != "" {
		for i := 0; i < len(columns); i++ {
			fields = append(fields, fmt.Sprintf("%s = %v", columns[i], values[i]))
		}
	}

	return strings.Join(fields, ",")
}

func GenerateTransactionAndRollback(queries []string) string {

	concatQuery := ""

	for i, query := range queries {
		concatQuery = concatQuery + query
		if i+1 == len(queries) {
			concatQuery += "\n"
		} else {
			concatQuery += "\n\t\n\t"
		}
	}

	return fmt.Sprintf(`
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		%s
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
	SELECT @ErrorMessage = N'Error %%d, Level %%d, State %%d, Procedure %%s, Line %%d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`,
		concatQuery,
	)
}

func GenerateTransactionAndRollbackWithRepeatableRead(queries []string) string {

	concatQuery := ""

	for i, query := range queries {
		concatQuery = concatQuery + query
		if i+1 == len(queries) {
			concatQuery += "\n"
		} else {
			concatQuery += "\n\t\n\t"
		}
	}

	return fmt.Sprintf(`
	SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
	BEGIN
	BEGIN TRY
	BEGIN TRANSACTION;
		%s
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
	SELECT @ErrorMessage = N'Error %%d, Level %%d, State %%d, Procedure %%s, Line %%d, ' + 'Message: ' + @ErrorMessage;
	RAISERROR (@ErrorMessage, @ErrorSeverity, 1, @ErrorNumber, @ErrorSeverity, @ErrorState, @ErrorProcedure, @ErrorLine)
	END CATCH
	END`,
		concatQuery,
	)
}
