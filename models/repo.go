package models

import (
	"database/sql"
	"encoding/json"

	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/gorp.v1"
)

// Database gin Middlware to select database
func Database(connString string) gin.HandlerFunc {
	dbmap := InitDb(connString)
	return func(c *gin.Context) {
		c.Set("DBmap", dbmap)
		c.Next()
	}
}

// Database gin Middlware to select database
func RedisPool(url string, password string, maxConnections int) gin.HandlerFunc {
	pool := InitRedisPool(url, password, maxConnections)
	return func(c *gin.Context) {
		c.Set("Pool", pool)
		c.Next()
	}
}

// InitDb set or create db
func InitDb(dbName string) *gorp.DbMap {

	db, err := sql.Open("mysql", dbName)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// XXX fix tables names
	dbmap.AddTableWithName(Agent{}, "agent").SetKeys(true, "Id")
	dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func InitRedisPool(url string, password string, maxConnections int) *redis.Pool {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", url)
		if err != nil {
			return nil, err
		}
		if password != "" {
			c.Do("AUTH", password)
		}
		return c, err
	}, maxConnections)
	return redisPool
}

// ParseQuery parse query to set select SQL query
func ParseQuery(q map[string][]string) string {
	query := " "
	if q["_filters"] != nil {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(q["_filters"][0]), &data)
		if err == nil {
			query = query + " WHERE "
			var searches []string
			for col, search := range data {
				valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
				if col != "" && search != "" && valid.MatchString(col) && valid.MatchString(search) {
					searches = append(searches, col+" LIKE \"%"+search+"%\"")
				}
			}
			query = query + strings.Join(searches, " AND ") // TODO join with OR for same keys
		}
	}
	if q["_sortField"] != nil && q["_sortDir"] != nil {
		sortField := q["_sortField"][0]
		// prevent SQLi
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if !valid.MatchString(sortField) {
			sortField = ""
		}
		if sortField == "created" || sortField == "updated" { // XXX trick for sqlite
			sortField = "datetime(" + sortField + ")"
		}
		sortOrder := q["_sortDir"][0]
		if sortOrder != "ASC" {
			sortOrder = "DESC"
		}
		if sortField != "" {
			query = query + " ORDER BY " + sortField + " " + sortOrder
		}
	}
	// _page, _perPage : LIMIT + OFFSET
	perPageInt := 0
	if q["_perPage"] != nil {
		perPage := q["_perPage"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		if valid.MatchString(perPage) {
			perPageInt, _ = strconv.Atoi(perPage)
			query = query + " LIMIT " + perPage
		}
	}
	if q["_page"] != nil {
		page := q["_page"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		pageInt, _ := strconv.Atoi(page)
		if valid.MatchString(page) && pageInt > 1 {
			offset := (pageInt-1)*perPageInt + 1
			query = query + " OFFSET " + strconv.Itoa(offset)
		}
	}
	return query
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
