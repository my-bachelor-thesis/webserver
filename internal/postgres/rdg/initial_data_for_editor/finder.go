package initial_data_for_editor

import (
	"webserver/internal/postgres"
)

func GetByTaskId(taskId int) (*InitialDataForEditor, error) {
	statement := `
	select
		t.title,
		t.difficulty,
		t.text,
		to_char(t.added_on, 'DD.MM.YY, HH24:MI:SS'),
		(select u.first_name || ' ' || u.last_name from users u where u.id = t.author_id) author,
		(select u.first_name || ' ' || u.last_name from users u where u.id = t.approver_id) approver
	from tasks t where id = $1`
	initData := InitialDataForEditor{}
	if err := postgres.GetPool().QueryRow(postgres.GetCtx(), statement, taskId).Scan(&initData.Title, &initData.Difficulty,
		&initData.Text, &initData.AddedOn, &initData.Author, &initData.Approver); err != nil {
		return nil, err
	}

	languagesStatement := `
	select
		language
	from tests where final = true and task_id = $1`
	rows, err := postgres.GetPool().Query(postgres.GetCtx(), languagesStatement, taskId)
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

	return &initData, err
}
