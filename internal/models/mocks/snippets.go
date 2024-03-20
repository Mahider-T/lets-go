package mocks

import (
	"oogway/first/snippetbox/internal/models"
	"time"
)

var mockSnippet = &models.Snippet{
	Id:      1,
	Title:   "Title",
	Content: "Content",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 1, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
