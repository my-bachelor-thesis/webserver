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
	Tests     map[int]*tests.Test `json:"tests"`
	Solutions map[int]*Solutions  `json:"solutions"`
}

func NewUserSolutionsWithTests() *UserSolutionsWithTests {
	return &UserSolutionsWithTests{Tests: map[int]*tests.Test{}, Solutions: map[int]*Solutions{}}
}
