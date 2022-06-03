package tasks

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetById(tx postgres.PoolInterface, id int) (*Task, error) {
	return getByCondition(tx, id, "id = $1", id)
}

func GetByIdAndAuthorId(tx postgres.PoolInterface, id, authorId int) (*Task, error) {
	return getByCondition(tx, id, "id = $1 and author_id = $2", id, authorId)
}

func getByCondition(tx postgres.PoolInterface, id int, condition string, args ...interface{}) (*Task, error) {
	if id == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	statement := fmt.Sprintf("select %s from tasks where %s", allFieldReplacedTimestamp, condition)
	task := Task{}
	err := load(tx.QueryRow(postgres.GetCtx(), statement, args...), &task)
	return &task, err
}

func GetUnapproved(tx postgres.PoolInterface, by *FilterBy) ([]*Task, error) {
	condition := "is_published = true and approver_id = 0"
	return getBySearchBarFilers(tx, condition, []interface{}{}, by)
}

func GetByAuthorIdAndFilter(tx postgres.PoolInterface, userId int, by *FilterBy) ([]*Task, error) {
	condition := "author_id = $1"
	return getBySearchBarFilers(tx, condition, []interface{}{userId}, by)
}

func GetApprovedAndPublishedByFilter(tx postgres.PoolInterface, by *FilterBy) ([]*Task, error) {
	condition := "is_published = true and approver_id != 0"
	return getBySearchBarFilers(tx, condition, []interface{}{}, by)
}

func getBySearchBarFilers(tx postgres.PoolInterface, condition string, conditionArgs []interface{}, by *FilterBy) ([]*Task, error) {
	if by.Search != "" {
		condition += fmt.Sprintf(" and (strpos(lower(title), $%d) > 0 or strpos(lower(text), $%d) > 0)",
			len(conditionArgs)+1, len(conditionArgs)+2)
		conditionArgs = append(conditionArgs, by.Search, by.Search)
	}

	switch by.Difficulty {
	case "easy":
		condition += " and difficulty = 'easy'"
	case "medium":
		condition += " and difficulty = 'medium'"
	case "hard":
		condition += " and difficulty = 'hard'"
	}

	if by.NotPublished != "" {
		condition += " and is_published = false"
	}

	sort := "order by added_on desc"
	if by.Date == "asc" {
		sort = "order by added_on asc"
	}

	if by.Name == "desc" {
		sort += ", title desc"
	} else {
		sort += ", title asc"
	}

	perPage := 7
	sort += fmt.Sprintf(" limit %d offset %d", perPage, perPage*by.Page-perPage)

	return getManyWithConditions(tx, condition, sort, conditionArgs...)
}

func getManyWithConditions(tx postgres.PoolInterface, conditions, sort string, args ...interface{}) ([]*Task, error) {
	statement := fmt.Sprintf("select %s from tasks where %s %s", allFieldReplacedTimestamp, conditions, sort)
	rows, err := tx.Query(postgres.GetCtx(), statement, args...)
	if err != nil {
		return nil, err
	}
	return loadTasks(rows)
}

func loadTasks(rows pgx.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		task := Task{}
		if err := load(rows, &task); err != nil {
			return nil, err
		}
		if task.Id == 0 {
			continue
		}
		tasks = append(tasks, &task)
	}
	if len(tasks) == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	return tasks, nil
}

func load(qr pgx.Row, task *Task) error {
	err := qr.Scan(&task.Id, &task.AuthorId, &task.ApproverId, &task.Title,
		&task.Difficulty, &task.IsPublished, &task.AddedOn, &task.Text)
	return err
}

type FilterBy struct {
	Search       string
	Date         string
	Name         string
	Difficulty   string
	Page         int
	NotPublished string // false if empty
}
