package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUser(t *testing.T) {
	defer deleteFile("_test.sqlite3")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Database("_test.sqlite3"))

	var urla = "/api/v1/users"
	router.POST(urla, PostUser)
	router.GET(urla, GetUsers)
	router.GET(urla+"/:id", GetUser)
	router.DELETE(urla+"/:id", DeleteUser)
	router.PUT(urla+"/:id", UpdateUser)

	// Add
	log.Println("= http POST User")
	var a = User{Name: "Name test"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second user
	log.Println("= http POST more User")
	var a2 = User{Name: "Name test2"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Users")
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []User
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	log.Println("= Test parsing query")
	s := "http://127.0.0.1:8080/api?_filters={\"name\":\"t\"}&_sortDir=ASC&_sortField=created"
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query := ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, "  WHERE name LIKE \"%t%\"  ORDER BY datetime(created) ASC", query, "Parse query")

	// Get one
	log.Println("= http GET one User")
	var a1 User
	req, err = http.NewRequest("GET", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one User")
	req, err = http.NewRequest("DELETE", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all for count success")
	//fmt.Println(resp.Body)
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 1, len(as), "1 result")

	// Update one
	log.Println("= http PUT one User")
	//var a4 = User{Name: "Name test2 updated"}
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 User
	req, err = http.NewRequest("GET", urla+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Name, a3.Name, "a2 Name updated")

}
