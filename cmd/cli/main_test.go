package main

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestCliModel(t *testing.T) {
	// Create a simple model for testing
	m := Model{
		messages: []string{"Welcome message"},
		input:    "",
	}

	// Test initial view
	view := m.View()
	assert.Contains(t, view, "Go Anon Kode")
	assert.Contains(t, view, "Welcome message")

	// Test input handling
	m.input = "test input"
	view = m.View()
	assert.Contains(t, view, "test input")

	// Test update with key press
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updatedModel, ok := newModel.(Model)
	assert.True(t, ok)
	assert.Equal(t, "test inputa", updatedModel.input)

	// Test update with backspace
	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	updatedModel, ok = newModel.(Model)
	assert.True(t, ok)
	assert.Equal(t, "test inpu", updatedModel.input)

	// Test quit command
	newModel, cmd := updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.NotNil(t, cmd)
}

// MockProgram is a mock implementation of the bubbletea.Program interface for testing
type MockProgram struct {
	output bytes.Buffer
}

func (m *MockProgram) Run() (tea.Model, error) {
	return nil, nil
}

func (m *MockProgram) Send(msg tea.Msg) {
	// Do nothing
}

func (m *MockProgram) Quit() {
	// Do nothing
}
