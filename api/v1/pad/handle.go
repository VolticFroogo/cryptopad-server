package pad

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
	"github.com/VolticFroogo/cryptopad-server/helper"
	"github.com/gorilla/mux"
)

var (
	errIncorrectProof     = errors.New("proofs do not match")
	errInvalidIDLen       = fmt.Errorf("ids must be between %v and %v in length", model.IDLen.Min, model.IDLen.Max)
	errInvalidContentLen  = fmt.Errorf("ids must be between %v and %v in length", model.ContentLen.Min, model.ContentLen.Max)
	errInvalidProofLen    = fmt.Errorf("proofs must be 0 or %v in length", model.ProofLen)
	errInvalidNewProofLen = fmt.Errorf("new proofs must be 0 or %v in length", model.ProofLen)
	errNoProof            = errors.New("proof can not be empty")
	errNoNewProof         = errors.New("new proof can not be empty")
)

// Get a pad.
func Get(w http.ResponseWriter, r *http.Request) {
	// Get data from the request.
	var data model.Pad

	vars := mux.Vars(r)
	data.ID = vars["id"]

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

	pad.Proof = ""
	pad.NewProof = ""

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
	if pLen != 0 && pLen != model.ProofLen {
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
		// Check again that the new proof is the right length.
		// This is necessary as 0 could pass earlier, but shouldn't now.
		if npLen != model.ProofLen {
			helper.ThrowErr(errNoNewProof, http.StatusBadRequest, w)
			return
		}

		// Just insert a new pad into the database.
		err = Insert(data)

		if err != nil {
			helper.ThrowErr(err, http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else { // If the pad does exist:
		// Update the pad if the proof matches.
		err = updateIfTrusted(pad, data)

		if err != nil {
			if err == errIncorrectProof {
				helper.ThrowErr(err, http.StatusForbidden, w)
				return
			}

			helper.ThrowErr(err, http.StatusInternalServerError, w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
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
	err = Update(data)
	return
}

// Delete a pad.
func Delete(w http.ResponseWriter, r *http.Request) {
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
		helper.ThrowErr(errNoProof, http.StatusBadRequest, w)
		return
	}

	// Get the pad (if it exists) from the database with a matching ID.
	pad, err := FromID(data.ID)
	if err != nil && err != sql.ErrNoRows {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	if err == sql.ErrNoRows {
		helper.ThrowErr(err, http.StatusNotFound, w)
		return
	}

	if data.Proof != pad.Proof {
		helper.ThrowErr(errIncorrectProof, http.StatusForbidden, w)
		return
	}

	err = Remove(data.ID)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
