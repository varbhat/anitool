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

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
		return textinput.Blink
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
	}
	switch m.uiState {
	case uiMainPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				items := []list.Item{
					item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
					item{title: "Nutella", desc: "It's good on toast"},
					item{title: "Bitter melon", desc: "It cools you down"},
					item{title: "Nice socks", desc: "And by that I mean socks without holes"},
					item{title: "Eight hours of sleep", desc: "I had this once"},
					item{title: "Cats", desc: "Usually"},
					item{title: "Plantasia, the album", desc: "My plants love it too"},
					item{title: "Pour over coffee", desc: "It takes forever to make though"},
					item{title: "VR", desc: "Virtual reality...what is there to say?"},
					item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
					item{title: "Linux", desc: "Pretty much the best OS"},
					item{title: "Business school", desc: "Just kidding"},
					item{title: "Pottery", desc: "Wet clay is a great feeling"},
					item{title: "Shampoo", desc: "Nothing like clean hair"},
					item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
					item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
					item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
					item{title: "Stickers", desc: "The thicker the vinyl the better"},
					item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
					item{title: "Warm light", desc: "Like around 2700 Kelvin"},
					item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
					item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
					item{title: "Terrycloth", desc: "In other words, towel fabric"},
				}

				m.animeList = list.New(items, list.NewDefaultDelegate(), 0, 0)
				m.animeList.Title = "My Fave Things"
				m.uiState = uiAnimeListPage
				var cmd tea.Cmd
				m.animeList, cmd = m.animeList.Update(msg)
				top, right, bottom, left := docStyle.GetMargin()
				m.animeList.SetSize(m.WindowSize.Width-left-right, m.WindowSize.Height-top-bottom)
				return m, cmd
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
