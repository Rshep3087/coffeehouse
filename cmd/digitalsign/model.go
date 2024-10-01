package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rshep3087/coffeehouse/postgres"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	msgChan chan postgres.Recipe
	recipes []postgres.Recipe
	table   table.Model
}

// newModel creates a new model for the tea program.
func newModel(msgChan chan postgres.Recipe) model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Brew Method", Width: 20},
		{Title: "Coffee (g)", Width: 10},
		{Title: "Water (g)", Width: 10},
		{Title: "Temp (°F)", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{msgChan: msgChan, table: t}
}

func (m model) Init() tea.Cmd {
	return waitForRecipe(m.msgChan)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case postgres.Recipe:
		m.recipes = append(m.recipes, msg)

		rows := make([]table.Row, len(m.recipes))
		for i, r := range m.recipes {
			rows[i] = table.Row{
				fmt.Sprintf("%d", r.ID),
				r.RecipeName,
				string(r.BrewMethod),
				fmt.Sprintf("%.2f", r.CoffeeWeight),
				fmt.Sprintf("%.2f", r.WaterWeight),
				fmt.Sprintf("%.2f", r.WaterTemp.Float64),
			}
		}

		m.table.SetRows(rows)

		m.table, cmd = m.table.Update(msg)
		return m, tea.Batch(waitForRecipe(m.msgChan), cmd)
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func waitForRecipe(msgChannel chan postgres.Recipe) tea.Cmd {
	return func() tea.Msg {
		r := <-msgChannel
		return r
	}
}
