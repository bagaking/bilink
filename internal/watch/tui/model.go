package tui

type Model struct {
	Status string
}

func NewModel() Model {
	return Model{Status: "starting"}
}
