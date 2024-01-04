package sqlite

import (
	"L2-task11/internal/storage"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(
		`
		CREATE TABLE IF NOT EXISTS event(
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    event_date DATE NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_uesr ON event(user_id); 
		`) // мб нужно создать составной индекс из сразу из двух полей, общий хэш или что-то еще
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveEvent(userId int, eventDate string) error {
	const fn = "storage.sqlite.SaveEvent"

	// check if already exist
	q := `SELECT COUNT(*) FROM event WHERE user_id = ? AND event_date = ?;`
	var count int
	err := s.db.QueryRow(q, userId, eventDate).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	if count > 0 {
		return storage.ErrEventExists
	}

	stmt, err := s.db.Prepare("INSERT INTO event(user_id, event_date) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(userId, eventDate)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) UpdateEvent(userId int, eventDate string) error {
	const fn = "storage.sqlite.UpdateEvent"

	// if doesnt exist -> err
	q := `SELECT COUNT(*) FROM event WHERE user_id = ? AND event_date = ?;`
	var count int
	err := s.db.QueryRow(q, userId, eventDate).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	if count == 0 {
		return storage.ErrEventNotExists
	}

	stmt, err := s.db.Prepare("UPDATE event SET user_id = ?, event_date = ? WHERE user_id = ? AND event_date = ?;")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(userId, eventDate, userId, eventDate)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) DeleteEvent(userId int, eventDate string) error {
	return nil
}

func (s *Storage) DayEvent(userId int, eventDate string) error {
	return nil
}

func (s *Storage) WeekEvent(userId int, eventDate string) error {
	return nil
}

func (s *Storage) MonthEvent(userId int, eventDate string) error {
	return nil
}
