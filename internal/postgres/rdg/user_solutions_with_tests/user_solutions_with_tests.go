package user_solutions_with_tests

type Test struct {
	LastModified string `json:"last_modified"`
	Final        bool   `json:"final"`
	Name         string `json:"name"`
	Public       bool   `json:"public"`
}

type Solution struct {
	LastModified string `json:"last_modified"`
	ExitCode     int    `json:"exit_code"`
	Name         string `json:"name"`
	Public       bool   `json:"public"`
}

type UserSolutionsWithTests struct {
	Tests     map[int]Test     `json:"tests"`
	Solutions map[int]Solution `json:"solutions"`
}

func NewUserSolutionsWithTests() *UserSolutionsWithTests {
	return &UserSolutionsWithTests{Tests: map[int]Test{}, Solutions: map[int]Solution{}}
}
