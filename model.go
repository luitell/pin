package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	createView
	timerView
)

type model struct {
	viewState uint
	store     *Store
	tasks     []Task
}

func NewModel(store *Store) model {
	tasks, err := store.GetTasks()
	if err != nil {
		log.Fatalf("Could not fetch tasks %v", err)
	}
	return model{
		viewState: listView,
		store:     store,
		tasks:     tasks,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		switch m.viewState {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}
