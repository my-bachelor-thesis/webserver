package tasks

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "author_id, approver_id, title, difficulty, is_published, added_on, text"
	allFields          = "id, " + allFieldsWithoutId
)

var allFieldReplacedTimestamp = postgresutil.CallToCharOnTimestamp(allFields, "added_on")

type Task struct {
	Id          int    `json:"id"` // 0 is default
	AuthorId    int    `json:"author_id"`
	ApproverId  int    `json:"approver_id"` // 0 means not approved
	Title       string `json:"title"`
	Difficulty  string `json:"difficulty"`
	IsPublished bool   `json:"is_published"`
	AddedOn     string `json:"added_on"`
	Text        string `json:"text"`
}

func (task *Task) Insert(tx postgres.PoolInterface) error {
	placeholders := postgresutil.GeneratePlaceholdersAndReplace(allFieldsWithoutId, map[int]string{5: "CURRENT_TIMESTAMP"})
	statement := fmt.Sprintf(`
	insert into tasks (%s)
	values (%s)
	returning id`, allFieldsWithoutId, placeholders)
	return tx.QueryRow(postgres.GetCtx(), statement, task.AuthorId, task.ApproverId, task.Title,
		task.Difficulty, task.IsPublished, task.Text).Scan(&task.Id)
}

func (task *Task) Publish(tx postgres.PoolInterface) error {
	statement := "update tasks set is_published = true where id = $1"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.Id)
	return err
}

func (task *Task) Unpublish(tx postgres.PoolInterface) error {
	statement := "update tasks set is_published = false, approver_id = 0 where id = $1"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.Id)
	return err
}

func (task *Task) Approve(tx postgres.PoolInterface) error {
	statement := "update tasks set approver_id = $1, added_on = CURRENT_TIMESTAMP where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.ApproverId, task.Id)
	return err
}

func (task *Task) Unapprove(tx postgres.PoolInterface) error {
	statement := "update tasks set approver_id = 0, is_published = false where id = $1"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.Id)
	return err
}

func (task *Task) ApproveAndPublish(tx postgres.PoolInterface) error {
	statement := "update tasks set approver_id = $1, is_published = true where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.ApproverId, task.Id)
	return err
}

func (task *Task) UpdateTitleDifficultyDateAndText(tx postgres.PoolInterface) error {
	statement := "update tasks set title = $1, difficulty = $2, text = $3, added_on = CURRENT_TIMESTAMP where id = $4 and author_id = $5"
	_, err := tx.Exec(postgres.GetCtx(), statement,
		task.Title, task.Difficulty, task.Text, task.Id, task.AuthorId)
	return err
}

func (task *Task) Delete(tx postgres.PoolInterface) error {
	statement := "delete from tasks where id = $1"
	_, err := tx.Exec(postgres.GetCtx(), statement, task.Id)
	return err
}
