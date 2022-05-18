package last_opened

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const allFields = "user_id, task_id, user_solution_id_for_language_1, language_1, user_solution_id_for_language_2, language_2"

type LastOpened struct {
	UserId                     int    `json:"user_id"`
	TaskId                     int    `json:"task_id" validate:"required"`
	UserSolutionIdForLanguage1 int    `json:"user_solution_id_for_language_1" validate:"required"`
	Language1                  string `json:"language_1" validate:"required"`
	UserSolutionIdForLanguage2 int    `json:"user_solution_id_for_language_2" validate:"required"`
	Language2                  string `json:"language_2" validate:"required"`
}

func (lo *LastOpened) Insert() error {
	statement := fmt.Sprintf(`
	insert into last_opened (%s)
	values (%s)`, allFields, postgresutil.GeneratePlaceholders(allFields))
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, lo.UserId, lo.TaskId,
		lo.UserSolutionIdForLanguage1, lo.Language1, lo.UserSolutionIdForLanguage2, lo.Language2)
	return err
}

func (lo *LastOpened) UpdateUserSolutionId() error {
	statement := `update last_opened set user_solution_id_for_language_1 = $1, language_1 = $2,
	user_solution_id_for_language_2 = $3, language_2 = $4
	where task_id = $5 and user_id = $6`
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, lo.UserSolutionIdForLanguage1,
		lo.Language1, lo.UserSolutionIdForLanguage2, lo.Language2, lo.TaskId, lo.UserId)
	return err
}

func (lo *LastOpened) Upsert() error {
	_, err := GetByUserIdAndTaskId(lo.UserId, lo.TaskId)
	if postgresutil.IsNoRowsInResultErr(err) {
		return lo.Insert()
	}
	return lo.UpdateUserSolutionId()
}
