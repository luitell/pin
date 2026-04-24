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
		m.height = msg.Height
		m.width = msg.Width
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
