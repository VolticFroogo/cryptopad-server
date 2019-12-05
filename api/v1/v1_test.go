package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/VolticFroogo/cryptopad-server/db"
	"github.com/gorilla/mux"
)

const (
	location = "http://localhost"
	port     = ":8080"
	baseURL  = location + port + urlPrefix
	dbCfgDir = "../../configs/db_test.ini"
	timeout  = time.Second * 2
)

// ErrorResponse is the type used for error JSON responses.
type ErrorResponse struct {
	Error string
}

func TestV1(t *testing.T) {
	// Initialise the DB.
	err := db.Init(dbCfgDir)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// Start handling incoming requests.
	// Create a new Mux Router with strict slash.
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Handle v1 of the API.
	Handle(r)

	// Start serving on a seperate thread.
	go http.ListenAndServe(port, r)

	// Create an HTTP client to make requests with.
	client := &http.Client{
		Timeout: timeout,
	}

	// Run all get related tests.
	get(t, client)
	getNoID(t, client)
	getIDTooShort(t, client)
	getIDTooLong(t, client)
	getNonExistant(t, client)

	// Run all put related tests.
	putNew(t, client)
	putUpdate(t, client)
	putIncorrectProof(t, client)
	putInvalidProofLen(t, client)
	putInvalidNewProofLen(t, client)
	putIDTooShort(t, client)
	putIDTooLong(t, client)
}

func request(t *testing.T, client *http.Client, body interface{}, output interface{}, method, url string) (res *http.Response, err error, errorResponse ErrorResponse) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return
	}

	res, err = client.Do(req)
	if err != nil {
		return
	}

	outputBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(outputBody, &errorResponse)
	if err != nil {
		return
	}

	if output != nil {
		err = json.Unmarshal(outputBody, &output)
		if err != nil {
			return
		}
	}

	return
}
