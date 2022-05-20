package task_statistics

import "webserver/internal/postgres/rdg/user_solutions_results"

type StatisticEntry struct {
	Username string `json:"username"`
	Language string `json:"language"`
	user_solutions_results.InfoForStatistic
}

type TaskStatistic []*StatisticEntry
