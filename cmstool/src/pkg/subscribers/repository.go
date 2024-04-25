package subscribers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"subscribetool/src/pkg/entity"
	subscribetoolError "subscribetool/src/pkg/utils/error"
	"subscribetool/src/pkg/utils/logger"
	sqlQuery "subscribetool/src/pkg/utils/sqlquery"

	sqlStruct "github.com/kisielk/sqlstruct"
)

type Repository interface {
	GetAllSubscribers() ([]entity.Subscribers, error)
}

type SqlRepository struct {
	Collection string
	Session    *sql.DB
	Logs       logger.Logger
}

func NewRepository(collection string, session *sql.DB, logs logger.Logger) *SqlRepository {
	return &SqlRepository{
		Collection: collection,
		Session:    session,
		Logs:       logs,
	}
}

func (repo *SqlRepository) GetAllSubscribers() ([]entity.Subscribers, error) {
	ctx := context.Background()
	session, sessionErr := repo.Session.Conn(ctx)
	if sessionErr != nil {
		return nil, sessionErr
	}
	defer session.Close()

	sql := fmt.Sprintf(`
	SELECT %[1]s 
	FROM %[2]s
	WHERE Delflag = 0
	AND IsSubscribed = 1
	`,
		sqlQuery.GenerateQueryColumnNames(entity.Subscribers{}, []string{}),
		repo.Collection,
	)
	rows, err := session.QueryContext(ctx, sql)
	if err != nil {
		go repo.Logs.Error("", "subscribers_Repo_GetAllSubscribers", "", subscribetoolError.NewError(subscribetoolError.TechnicalError, err.Error()))
		return nil, err
	}

	list := []entity.Subscribers{}
	for rows.Next() {
		var entity entity.Subscribers
		if err := sqlStruct.Scan(&entity, rows); err != nil {
			log.Println(err.Error())
		}
		list = append(list, entity)
	}
	return list, nil
}
