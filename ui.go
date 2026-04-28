package main

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	archiver  *SimpleArchiver
	cursor    int
	choices   []string
	state     string
	err       error
	inputPath string
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Сжать файл", "Распаковать файл", "Выход"},
	}
}
func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m model) View() string {
	return ""
}
