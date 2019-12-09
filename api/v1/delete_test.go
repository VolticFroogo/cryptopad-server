package v1

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
	"github.com/VolticFroogo/cryptopad-server/api/v1/pad"
)

// deletePad tests a valid request for deleting a pad.
// Note: I can't use the name delete as that's a default Go func.
func deletePad(t *testing.T, client *http.Client) {
	new := model.Pad{
		ID:       "delete-pad",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Insert(new)
	if err != nil {
		t.Error(err.Error())
	}

	body := model.Pad{
		ID:    "delete-pad",
		Proof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("delete pad: could not put delete pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(new.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err != sql.ErrNoRows {
		t.Errorf("delete pad: pad still exists (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad: success (%v, %v)", res.Status, errorResponse.Error)
}

func deletePadIncorretProof(t *testing.T, client *http.Client) {
	new := model.Pad{
		ID:       "del-inc-proof",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(new.ID)
	if err != nil {
		t.Error(err.Error())
	}

	err = pad.Insert(new)
	if err != nil {
		t.Error(err.Error())
	}

	body := model.Pad{
		ID:    "del-inc-proof",
		Proof: "OTHER-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusForbidden {
		t.Errorf("delete pad incorrect proof: expected status forbidden (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(new.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err == sql.ErrNoRows {
		t.Errorf("delete pad incorrect proof: pad doesn't exist (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad incorrect proof: success (%v, %v)", res.Status, errorResponse.Error)
}

func deletePadNonExistant(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID:    "non-existant",
		Proof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("delete pad non existant: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad non existant: success (%v, %v)", res.Status, errorResponse.Error)
}

func deletePadInvalidProofLen(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID:    "del-proof-len",
		Proof: "INVALID-LEN-PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("delete pad invalid proof len: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad invalid proof len: success (%v, %v)", res.Status, errorResponse.Error)
}

func deletePadIDTooShort(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID:    "abc",
		Proof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("delete pad id too short: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad id too short: success (%v, %v)", res.Status, errorResponse.Error)
}

func deletePadIDTooLong(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID:    "abcdefghijklmnopq",
		Proof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodDelete, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("delete pad id too long: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("delete pad id too long: success (%v, %v)", res.Status, errorResponse.Error)
}
