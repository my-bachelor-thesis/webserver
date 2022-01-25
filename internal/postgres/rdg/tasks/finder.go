package tasks

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetById(id int) (*Task, error) {
	statement := fmt.Sprintf("select %s from tasks where id = $1", allFieldReplacedTimestamp)
	task := Task{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, id), &task)
	return &task, err
}

func GetAllApprovedAndPublished() ([]*Task, error) {
	statement := fmt.Sprintf("select %s from tasks where is_published=true and is_approved=true", allFieldReplacedTimestamp)
	rows, err := postgres.GetPool().Query(postgres.GetCtx(), statement)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for rows.Next() {
		task := Task{}
		if err = load(rows, &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, err
}

func load(qr pgx.Row, task *Task) error {
	err := qr.Scan(&task.Id, &task.AuthorId, &task.ApproverId, &task.FinalTestId, &task.Title, &task.Difficulty,
		&task.Description, &task.IsPublished, &task.IsApproved, &task.AddedOn, &task.Text)
	return err
}
