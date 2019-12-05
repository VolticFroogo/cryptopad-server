package v1

import (
	"net/http"
	"testing"
)

// putNew tests a valid request for putting a new pad into database.
func putNew(t *testing.T, client *http.Client) {

}

// putUpdate tests a valid request for updating a pad in the database.
func putUpdate(t *testing.T, client *http.Client) {

}

// putIncorrectProof tests a put request with a proof that doesn't match.
func putIncorrectProof(t *testing.T, client *http.Client) {

}

// putInvalidProofLen tests a put request with a proof with an invalid length.
func putInvalidProofLen(t *testing.T, client *http.Client) {

}

// putInvalidProofLen tests a put request with a new proof with an invalid length.
func putInvalidNewProofLen(t *testing.T, client *http.Client) {

}

// putIDTooLong tests a put request with an ID which is too short.
func putIDTooShort(t *testing.T, client *http.Client) {

}

// putIDTooLong tests a put request with an ID which is too long.
func putIDTooLong(t *testing.T, client *http.Client) {

}
