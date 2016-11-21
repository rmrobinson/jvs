package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DeviceData struct {
	Id       string
	Name     string
	Address  string
	Type     string
	IsOn     bool
	IsActive bool
}

type BridgeData struct {
	Id   string
	Name string
}

type storage struct {
	db *sql.DB
}

func (s *storage) Open(fname string) (err error) {
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

func (s *storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *storage) setupDb() (err error) {
	cmd := `CREATE TABLE IF NOT EXISTS devices(
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		address TEXT,
		type TEXT,
		is_on BOOLEAN,
		is_active BOOLEAN
		);
		CREATE TABLE IF NOT EXISTS bridges(
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT
		);`

	_, err = s.db.Exec(cmd)
	return
}

func (s *storage) Devices() (devices []DeviceData, err error) {
	cmd := "SELECT id, name, address, type, is_on, is_active FROM devices;"

	rows, err := s.db.Query(cmd)

	if err != nil {
		return
	}

	devices = []DeviceData{}

	for rows.Next() {
		var d DeviceData
		err = rows.Scan(&d.Id, &d.Name, &d.Address, &d.Type, &d.IsOn, &d.IsActive)

		if err != nil {
			return
		}

		devices = append(devices, d)
	}

	return
}

func (s *storage) Device(Id string) (d DeviceData, err error) {
	cmd := `SELECT id, name, address, type, is_on, is_active FROM devices
		WHERE id=?;`

	err = s.db.QueryRow(cmd, Id).Scan(&d.Id, &d.Name, &d.Address, &d.Type, &d.IsOn, &d.IsActive)

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

func (s *storage) SetDevice(d DeviceData) (err error) {
	cmd := `INSERT OR REPLACE INTO devices(
		id,
		name,
		address,
		type,
		is_on,
		is_active
		) VALUES
		(?, ?, ?, ?, ?, ?);`

	stmt, err := s.db.Prepare(cmd)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(d.Id, d.Name, d.Address, d.Type, d.IsOn, d.IsActive)

	return
}

func (s *storage) Bridges() (bridges []BridgeData, err error) {
	cmd := "SELECT id, name FROM bridges;"

	rows, err := s.db.Query(cmd)

	if err != nil {
		return
	}

	bridges = []BridgeData{}

	for rows.Next() {
		var b BridgeData
		err = rows.Scan(&b.Id, &b.Name)

		if err != nil {
			return
		}

		bridges = append(bridges, b)
	}

	return
}

func (s *storage) Bridge(Id string) (b BridgeData, err error) {
	cmd := `SELECT id, name FROM bridges
		WHERE id=?;`

	err = s.db.QueryRow(cmd, Id).Scan(&b.Id, &b.Name)

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

func (s *storage) SetBridge(b BridgeData) (err error) {
	cmd := `INSERT OR REPLACE INTO bridges(
		id,
		name
		) VALUES
		(?, ?);`

	stmt, err := s.db.Prepare(cmd)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(b.Id, b.Name)

	return
}
