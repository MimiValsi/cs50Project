package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Info struct {
	ID       int // primary key
	Priority int
	SourceID int // FK references to source id (PK)
	Agent    string
	Material string
	Detail   string
	Estimate string
	Status   string

	ZeroTime time.Time
	Created  time.Time
	Updated  time.Time
}

// It sends data to DB
func (i *Info) Insert(id int, conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO info
    (source_id, agent, material, details, priority,
	estimate, status, created)
	  VALUES
	    ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
`
	err := conn.QueryRow(ctx, query, id, i.Agent,
		i.Material, i.Detail, i.Priority,
		i.Estimate, i.Status,
		time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return -1, err
	}

	return i.ID, nil
}

// Retrieve data from a choosen info
func (i *Info) InfoGet(id int, conn *pgxpool.Conn) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, agent, material, priority, details, estimate,
       source_id, created, updated, status
FROM info
  WHERE id = $1
`
	var estimate *string
	var updated *time.Time

	iObj := &Info{}
	err := conn.QueryRow(ctx, query, id).Scan(&iObj.ID, &iObj.Agent,
		&iObj.Material, &iObj.Priority, &iObj.Detail,
		&estimate, &iObj.SourceID,
		&iObj.Created, &updated, &iObj.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0001, time.January,
		1, 0, 0, 0, 0, time.UTC)

	if updated != nil {
		iObj.Updated = *updated
	}

	if estimate != nil {
		iObj.Estimate = *estimate
	}

	return iObj, nil
}

// Fetch a list with the minimum data from a info.
// It's specialy used within source view web page
func (i *Info) InfoList(id int, conn *pgxpool.Conn) ([]*Info, error) {
	ctx := context.Background()
	query := `
SELECT id,
       material,
       created,
       status,
       source_id,
       priority
FROM info
  WHERE source_id = $1
  ORDER BY priority ASC
`
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []*Info{}

	for rows.Next() {
		iObj := &Info{}

		err = rows.Scan(&iObj.ID, &iObj.Material,
			&iObj.Created, &iObj.Status,
			&iObj.SourceID, &iObj.Priority)
		if err != nil {
			return nil, err
		}

		infos = append(infos, iObj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return infos, nil
}

func (i *Info) InfoDelete(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
DELETE FROM info
  WHERE id = $1
`
	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// info update
func (i *Info) InfoUpdate(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE info
SET agent = $1, material = $2, priority = $3, details = $4,
	estimate = $5, updated = $6, status = $7
WHERE id = $8
`
	_, err := conn.Exec(ctx, query, i.Agent, i.Material,
		i.Priority, i.Detail, i.Estimate,
		time.Now().UTC(), i.Status, id)
	if err != nil {
		return err
	}

	return nil
}
