package task_statistics

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetByTaskId(taskId int) (TaskStatistic, error) {
	// task with id 0 is default task
	if taskId == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	statement := fmt.Sprintf(`select t_res.*
from (
	select distinct t.language from tests t
	where t.final = true and t.task_id = $1
    ) t_groups
    join lateral (
        select
        	u.username,
			usr.compilation_time,
			usr.real_time,
			usr.kernel_time,
			usr.user_time,
			usr.max_ram_usage,
			usr.binary_size,
			t2.language
        from tests t2
        join user_solutions_results usr on usr.test_id = t2.id
        and t2.language = t_groups.language and t2.task_id = $2
		and usr.exit_code = 0
        join users u on u.id = t2.user_id
		order by usr.real_time asc
        limit 5
    ) t_res on true`)

	rows, err := postgres.GetPool().Query(postgres.GetCtx(), statement, taskId, taskId)
	if err != nil {
		return nil, err
	}

	var ts TaskStatistic
	for rows.Next() {
		se := &StatisticEntry{}
		if err = rows.Scan(&se.Username, &se.CompilationTime, &se.RealTime, &se.KernelTime,
			&se.UserTime, &se.MaxRamUsage, &se.BinarySize, &se.Language); err != nil {
			return nil, err
		}
		ts = append(ts, se)
	}
	return ts, nil
}
