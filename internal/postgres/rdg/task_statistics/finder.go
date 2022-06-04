package task_statistics

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetByTaskId(tx postgres.PoolInterface, taskId int) (TaskStatistic, error) {
	// task with id 0 is default task
	if taskId == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	statement := fmt.Sprintf(`with all_fastest as (
			select t_res.*
			from (
				select distinct t.language from tests t
				where t.final = true and t.task_id = $1
				) t_groups
				join lateral (
					select
						u.username,
						row_number() over (partition by u.username order by usr.real_time asc) n,
						usr.compilation_time,
						usr.real_time,
						usr.kernel_time,
						usr.user_time,
						usr.max_ram_usage,
						usr.binary_size,
						t2.language
					from tests t2
					join user_solutions_results usr on usr.test_id = t2.id and t2.final = true
					and t2.language = t_groups.language and t2.task_id = $2
					and usr.exit_code = 0
					join users u on u.id = usr.user_id
				) t_res on true and n = 1
		),
		
		all_fastest_partitioned_by_language as (
			select
				all_fastest.*,
				row_number() over (partition by all_fastest.language order by all_fastest.real_time asc) l
			from all_fastest
		)
		
		select
			username,
			compilation_time,
			real_time,
			kernel_time,
			user_time,
			max_ram_usage,
			binary_size,
			language
		from all_fastest_partitioned_by_language where l between 1 and 5
		order by language asc, real_time asc`)

	rows, err := tx.Query(postgres.GetCtx(), statement, taskId, taskId)
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
