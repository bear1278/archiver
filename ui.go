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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == "menu" {
			return m.updateMenu(msg)
		}
	}
	return m, nil
}
func (m model) View() string {
	if m.state == "menu" {
		return m.viewMenu()
	}
	if m.state == "compress" || m.state == "decompress" {
		return m.viewInput()
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

func (m model) viewInput() string {
	builder := strings.Builder{}
	builder.WriteString("===== Архиватор =====\n\n")
	builder.WriteString("Введите к файлу для ")
	if m.state == "compress" {
		builder.WriteString("сжатия:\n")
	} else if m.state == "decompress" {
		builder.WriteString("распаковки:\n")
	}
	builder.WriteString(m.inputPath)
	builder.WriteString("_ \n")
	if m.err != nil {
		builder.WriteString(m.err.Error())
	}
	builder.WriteString("Enter для подтверждения, Esc для возврата в меню")
	return builder.String()
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
		return m, nil
	case "enter":
		switch m.cursor {
		case 0:
			m.state = "compress"
			return m, nil
		case 1:
			m.state = "decompress"
			return m, nil
		case 2:
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}
