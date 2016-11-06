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

func TestAgent(t *testing.T) {
	defer deleteFile("_test.sqlite3")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Database("_test.sqlite3"))

	var urla = "/api/v1/agents"
	router.POST(urla, PostAgent)
	router.GET(urla, GetAgents)
	router.GET(urla+"/:id", GetAgent)
	router.DELETE(urla+"/:id", DeleteAgent)
	router.PUT(urla+"/:id", UpdateAgent)

	// Add
	log.Println("= http POST Agent")
	var a = Agent{Name: "Name test", IP: "Ip test"}
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

	// Add second agent
	log.Println("= http POST more Agent")
	var a2 = Agent{Name: "Name test2", IP: "Ip test2"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Agents")
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Agent
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
	log.Println("= http GET one Agent")
	var a1 Agent
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
	log.Println("= http DELETE one Agent")
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
	log.Println("= http PUT one Agent")
	//var a4 = Agent{Name: "Name test2 updated", IP: "Ip test2"}
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 Agent
	req, err = http.NewRequest("GET", urla+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Name, a3.Name, "a2 Name updated")

}
