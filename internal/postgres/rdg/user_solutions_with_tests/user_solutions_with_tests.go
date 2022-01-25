package user_solutions_with_tests

type Test struct {
	Date string `json:"date"`
}

type Solution struct {
	Date     string `json:"date"`
	ExitCode int    `json:"exit_code"`
}

type UserSolutionsWithTests struct {
	FinalTest int              `json:"final_test"`
	Tests     map[int]Test     `json:"tests"`
	Solutions map[int]Solution `json:"solutions"`
}
