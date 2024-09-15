package entities

type DoraemonGadget struct {
	EngName string `json:"eng_name"`
	JpName string `json:"jp_name"`
	Function string `json:"function"`
	AppearsIn []string `json:"appears_in"`
	ImageUrl string `json:"image_url"`
}