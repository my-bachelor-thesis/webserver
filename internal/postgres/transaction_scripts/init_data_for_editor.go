package transaction_scripts

import (
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/initial_data_for_editor"
)

func GetInitDataForEditorByTaskId(taskId int) (*initial_data_for_editor.InitialDataForEditor, error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	statement := `
	select
		t.title,
		t.difficulty,
		t.text,
		t.author_id,
		to_char(t.added_on, 'DD.MM.YY, HH24:MI:SS'),
		(select u.first_name || ' ' || u.last_name from users u where u.id = t.author_id) author,
		(select u.first_name || ' ' || u.last_name from users u where u.id = t.approver_id) approver,
		t.approver_id
	from tasks t where id = $1`
	initData := initial_data_for_editor.InitialDataForEditor{}
	if err := tx.QueryRow(postgres.GetCtx(), statement, taskId).Scan(&initData.Title, &initData.Difficulty, &initData.Text,
		&initData.AuthorId, &initData.AddedOn, &initData.Author, &initData.Approver, &initData.ApproverId); err != nil {
		return nil, err
	}

	languagesStatement := `
	select
		language
	from tests where final = true and task_id = $1`
	rows, err := tx.Query(postgres.GetCtx(), languagesStatement, taskId)
	if err != nil {
		return nil, err
	}
	var row string
	for rows.Next() {
		if err = rows.Scan(&row); err != nil {
			return nil, err
		}
		initData.Languages = append(initData.Languages, row)
	}

	err = tx.Commit(postgres.GetCtx())

	return &initData, err
}
