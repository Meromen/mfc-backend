package db

import (
	"database/sql"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"io/ioutil"
)

const (
	createTableQuery string = `
		CREATE TABLE IF NOT EXISTS mfc (
			id                      TEXT,
			name                    TEXT,
			organization_full_name  TEXT,
			organization_address    TEXT,
			completed_tickets_count INTEGER,
			pending_tickets_count   INTEGER,
			lat                     FLOAT,
			lan                     FLOAT
		);	
	`
	dropTableQuery string = `
		DROP TABLE IF EXISTS mfc;
	`
)

type MfcStorage struct {
	conn      *sql.DB
	tableName string
}

func (s *MfcStorage) initializeTable() error {
	bdTx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	_, err = bdTx.Exec(dropTableQuery)
	if err != nil {
		bdTx.Rollback()
		return err
	}

	_, err = bdTx.Exec(createTableQuery)
	if err != nil {
		bdTx.Rollback()
		return err
	}

	body, err := ioutil.ReadFile("mfc-list.json")
	if err != nil {
		bdTx.Rollback()
		return err
	}

	mfcs := make([]Mfc, 0)
	err = json.Unmarshal(body, &mfcs)
	if err != nil {
		bdTx.Rollback()
		return err
	}

	sqlq, _, err := sq.Insert(s.tableName).
		Columns("id", "name", "organization_full_name",
		"organization_address", "lat", "lan").
		Values("", "", "", "", 0, 0).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	stmt, err := bdTx.Prepare(sqlq)
	if err != nil {
		bdTx.Rollback()
		return err
	}

	for _, mfc := range mfcs {
		_, err = stmt.Exec(mfc.Id, mfc.Name, mfc.OrganizationFullName, mfc.OrganizationAddress, mfc.Lat, mfc.Lan)
		if err != nil {
			bdTx.Rollback()
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		bdTx.Rollback()
		return err
	}

	err = bdTx.Commit()
	if err != nil {
		bdTx.Rollback()
		return err
	}
	return nil
}

func (s *MfcStorage) SelectAll() ([]DBRow, error) {
	sqlQ, _, err := sq.Select("id", "name", "organization_full_name",
		"organization_address", "completed_tickets_count",
		"pending_tickets_count", "lat", "lan").
		From(s.tableName).
		ToSql()

	if err != nil {
		return make([]DBRow, 0), err
	}

	mfcs := make([]DBRow, 0)
	rows, err := s.conn.Query(sqlQ)
	if err != nil {
		return make([]DBRow, 0), err
	}

	for rows.Next() {
		mfc := Mfc{}
		err := rows.Scan(&mfc.Id, &mfc.Name, &mfc.OrganizationFullName, &mfc.OrganizationAddress, &mfc.CompletedTicketsCount,
			&mfc.PendingTicketsCount, &mfc.Lat, &mfc.Lan)
		if err != nil {
			return make([]DBRow, 0), err
		}
		mfcs = append(mfcs, mfc)
	}

	return mfcs, nil
}

func (s *MfcStorage) UpdateAll(mfcs []DBRow) error {
	bdTx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	sqlq, _, err := sq.Update(s.tableName).
		Set("completed_tickets_count", 0).
		Set("pending_tickets_count", 0).
		Where("id = ?", "").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	stmt, err := bdTx.Prepare(sqlq)
	if err != nil {
		bdTx.Rollback()
		return err
	}

	for _, row := range mfcs {
		mfc := (row).(*Mfc)

		_, err = stmt.Exec(mfc.CompletedTicketsCount, mfc.PendingTicketsCount, mfc.Id)
		if err != nil {
			bdTx.Rollback()
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		bdTx.Rollback()
		return err
	}

	err = bdTx.Commit()
	if err != nil {
		bdTx.Rollback()
		return err
	}

	return nil
}

func NewMfcStorage(conn *sql.DB) (Storage, error) {
	s := MfcStorage{
		conn:      conn,
		tableName: "mfc",
	}

	err := s.initializeTable()
	if err != nil {
		return nil, err
	}

	return &s, nil
}
