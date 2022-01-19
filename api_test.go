package main

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	var err error
	db, err = sql.Open("mysql", "surajit:Spore@0020@tcp(127.0.0.1:3306)/customers")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)

	}

	req := httptest.NewRequest(http.MethodGet, "/customer?id=2", nil)
	w := httptest.NewRecorder()

	get(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	tc := `{"ID":2,"Name":"sam","Phone":979948379,"Address":"test2"}`

	if string(data) != tc {
		t.Errorf("expected %v got %v", tc, string(data))
	}
}

func TestHandler(t *testing.T) {
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
		id     string
		// output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"GET", nil, "2", http.StatusOK, []byte(`{"ID":2,"Name":"sam","Phone":979948379,"Address":"test2"}`)},
		//{"GET", nil, "100000", http.StatusNotFound, nil},
		//{"POST", []byte(`{"ID":11,"Name":"sam","Phone":979948379,"Address":"test2"}`), http.StatusCreated, []byte(`success`)},
	}

	for _, tc := range testcases {
		req := httptest.NewRequest(tc.method, "/customer?id=2", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", tc.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(tc.expectedResponse)
		if !reflect.DeepEqual(w.Body, expected) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}

	}
}

func TestDeleteById(t *testing.T) {
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
		id     int
		// output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		//{"DELETE", nil, 3, http.StatusOK, nil},
		{"DELETE", nil, 3, http.StatusOK, nil},
		//{"DELETE", nil, 3, http.StatusInternalServerError, nil},
	}

	for _, tc := range testcases {
		query := "/customer?id=" + strconv.Itoa(tc.id)
		req := httptest.NewRequest(tc.method, query, bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", tc.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(tc.expectedResponse)
		if !reflect.DeepEqual(w.Body, expected) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}

	}
}

func TestPost(t *testing.T) {
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
		//{"DELETE", nil, 3, http.StatusOK, nil},
		{"POST", []byte(`{"ID":3,"Name":"sam","Phone":979948379,"Address":"test2"}`), http.StatusCreated, []byte(`success`)},
		//{"POST", []byte(`{"ID":,"Name":"sam","Phone":979948379,"Address":"test2"}`), http.StatusBadRequest, nil},
		//{"POST", []byte(`{"ID":21,"Name":"sam","Phone":979948379,"Address":"test2"}`), http.StatusInternalServerError, nil},
	}

	for _, tc := range testcases {
		req := httptest.NewRequest(tc.method, "/customer", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", tc.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(tc.expectedResponse)
		if !reflect.DeepEqual(w.Body, expected) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}

	}
}

func TestUpdateById(t *testing.T) {
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
		id     string
		// output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		//{"DELETE", nil, 3, http.StatusOK, nil},
		{"PUT", []byte(`{"Name":"sam","Phone":979948379,"Address":"test2"}`), "", http.StatusInternalServerError, nil},
		{"PUT", []byte(`{"Name":"sam","Phone":979948379,"Address":"test2"}`), "2", http.StatusCreated, nil},
	}

	for _, tc := range testcases {
		query := "/customer?id=" + tc.id
		req := httptest.NewRequest(tc.method, query, bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("Expected %v\tGot %v", tc.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(tc.expectedResponse)
		if !reflect.DeepEqual(w.Body, expected) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}

	}
}
