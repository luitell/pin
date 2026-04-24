package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

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
