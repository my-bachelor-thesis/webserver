package initial_data_for_editor

type InitialDataForEditor struct {
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Text       string   `json:"text"`
	AddedOn    string   `json:"added_on"`
	Author     string   `json:"author"`
	AuthorId   int      `json:"author_id"`
	Approver   string   `json:"approver"`
	Languages  []string `json:"languages"`
}
