package task_with_solutions_and_tests

type NameAndCode struct {
	Name     string `json:"name"`
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type TaskWithSolutionsAndTests struct {
	TaskId          int            `json:"task_id"`
	IsPublished     bool           `json:"is_published"`
	Title           string         `json:"title" validate:"required"`
	Difficulty      string         `json:"difficulty" validate:"required"`
	Description     string         `json:"description" validate:"required"`
	FinalTests      []*NameAndCode `json:"final_tests" validate:"required"`
	PublicTests     []*NameAndCode `json:"public_tests"`
	PublicSolutions []*NameAndCode `json:"public_solutions"`
}
