package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type column struct {
	ColumnName string
	Type       string
	Nullable   string
	Json       string
}

var db *sql.DB

//map for converting mysql type to golang types
var go_mysql_typemap = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "string",
	"datetime":           "string",
	"timestamp":          "string",
	"time":               "string",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

//function for generating golang struct
func generateModel(table_name string) {
	err, columns := getColumns(table_name)
	if err != nil {
		return
	}

	table_name = camelCase(table_name)
	depth := 1
	fmt.Print("type " + table_name + " struct {\n")
	for _, v := range columns {
		fmt.Print(tab(depth) + v.ColumnName + " " + v.Type + " " + v.Json)
		fmt.Print("\n")
	}
	fmt.Print(tab(depth-1) + "}\n")
}

// Function for fetching schema definition of passed table
func getColumns(table string) (errr error, columns []column) {
	rows, err := db.Query(`
		SELECT COLUMN_NAME,DATA_TYPE, IS_NULLABLE
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE() 
			AND TABLE_NAME = ? order by ORDINAL_POSITION`, table)
	if err != nil {
		fmt.Println("Error reading table information: ", err.Error())
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		col := column{}
		err := rows.Scan(&col.ColumnName, &col.Type, &col.Nullable)

		if err != nil {
			fmt.Println(err.Error())
			return err, nil
		}

		col.Json = strings.ToLower(col.ColumnName)
		col.ColumnName = camelCase(col.ColumnName)
		p := strings.Index(col.ColumnName, ", size")
		if p >= 0 {
			col.ColumnName = col.ColumnName[:p]
		}

		col.Type = go_mysql_typemap[col.Type]
		col.Json = fmt.Sprintf("`db:\"%s\" json:\"%s\"`", col.Json, strings.ToLower(col.ColumnName))

		columns = append(columns, col)
	}
	return err, columns
}

func camelCase(str string) string {
	name := strings.ToLower(str)
	var text string
	for _, p := range strings.Split(name, "_") {
		text += strings.ToUpper(p[0:1]) + p[1:]
	}
	return text
}

func tab(depth int) string {
	return strings.Repeat("\t", depth)
}

func main() {
	table_name := flag.String("t", "test_table", "table name")
	usr := flag.String("u", "root", "DB user name")
	pwd := flag.String("p", "", "DB password")
	db_name := flag.String("d", "test", "DB name")
	url := flag.String("url", "dburl", "dburl: localhost:3306")
	flag.Parse()
	conn := *usr + ":" + *pwd + "@tcp(" + *url + ")/" + *db_name
	var err error
	if conn != "" {
		db, err = sql.Open("mysql", conn)
		if err != nil {
			fmt.Println("[ERROR] Could not connect to database: ", err)
			os.Exit(1)
		}
	}
	generateModel(*table_name)
}

//Usage:   go run gorpgen.go -url=localhost:3306  -u=root -p=spwx -d=todolist -t=user
