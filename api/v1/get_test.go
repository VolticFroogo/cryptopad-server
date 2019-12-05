package v1

import (
	"net/http"
	"testing"

	"github.com/VolticFroogo/cryptopad-server/api/v1/model"
)

// get tests a valid get request.
func get(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID: "test",
	}

	expected := model.Pad{
		ID:      "test",
		Content: "ENCRYPTED-STUFF-HERE",
		Proof:   "PROOF-KEY",
	}

	output := &model.Pad{}

	res, err, errorResponse := request(t, client, body, output, http.MethodGet, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("get pad: could not get valid pad (%v, %v)", res.Status, errorResponse.Error)
	} else if expected != *output {
		t.Errorf("get pad: output differs from expected while getting pad (%v, %v)", res.Status, errorResponse.Error)
	} else {
		t.Logf("get pad: success (%v, %v)", res.Status, errorResponse.Error)
	}
}

// getNoID tests a get request without providing an ID.
func getNoID(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID: "",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodGet, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("get pad no id: expected status bad request, got %v, %v", res.Status, errorResponse.Error)
	} else {
		t.Logf("get pad no id: success (%v, %v)", res.Status, errorResponse.Error)
	}
}

// getIDTooShort tests a get request with an ID which is too short.
func getIDTooShort(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID: "abc",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodGet, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("get pad id too short: expected status bad request, got %v, %v", res.Status, errorResponse.Error)
	} else {
		t.Logf("get pad id too short: success (%v, %v)", res.Status, errorResponse.Error)
	}
}

// getIDTooLong tests a get request with an ID which is too long.
func getIDTooLong(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID: "abcdefghijklmnopq",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodGet, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("get pad id too long: expected status bad request, got %v, %v", res.Status, errorResponse.Error)
	} else {
		t.Logf("get pad id too long: success (%v, %v)", res.Status, errorResponse.Error)
	}
}

// getNonExistant tests a valid request but the ID it's requesting doesn't exist.
func getNonExistant(t *testing.T, client *http.Client) {
	body := model.Pad{
		ID: "non-existant",
	}

	res, err, errorResponse := request(t, client, body, nil, http.MethodGet, baseURL+"pad")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("get pad non existant: success (%v, %v)", res.Status, errorResponse.Error)
	} else {
		t.Logf("get pad non existant: success (%v, %v)", res.Status, errorResponse.Error)
	}
}