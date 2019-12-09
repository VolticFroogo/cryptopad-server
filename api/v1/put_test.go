package v1

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
	"github.com/VolticFroogo/cryptopad-server/api/v1/pad"
)

// putNew tests a valid request for putting a new pad into database.
func putNew(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID:       "new-pad",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	expected := model.Pad{
		ID:      "new-pad",
		Content: "ENCRYPTED-STUFF-HERE",
		Proof:   "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(body.ID)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("put new: could not put new pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	output, err := pad.FromID(body.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected != output {
		t.Errorf("put new: output differs from expected while putting new pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put new: success (%v, %v)", res.Status, errorResponse.Error)
}

// putUpdate tests a valid request for updating a pad in the database.
func putUpdate(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "update-pad",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	updated := model.Pad{
		ID:       "update-pad",
		Content:  "OTHER-ENCRYPTED-STUFF-HERE",
		Proof:    "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
		NewProof: "OTHER-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	expected := model.Pad{
		ID:      "update-pad",
		Content: "OTHER-ENCRYPTED-STUFF-HERE",
		Proof:   "OTHER-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	err = pad.Insert(original)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, updated, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("put update: could not put new pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	output, err := pad.FromID(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected != output {
		t.Errorf("put update: output differs from expected while putting updated pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put update: success (%v, %v)", res.Status, errorResponse.Error)
}

// putIncorrectProof tests a put request with a proof that doesn't match.
func putIncorrectProof(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "incorrect-proof",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	updated := model.Pad{
		ID:      "incorrect-proof",
		Content: "OTHER-ENCRYPTED-STUFF-HERE",
		Proof:   "OTHER-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	expected := model.Pad{
		ID:      "incorrect-proof",
		Content: "ENCRYPTED-STUFF-HERE",
		Proof:   "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	err = pad.Insert(original)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, updated, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusForbidden {
		t.Errorf("put incorrect proof: expected status forbidden (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	output, err := pad.FromID(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected != output {
		t.Errorf("put incorrect proof: pad still got updated with invalid incorrect proof (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put incorrect proof: success (%v, %v)", res.Status, errorResponse.Error)
}

// putInvalidProofLen tests a put request with a proof with an invalid length.
func putInvalidProofLen(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "proof-len",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	updated := model.Pad{
		ID:      "proof-len",
		Content: "OTHER-ENCRYPTED-STUFF-HERE",
		Proof:   "INVALID-LEN-PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	expected := model.Pad{
		ID:      "proof-len",
		Content: "ENCRYPTED-STUFF-HERE",
		Proof:   "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	err = pad.Insert(original)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, updated, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("put invalid proof len: could not put new pad (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	output, err := pad.FromID(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected != output {
		t.Errorf("put invalid proof len: pad still got updated with invalid proof len (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put invalid proof len: success (%v, %v)", res.Status, errorResponse.Error)
}

// putInvalidProofLen tests a put request with a new proof with an invalid length.
func putInvalidNewProofLen(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "proof-len",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "INVALID-LEN-PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, original, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("put invalid new proof len: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(original.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err != sql.ErrNoRows {
		t.Errorf("put invalid new proof len: pad still got inserted with invalid proof len (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put invalid new proof len: success (%v, %v)", res.Status, errorResponse.Error)
}

func putNoNewProof(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "proof-len",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, original, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("put no new proof: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(original.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err != sql.ErrNoRows {
		t.Errorf("put no new proof: pad still got inserted with no new proof (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put no new proof: success (%v, %v)", res.Status, errorResponse.Error)
}

// putIDTooLong tests a put request with an ID which is too short.
func putIDTooShort(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "abc",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, original, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("put id too short: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(original.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err != sql.ErrNoRows {
		t.Errorf("put id too shortn: pad still got inserted with id too short (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put id too short: success (%v, %v)", res.Status, errorResponse.Error)
}

// putIDTooLong tests a put request with an ID which is too long.
func putIDTooLong(t *testing.T, client *http.Client) {
	original := model.Pad{
		ID:       "abcdefghijklmnopq",
		Content:  "ENCRYPTED-STUFF-HERE",
		NewProof: "INVALID-LEN-PROOF-KEY-ABCDEFGHIJKLMNOPQRSTUV",
	}

	err := pad.Remove(original.ID)
	if err != nil {
		t.Error(err.Error())
	}

	res, err, errorResponse := request(t, client, original, nil, http.MethodPut, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("put id too long: expected status bad request (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	_, err = pad.FromID(original.ID)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err.Error())
	}

	if err != sql.ErrNoRows {
		t.Errorf("put id too long: pad still got inserted with id too long (%v, %v)", res.Status, errorResponse.Error)
		return
	}

	t.Logf("put id too long: success (%v, %v)", res.Status, errorResponse.Error)
}
