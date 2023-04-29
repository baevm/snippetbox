package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetRepo interface {
	Create(title, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Create(title, content string, expires int) (int, error) {
	// ? used as placeholder to avoid SQL injections
	query := `
	INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`
	res, err := s.DB.Exec(query, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (*Snippet, error) {
	snip := &Snippet{}

	query := `
	SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?
	`
	err := s.DB.
		QueryRow(query, id).
		Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snip, nil
}

func (s *SnippetModel) Latest() ([]*Snippet, error) {
	snippets := []*Snippet{}

	query := `
	SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() 
	ORDER BY id DESC LIMIT 10
	`

	rows, err := s.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		snip := &Snippet{}

		err := rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snip)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (s *SnippetModel) Update(title, content string, expires int) (int, error) {
	return 0, nil
}
