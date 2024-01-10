package sqlite

import (
	"L2-task11/internal/domain"
	"L2-task11/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
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

	q := `
		CREATE TABLE IF NOT EXISTS event(
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    event_date DATE NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_uesr ON event(user_id); 
		` // мб нужно создать составной индекс из сразу из двух полей, общий хэш или что-то еще

	_, err = db.Exec(q)
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

	q = `INSERT INTO event(user_id, event_date) VALUES (?, ?)`
	_, err = s.db.Exec(q, userId, eventDate)
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

	// i store only date and user id what is the point if updating idk ¯\_(ツ)_/¯
	q = `UPDATE event SET user_id = ?, event_date = ? WHERE user_id = ? AND event_date = ?`
	_, err = s.db.Exec(q, userId, eventDate, userId, eventDate)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) DeleteEvent(userId int, eventDate string) error {
	const fn = "storage.sqlite.DeleteEvent"

	q := `DELETE FROM event WHERE user_id = ? AND event_date = ?;`
	result, err := s.db.Exec(q, userId, eventDate)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAffected == 0 {
		return storage.ErrEventNotExists
	}

	return nil
}

func (s *Storage) EventsForDay(userId int, eventDate string) (*domain.ReqParameters, error) {
	const fn = "storage.sqlite.EventsForDay"

	// Выполнение запроса SELECT
	q := `SELECT user_id, event_date FROM event WHERE user_id = ? AND event_date = ?`
	row := s.db.QueryRow(q, userId, eventDate)

	res := domain.ReqParameters{}

	err := row.Scan(&res.UserId, &res.Date)
	if errors.Is(sql.ErrNoRows, err) {
		return nil, storage.ErrEventNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	CorrectDate, err := time.Parse(time.RFC3339, res.Date)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	res.Date = CorrectDate.Format(time.DateOnly)

	return &res, nil
}

func (s *Storage) EventsForWeek(userId int, eventDate string) (*[]domain.ReqParameters, error) {
	const fn = "storage.sqlite.EventsForWeek"

	q := `SELECT user_id, event_date FROM event WHERE user_id = ? AND strftime('%W', event_date) = strftime('%W', ?);`

	rows, err := s.db.Query(q, userId, eventDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	events := make([]domain.ReqParameters, 0)
	for rows.Next() {
		evnt := domain.ReqParameters{}
		err = rows.Scan(&evnt.UserId, &evnt.Date)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		CorrectDate, err := time.Parse(time.RFC3339, evnt.Date)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		evnt.Date = CorrectDate.Format(time.DateOnly)

		events = append(events, evnt)
	}

	return &events, nil
}

func (s *Storage) EventsForMonth(userId int, eventDate string) (*[]domain.ReqParameters, error) {
	const fn = "storage.sqlite.EventsForWeek"

	q := `SELECT user_id, event_date FROM event WHERE user_id = ? AND strftime('%m', event_date) = strftime('%m', ?);`

	rows, err := s.db.Query(q, userId, eventDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	events := make([]domain.ReqParameters, 0)
	for rows.Next() {
		evnt := domain.ReqParameters{}
		err = rows.Scan(&evnt.UserId, &evnt.Date)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		CorrectDate, err := time.Parse(time.RFC3339, evnt.Date)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		evnt.Date = CorrectDate.Format(time.DateOnly)

		events = append(events, evnt)
	}

	return &events, nil
}
