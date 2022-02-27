package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type uiState int

const (
	uiMainPage uiState = iota
	uiAnimeListPage
	uiSettingsPage
	uiEpisodePage
	uiPlayPage
)

type model struct {
	uiState    uiState
	textInput  textinput.Model
	animeList  list.Model
	err        error
	WindowSize tea.WindowSizeMsg
}

func mainModel() model {
	ti := textinput.New()
	ti.Focus()
	//ti.CharLimit = 156
	//ti.Width = 20

	return model{
		uiState:   uiMainPage,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	switch m.uiState {
	case uiMainPage:
		return tea.Batch(tea.EnterAltScreen, textinput.Blink)
	case uiAnimeListPage:
		return nil
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.WindowSize = msg
	case chan GSRes:
		itemsg, notclosed := <-msg
		if !notclosed {
			return m, nil
		}
		m.animeList.InsertItem(-1, itemsg)
		return m, func() tea.Msg { return msg }
	}
	switch m.uiState {
	case uiMainPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				items := []list.Item{}
				m.animeList = list.New(items, list.NewDefaultDelegate(), 0, 0)
				m.animeList.Title = "Anime Search Results"
				m.uiState = uiAnimeListPage
				top, right, bottom, left := docStyle.GetMargin()
				m.animeList.SetSize(m.WindowSize.Width-left-right, m.WindowSize.Height-top-bottom)
				return m, func() tea.Msg {
					return searchGogoAll("https://gogoanime.fi", m.textInput.Value())
				}
			}
		}

		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	case uiAnimeListPage:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			top, right, bottom, left := docStyle.GetMargin()
			m.animeList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		}

		var cmd tea.Cmd
		m.animeList, cmd = m.animeList.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	switch m.uiState {
	case uiMainPage:
		return fmt.Sprintf(
			"Search: \n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	case uiAnimeListPage:
		return docStyle.Render(m.animeList.View())
	case uiPlayPage:
		return "hello"
	default:
		return ""
	}
}
