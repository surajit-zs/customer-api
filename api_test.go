package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var db *sql.DB

func Test_Handler(t *testing.T) {

	var err error
	db, err = sql.Open("mysql", "surajit:Spore@0020@tcp(127.0.0.1:3306)/customers")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)

	}

	testcases := []struct {
		// input
		method string
		body   []byte
		// output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"GET", nil, http.StatusOK, []byte(`[{"Id":"1","Name":"ram","Phone":979948379,"Address":"test2"}]`)},
		{"POST", []byte(`{"Id":"12","Name":"ram","Phone":979948379,"Address":"test2"}`), http.StatusOK, []byte(`success`)},
		{"DELETE", nil, http.StatusMethodNotAllowed, nil},
	}

	for _, v := range testcases {
		req := httptest.NewRequest(v.method, "/customer", bytes.NewReader(v.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != v.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", v.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(v.expectedResponse)
		if !reflect.DeepEqual(w.Body, expected) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}
	}
}
