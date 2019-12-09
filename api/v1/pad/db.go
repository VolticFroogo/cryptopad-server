package pad

import (
	"database/sql"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
	"github.com/VolticFroogo/cryptopad-server/db"
)

// FromID gets a pad from the database given an ID.
func FromID(id string) (pad model.Pad, err error) {
	// Query a row from our ID.
	row, err := db.Dot.QueryRow(
		db.SQL,
		"v1-pad-from-id",
		id,
	)

	if err != nil {
		return
	}

	// Scan the row into our profile.
	err = scan(&pad, row)
	return
}

// Insert a pad into the database.
func Insert(pad model.Pad) (err error) {
	_, err = db.Dot.Exec(
		db.SQL,
		"v1-insert-pad",
		pad.ID,
		pad.Content,
		pad.NewProof,
	)

	return
}

// Update a pad in the database.
func Update(pad model.Pad) (err error) {
	_, err = db.Dot.Exec(
		db.SQL,
		"v1-update-pad",
		pad.Content,
		pad.NewProof,
		pad.ID,
	)

	return
}

// Remove a pad from the database.
func Remove(id string) (err error) {
	_, err = db.Dot.Exec(
		db.SQL,
		"v1-remove-pad",
		id,
	)

	return
}

func scan(pad *model.Pad, row *sql.Row) error {
	return row.Scan(
		&pad.ID,
		&pad.Content,
		&pad.Proof,
	)
}
