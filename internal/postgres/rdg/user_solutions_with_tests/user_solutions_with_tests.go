package user_solutions_with_tests

import (
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
)

type Solutions struct {
	user_solutions.UserSolution
	TestId int `json:"test_id"`
}

type UserSolutionsWithTests struct {
	Tests     []*tests.Test `json:"tests"`
	Solutions []*Solutions  `json:"solutions"`
}