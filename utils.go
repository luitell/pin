package main

func switchField(m model, field int) model {
	switch field {
	case 0:
		m.activeField = 0
		m.titleField.Focus()
		m.bodyField.Blur()

	case 1:
		m.activeField = 1
		m.bodyField.Focus()
		m.titleField.Blur()
	}
	return m
}
