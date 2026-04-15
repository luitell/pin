package main

import (
	"log"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	createView
	timerView
)

type model struct {
	viewState   uint
	width       int
	height      int
	activeField int
	store       *Store
	tasks       []Task
	titleField  textinput.Model
	bodyField   textinput.Model
}

func NewModel(store *Store) model {
	tasks, err := store.GetTasks()
	if err != nil {
		log.Fatalf("Could not fetch tasks %v", err)
	}
	return model{
		viewState:   listView,
		store:       store,
		tasks:       tasks,
		titleField:  textinput.New(),
		bodyField:   textinput.New(),
		activeField: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.titleField, cmd = m.titleField.Update(msg)
	cmds = append(cmds, cmd)

	m.bodyField, cmd = m.bodyField.Update(msg)

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Width
		m.width = msg.Height
	case tea.KeyPressMsg:
		key := msg.String()
		switch m.viewState {
		case listView:
			return ListViewActions(key, m)
		case createView:
			return CreateViewActions(key, m)
		}
	}
	return m, tea.Batch(cmds...)
}

func ListViewActions(key string, m model) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "a":
		m.titleField.Focus()
		m.viewState = createView
		return m, nil
	}
	return m, nil
}

func CreateViewActions(key string, m model) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.viewState = listView
		return m, nil

	case "enter":
		if m.activeField == 0 {
			// Move from title → body
			if m.titleField.Value() == "" {
				return m, nil // don't proceed if title empty
			}
			m.activeField = 1
			m.bodyField.Focus()
			m.titleField.Blur()
			return m, nil
		}

		// Submit when on body field
		title := m.titleField.Value()
		body := m.bodyField.Value()

		task, err := m.store.SaveTask(title, body)
		if err != nil {
			log.Fatal(err)
		}

		m.tasks = append(m.tasks, task)

		// Reset fields
		m.titleField.SetValue("")
		m.bodyField.SetValue("")
		m.activeField = 0

		m.viewState = listView
		return m, nil
	}
	return m, nil
}
