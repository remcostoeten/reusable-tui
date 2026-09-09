package example

import (
	"context"
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func createItem(db *sql.DB, title string) (Item, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return Item{}, store.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	result, err := db.ExecContext(ctx, `INSERT INTO example_items (title) VALUES (?)`, trimmed)
	if err != nil {
		return Item{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Item{}, err
	}
	return findItem(db, id)
}

func toggleItem(db *sql.DB, id int64) (Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	result, err := db.ExecContext(ctx, `UPDATE example_items SET done = 1 - done WHERE id = ?`, id)
	if err != nil {
		return Item{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Item{}, err
	}
	if affected == 0 {
		return Item{}, store.ErrNotFound
	}
	return findItem(db, id)
}

func deleteItem(db *sql.DB, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	result, err := db.ExecContext(ctx, `DELETE FROM example_items WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return store.ErrNotFound
	}
	return nil
}

func addItem(db *sql.DB, title string) tea.Cmd {
	return func() tea.Msg { return addResult(db, title) }
}

func addResult(db *sql.DB, title string) tea.Msg {
	item, err := createItem(db, title)
	if err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return itemChangedMsg{action: "added", title: item.Title}
}

func flipItem(db *sql.DB, id int64) tea.Cmd {
	return func() tea.Msg { return flipResult(db, id) }
}

func flipResult(db *sql.DB, id int64) tea.Msg {
	item, err := toggleItem(db, id)
	if err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return itemChangedMsg{action: "toggled", title: item.Title}
}

func removeItem(db *sql.DB, id int64, title string) tea.Cmd {
	return func() tea.Msg { return removeResult(db, id, title) }
}

func removeResult(db *sql.DB, id int64, title string) tea.Msg {
	if err := deleteItem(db, id); err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return itemChangedMsg{action: "removed", title: title}
}
