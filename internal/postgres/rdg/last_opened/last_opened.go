package last_opened

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const allFields = "user_id, task_id, user_solution_id_for_language_1, language_1, user_solution_id_for_language_2, language_2"

type LastOpened struct {
	UserId                     int    `json:"user_id"`
	TaskId                     int    `json:"task_id"`
	UserSolutionIdForLanguage1 int    `json:"user_solution_id_for_language_1"`
	Language1                  string `json:"language_1"`
	UserSolutionIdForLanguage2 int    `json:"user_solution_id_for_language_2"`
	Language2                  string `json:"language_2"`
}

func (lo *LastOpened) Insert(tx postgres.PoolInterface) error {
	statement := fmt.Sprintf(`
	insert into last_opened (%s)
	values (%s)`, allFields, postgresutil.GeneratePlaceholders(allFields))
	_, err := tx.Exec(postgres.GetCtx(), statement, lo.UserId, lo.TaskId,
		lo.UserSolutionIdForLanguage1, lo.Language1, lo.UserSolutionIdForLanguage2, lo.Language2)
	return err
}

func (lo *LastOpened) UpdateUserSolutionId(tx postgres.PoolInterface) error {
	statement := `update last_opened set user_solution_id_for_language_1 = $1, language_1 = $2,
	user_solution_id_for_language_2 = $3, language_2 = $4
	where task_id = $5 and user_id = $6`
	_, err := tx.Exec(postgres.GetCtx(), statement, lo.UserSolutionIdForLanguage1,
		lo.Language1, lo.UserSolutionIdForLanguage2, lo.Language2, lo.TaskId, lo.UserId)
	return err
}

func (lo *LastOpened) Upsert(tx postgres.PoolInterface) error {
	_, err := GetByUserIdAndTaskId(tx, lo.UserId, lo.TaskId)
	if postgresutil.IsNoRowsInResultErr(err) {
		return lo.Insert(tx)
	}
	return lo.UpdateUserSolutionId(tx)
}
