package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Profile struct {
	Id             string
	Username       string
	LastModifiedAt time.Time
}

type Storage struct {
	db *sql.DB
}

func (s *Storage) Open(fname string) (err error) {
	db, err := sql.Open("sqlite3", fname)

	if err != nil {
		err = fmt.Errorf("Unable to open config db: %s\n", err)
		return
	}

	// Not sure why bit sql.Open doesn't like s.db
	s.db = db
	err = s.setupDb()
	return
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Storage) setupDb() (err error) {
	cmd := `CREATE TABLE IF NOT EXISTS profiles(
		id TEXT NOT NULL PRIMARY KEY,
		username TEXT,
		lastModifiedTime DATETIME
		);`

	_, err = s.db.Exec(cmd)
	return
}

func (s *Storage) Profile(Id string) (p Profile, err error) {
	cmd := `SELECT id, username, lastModifiedTime FROM profiles
		WHERE id=?;`

	err = s.db.QueryRow(cmd, Id).Scan(&p.Id, &p.Username, &p.LastModifiedAt)

	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("Id not present: %s", Id)
		return
	case err != nil:
		return
	default:
		return
	}
}

func (s *Storage) SetProfile(p Profile) (err error) {
	cmd := `INSERT OR REPLACE INTO profiles(
		id,
		username,
		lastModifiedTime
		) VALUES
		(?, ?, CURRENT_TIMESTAMP);`

	stmt, err := s.db.Prepare(cmd)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(p.Id, p.Username)

	return
}
