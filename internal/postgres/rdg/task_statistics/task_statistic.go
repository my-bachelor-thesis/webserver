package task_statistics

import "webserver/internal/postgres/rdg/user_solutions_results"

type StatisticEntry struct {
	SolutionId int    `json:"solution_id"`
	Username   string `json:"username"`
	Language   string `json:"language"`
	user_solutions_results.InfoForStatistic
}

type TaskStatistic []*StatisticEntry
