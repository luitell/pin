package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

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
