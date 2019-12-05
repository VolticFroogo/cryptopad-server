package pad

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
	"github.com/VolticFroogo/cryptopad-server/helper"
)

var (
	errIncorrectProof     = errors.New("proofs do not match")
	errInvalidIDLen       = fmt.Errorf("ids must be between %v and %v in length", model.IDLen.Min, model.IDLen.Max)
	errInvalidContentLen  = fmt.Errorf("ids must be between %v and %v in length", model.ContentLen.Min, model.ContentLen.Max)
	errInvalidProofLen    = fmt.Errorf("proofs must be %v in length", model.ProofLen)
	errInvalidNewProofLen = fmt.Errorf("new proofs must be 0 or %v in length", model.ProofLen)
)

// Get a pad.
func Get(w http.ResponseWriter, r *http.Request) {
	// Get data from the JSON request.
	var data model.Pad
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Check if the ID is a valid length.
	if !model.IDLen.Check(data.ID) {
		helper.ThrowErr(errInvalidIDLen, http.StatusBadRequest, w)
		return
	}

	// Get the pad (if it exists) from the database with a matching ID.
	pad, err := FromID(data.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.ThrowErr(err, http.StatusNotFound, w)
			return
		}

		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Return the pad to the client.
	helper.JSONResponse(pad, http.StatusOK, w)
}

// Put (overwrite / create) a pad.
func Put(w http.ResponseWriter, r *http.Request) {
	// Get data from the JSON request.
	var data model.Pad
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Check if the ID is a valid length.
	if !model.IDLen.Check(data.ID) {
		helper.ThrowErr(errInvalidIDLen, http.StatusBadRequest, w)
		return
	}

	// Check if the proof is a valid length.
	pLen := len(data.Proof)
	if pLen != model.ProofLen {
		helper.ThrowErr(errInvalidProofLen, http.StatusBadRequest, w)
		return
	}

	// Check if the new proof is a valid length.
	npLen := len(data.NewProof)
	if npLen != 0 && npLen != model.ProofLen {
		helper.ThrowErr(errInvalidNewProofLen, http.StatusBadRequest, w)
		return
	}

	// Check if the content is a valid length.
	if !model.ContentLen.Check(data.Content) {
		helper.ThrowErr(errInvalidContentLen, http.StatusBadRequest, w)
		return
	}

	// Get the pad (if it exists) from the database with a matching ID.
	pad, err := FromID(data.ID)
	if err != nil && err != sql.ErrNoRows {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	if err == sql.ErrNoRows { // If the pad doesn't exist:
		// Just insert a new pad into the database.
		err = Insert(data)

		if err == nil {
			helper.ThrowErr(err, http.StatusInternalServerError, w)
			return
		}
	} else { // If the pad does exist:
		// Update the pad if the proof matches.
		err = updateIfTrusted(pad, data)

		if err == nil {
			if err == errIncorrectProof {
				helper.ThrowErr(err, http.StatusForbidden, w)
				return
			}

			helper.ThrowErr(err, http.StatusInternalServerError, w)
			return
		}
	}

	// Tell the client the pad has been created or created.
	w.WriteHeader(http.StatusCreated)
}

func updateIfTrusted(pad, data model.Pad) (err error) {
	// Check if the proofs match.
	if pad.Proof != data.Proof {
		err = errIncorrectProof
		return
	}

	// If the new proof is empty, set it to the current proof.
	if data.NewProof == "" {
		data.NewProof = data.Proof
	}

	// Update the pad in the database.
	err = Update(pad)
	return
}