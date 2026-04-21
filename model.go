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
	activeIndex int
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
		activeIndex: 0,
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

	if m.viewState == createView {
		switch m.activeField {
		case 0:
			m.titleField, cmd = m.titleField.Update(msg)
			cmds = append(cmds, cmd)
		case 1:
			m.bodyField, cmd = m.bodyField.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

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
		m.titleField.SetValue("")
		m.bodyField.SetValue("")
		m = switchField(m, 0)
		m.viewState = createView
		return m, nil

	case "d":

		task := m.tasks[m.activeIndex]
		err := m.store.DeleteTask(int(task.ID))
		if err != nil {
			log.Fatal(err)
		}
		m.tasks = append(m.tasks[:m.activeIndex], m.tasks[m.activeIndex+1:]...)
		if m.activeIndex > len(m.tasks)-1 {
			m.activeIndex = len(m.tasks) - 1
		}

		return m, nil
	case "j":
		if m.activeIndex < len(m.tasks)-1 {
			m.activeIndex += 1
			return m, nil
		}
	case "k":
		if m.activeIndex > 0 {
			m.activeIndex -= 1
			return m, nil
		}
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

	case "tab":
		if m.activeField < 1 {
			m = switchField(m, m.activeField+1)
		} else if m.activeField > 0 {
			m = switchField(m, m.activeField-1)
		}

	case "enter":
		if m.activeField == 0 {
			if m.titleField.Value() == "" {
				return m, nil
			}
			m = switchField(m, 1)
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
		if m.activeField == -1 {
			m.activeField = 0
		}
		// Reset fields
		m.titleField.SetValue("")
		m.bodyField.SetValue("")
		m.activeField = 0

		m.viewState = listView
		return m, nil
	}
	return m, nil
}
