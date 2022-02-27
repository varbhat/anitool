package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type uiState int

const (
	uiMainPage uiState = iota
	uiSearchPage
	uiSettingsPage
	uiEpisodePage
	uiPlayPage
)

type model struct {
	uiState   uiState
	textInput textinput.Model
	err       error
}

func mainModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.uiState {
	case uiMainPage, uiSearchPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}
		}

		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	switch m.uiState {
	case uiMainPage, uiSearchPage:
		return fmt.Sprintf(
			"Search Anime: \n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	case uiPlayPage:
		return "hello"
	default:
		return ""
	}
}
