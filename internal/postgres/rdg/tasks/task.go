package tasks

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "author_id, approver_id, title, difficulty, description, is_published, is_approved, added_on, text"
	allFields          = "id, " + allFieldsWithoutId
)

var allFieldReplacedTimestamp = postgresutil.CallToCharOnTimestamp(allFields, "added_on")

type Task struct {
	Id          int    `json:"id"`
	AuthorId    int    `json:"author_id"`
	ApproverId  int    `json:"approver_id"`
	Title       string `json:"title"`
	Difficulty  string `json:"difficulty"`
	Description string `json:"description"`
	IsPublished bool   `json:"is_published"`
	IsApproved  bool   `json:"is_approved"`
	AddedOn     string `json:"added_on"`
	Text        string `json:"text"`
}

func Insert(task *Task) error {
	statement := fmt.Sprintf(`
	insert into tasks %s
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholders(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, task.AuthorId, task.ApproverId, task.Title,
		task.Difficulty, task.Description, task.IsPublished, task.IsApproved, task.AddedOn, task.Text).Scan(&task.Id)
}
