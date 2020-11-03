package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rknizzle/faas/internal/function"
	"github.com/sirupsen/logrus"
)

// interacts with the functions database table
type functionRepo struct {
	Conn *sql.DB
}

func (repo *functionRepo) GetByName(ctx context.Context, name string) (fxn function.Function, err error) {
	query := `SELECT * FROM functions WHERE name = '$1'`

	list, err := repo.fetch(ctx, query, name)
	if err != nil {
		return
	}

	if len(list) > 0 {
		fxn = list[0]
	} else {
		return fxn, errors.New("Function not found")
	}
	return
}

// convert the result of a sql query into function.Function instances
func (r *functionRepo) fetch(ctx context.Context, query string, args ...interface{}) (result []function.Function, err error) {

	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]function.Function, 0)
	for rows.Next() {
		t := function.Function{}
		err = rows.Scan(
			&t.Name,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
