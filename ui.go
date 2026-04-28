package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

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
		state:   "menu",
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
	if m.state == "menu" {
		return m.viewMenu()
	}
	return ""
}

func (m model) viewMenu() string {
	builder := strings.Builder{}
	builder.WriteString("===== Архиватор =====\n\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		builder.WriteString(fmt.Sprintf("%s %s \n", cursor, choice))
	}
	builder.WriteString("Используйте стрелки для навигации и enter для выбора\nНажмите q для выхода")
	return builder.String()
}
