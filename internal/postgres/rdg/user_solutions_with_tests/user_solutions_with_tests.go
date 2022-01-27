package user_solutions_with_tests

type Test struct {
	Date  string `json:"date"`
	Final bool   `json:"final"`
}

type Solution struct {
	Date     string `json:"date"`
	ExitCode int    `json:"exit_code"`
}

type UserSolutionsWithTests struct {
	Tests     map[int]Test     `json:"tests"`
	Solutions map[int]Solution `json:"solutions"`
}

func NewUserSolutionsWithTests() *UserSolutionsWithTests {
	return &UserSolutionsWithTests{Tests: map[int]Test{}, Solutions: map[int]Solution{}}
}
