package device

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rmrobinson/jvs/service/device/pb"
)

type brDevice struct {
	ID       string
	Name     string
	IsOn     bool
	IsActive bool
}

type brBridge struct {
	ID   string
	Name string
}

type BottlerocketDB struct {
	db *sql.DB
}

func (db *BottlerocketDB) Open(fname string) (err error) {
	sqldb, err := sql.Open("sqlite3", fname)

	if err != nil {
		err = fmt.Errorf("Unable to open config db: %db\n", err)
		return
	}

	db.db = sqldb
	err = db.setupDb()
	return
}

func (db *BottlerocketDB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

func (db *BottlerocketDB) setupDb() error {
	setupCmd := `CREATE TABLE IF NOT EXISTS devices(
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		is_on BOOLEAN,
		is_active BOOLEAN
		);
		CREATE TABLE IF NOT EXISTS bridges(
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT
		);`

	_, err := db.db.Exec(setupCmd)

	if err != nil {
		return err
	}

	findBridgeCmd := `SELECT id, name FROM bridges;`

	b := brBridge{}
	err = db.db.QueryRow(findBridgeCmd).Scan(&b.ID, &b.Name)

	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == nil {
		return nil
	}

	seedDevicesCmd := `INSERT INTO devices(
		id,
		name,
		is_on,
		is_active
		) VALUES
		(?, ?, ?, ?);`

	seedDevicesStmt, err := db.db.Prepare(seedDevicesCmd)

	if err != nil {
		return err
	}

	defer seedDevicesStmt.Close()

	// Populate the devices
	houses := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
	const maxDeviceID = 16
	for _, houseID := range houses {
		for deviceID := 0; deviceID < maxDeviceID; deviceID++ {
			id := fmt.Sprintf("%s%d", houseID, deviceID)
			_, err = seedDevicesStmt.Exec(id, "X10 device", false, false)
		}
		if err != nil {
			return err
		}
	}

	// Populate the bridge
	seedBridgeCmd := `INSERT INTO bridges(
		id,
		name
		) VALUES
		(?, ?);`

	seedBridgeStmt, err := db.db.Prepare(seedBridgeCmd)

	if err != nil {
		return err
	}

	defer seedBridgeStmt.Close()

	_, err = seedBridgeStmt.Exec(uuid.New().String(), "X10 bridge")

	return nil
}

func (db *BottlerocketDB) Devices() ([]pb.Device, error) {
	cmd := "SELECT id, name, is_on, is_active FROM devices;"

	rows, err := db.db.Query(cmd)

	if err != nil {
		return nil, err
	}

	var brDevices []brDevice
	for rows.Next() {
		var d brDevice
		err = rows.Scan(&d.ID, &d.Name, &d.IsOn, &d.IsActive)

		if err != nil {
			return nil, err
		}

		brDevices = append(brDevices, d)
	}

	var devices []pb.Device
	for _, brd := range brDevices {
		d := pb.Device{
			Path:     brd.ID,
			IsActive: brd.IsActive,

			Config: &pb.DeviceConfig{
				Name: brd.Name,
			},
			State: &pb.DeviceState{
				Binary: &pb.DeviceState_BinaryState{
					IsOn: brd.IsOn,
				},
			},
		}

		devices = append(devices, d)
	}

	return devices, err
}

func (db *BottlerocketDB) SetDevice(d pb.Device) (err error) {
	cmd := `INSERT OR REPLACE INTO devices(
		id,
		name,
		is_on,
		is_active
		) VALUES
		(?, ?, ?, ?);`

	stmt, err := db.db.Prepare(cmd)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(d.Id, d.Config.Name, d.State.Binary.IsOn, d.IsActive)

	return
}

func (db *BottlerocketDB) ID() (string, error) {
	b, err := db.bridge()

	return b.ID, err
}

func (db *BottlerocketDB) Name() (string, error) {
	b, err := db.bridge()

	return b.Name, err
}

func (db *BottlerocketDB) SetName(name string) error {
	cmd := `UPDATE bridges SET name=?;`

	stmt, err := db.db.Prepare(cmd)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name)

	return err
}

func (db *BottlerocketDB) bridge() (b brBridge, err error) {
	cmd := `SELECT id, name FROM bridges;`

	err = db.db.QueryRow(cmd).Scan(&b.ID, &b.Name)

	switch {
	case err == sql.ErrNoRows:
		err = errors.New("database not setup")
		return
	case err != nil:
		return
	default:
		return
	}
}
