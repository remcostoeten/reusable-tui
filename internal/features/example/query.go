package example

import (
	"context"
	"database/sql"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const queryTimeout = 5 * time.Second

type Item struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt string
}

func listItems(db *sql.DB) ([]Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT id, title, done, created_at FROM example_items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func findItem(db *sql.DB, id int64) (Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	row := db.QueryRowContext(ctx, `SELECT id, title, done, created_at FROM example_items WHERE id = ?`, id)
	var item Item
	err := row.Scan(&item.ID, &item.Title, &item.Done, &item.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Item{}, store.ErrNotFound
	}
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func loadItems(db *sql.DB) tea.Cmd {
	return func() tea.Msg { return itemsResult(db) }
}

func itemsResult(db *sql.DB) tea.Msg {
	items, err := listItems(db)
	if err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return itemsLoadedMsg{items: items}
}
