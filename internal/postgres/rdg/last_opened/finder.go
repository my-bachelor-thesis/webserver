package last_opened

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetByUserIdAndTaskId(userId int, taskId int) (*LastOpened, error) {
	statement := fmt.Sprintf("select %s from last_opened where user_id = $1 and task_id = $2", allFields)
	lo := LastOpened{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, userId, taskId), &lo)
	return &lo, err
}

func load(qr pgx.Row, lo *LastOpened) error {
	return qr.Scan(&lo.UserId, &lo.TaskId, &lo.UserSolutionIdForLanguage1, &lo.Language1,
		&lo.UserSolutionIdForLanguage2, &lo.Language2)
}
