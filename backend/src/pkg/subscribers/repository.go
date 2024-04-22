package subscribers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"newsletter/src/pkg/entity"
	newsletterError "newsletter/src/pkg/utils/error"
	"newsletter/src/pkg/utils/logger"
	sqlQuery "newsletter/src/pkg/utils/sqlquery"

	sqlStruct "github.com/kisielk/sqlstruct"
)

type Repository interface {
	GetAllSubscribers() ([]entity.Subscribers, error)
	FindByEmail(email string) ([]entity.Subscribers, error)
	Insert(subscriber entity.Subscribers) error
	UpdateByEmail(subscriber entity.Subscribers) error
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
		go repo.Logs.Error("", "subscribers_Repo_GetAllSubscribers", "", newsletterError.NewError(newsletterError.TechnicalError, err.Error()))
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

func (repo *SqlRepository) FindByEmail(email string) ([]entity.Subscribers, error) {
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
	AND Email = N'%[3]s'
	`,
		sqlQuery.GenerateQueryColumnNames(entity.Subscribers{}, []string{}),
		repo.Collection,
		email,
	)
	rows, err := session.QueryContext(ctx, sql)
	if err != nil {
		go repo.Logs.Error("", "subscribers_Repo_FindByEmail", "", newsletterError.NewError(newsletterError.TechnicalError, err.Error()))
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

func (repo *SqlRepository) Insert(subscriber entity.Subscribers) error {
	ctx := context.Background()
	session, sessionErr := repo.Session.Conn(ctx)
	if sessionErr != nil {
		return sessionErr
	}
	defer session.Close()

	sql := fmt.Sprintf(`
	INSERT INTO %[1]s
	(
		%[2]s,
		[IsSubscribed],
		[SubscribedDate],
		[UnsubscribedDate],
		[Delflag]
	)
	VALUES
	(
		%[3]s,
		1,
		GETDATE(),
		Null,
		0
	)
	`,
		repo.Collection,
		sqlQuery.GenerateQueryColumnNames(entity.Subscribers{}, []string{"ID", "IsSubscribed", "SubscribedDate", "UnsubscribedDate", "DelFlag"}),
		sqlQuery.GenerateQueryColumnValues(subscriber, []string{"ID", "IsSubscribed", "SubscribedDate", "UnsubscribedDate", "DelFlag"}),
	)

	_, err := session.ExecContext(ctx, sql)
	if err != nil {
		go repo.Logs.Error("", "subscriber_Repo_Insert", subscriber,
			newsletterError.NewError(newsletterError.TechnicalError, err.Error()))
		return err
	}

	return nil
}

func (repo *SqlRepository) UpdateByEmail(subscriber entity.Subscribers) error {

	ctx := context.Background()
	session, sessionErr := repo.Session.Conn(ctx)
	if sessionErr != nil {
		return sessionErr
	}
	defer session.Close()

	setDate := `IsSubscribed = 1,
	SubscribedDate = GETDATE(),`

	if !subscriber.IsSubscribed {
		setDate = `IsSubscribed = 0,
		UnsubscribedDate = GETDATE()`
	}

	sql := fmt.Sprintf(`
	UPDATE %[1]s
	SET 
		%[2]s,
		%[4]s
	WHERE Email = N'%[3]s'
	`,
		repo.Collection,
		sqlQuery.GenerateQueryUpdateFields(subscriber, []string{"ID", "SubscribedDate", "UnsubscribedDate", "IsSubscribed", "DelFlag"}),
		subscriber.Email,
		setDate,
	)

	_, err := session.ExecContext(ctx, sql)
	if err != nil {
		return err
	}
	return nil
}
