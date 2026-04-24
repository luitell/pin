package main

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	s := "Welcome to pin\n\n"
	switch m.viewState {
	case listView:
		for index, task := range m.tasks {
			if index == m.activeIndex {
				s += "> "
			}
			s += task.Title
			if task.Description != "" {
				s += "\n" + task.Description
			}
			s += "\n\n"
		}
	case createView:
		s += "Task Title\n"
		s += m.titleField.View() + "\n"
		s += "Task Description (optional)\n"
		s += m.bodyField.View() + "\n"
	}
	s += "q, ctrl+c->quit\na->new task"
	return tea.NewView(s)
}
