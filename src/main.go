package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
)

//TODO: translate sql

var sqlString string = "root:root@/soc"

func handler(w http.ResponseWriter, r *http.Request) {
	queries := parseUrl(r)
	where := generateWhere(queries)

	var gradeType string = "average"
	if queries["gradeType"] != nil {
		gradeType = queries["gradeType"][0]
	}

	var optionalAnd string = ""
	if where != "" {
		optionalAnd = " and "
	}

	var tempExcusableAbsence float64 = average("opravicene", where)
	var tempInexcusableAbsence float64 = average("neopravicene", where)

	response := &Response{
		Facts: Facts{
			AllStudents:     numberOf("*", ""),
			CurrentStudents: numberOf("*", where),
			AllMale:         numberOf("*", "gender = 'M'"),
			CurrentMale:     numberOf("*", where+optionalAnd+"gender = 'M'"),
			AllFemale:       numberOf("*", "gender = 'Z'"),
			CurrentFemale:   numberOf("*", where+optionalAnd+"gender = 'Z'"),
		},
		Stats: Stats{
			AverageGrade:        average(gradeType, where),
			ExcusableAbsence:    tempExcusableAbsence,
			InexcusableAbsence:  tempInexcusableAbsence,
			AverageAbsence:      tempExcusableAbsence + tempInexcusableAbsence,
			ExcusableStudents:   numberOf("*", where+optionalAnd+"opravicene != 0"),
			InexcusableStudents: numberOf("*", where+optionalAnd+"neopravicene != 0"),
		},
	}
	responseStr, err := json.Marshal(response)
	check(err)
	fmt.Fprint(w, string(responseStr))
}

func generateWhere(queries map[string][]string) string {
	//years not included yet
	var where string = ""

	//for grades (1., 2., 3., 4. grade) - default is true
	if queries["grade1"] != nil && queries["grade1"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 1"
	}
	if queries["grade2"] != nil && queries["grade2"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 2"
	}
	if queries["grade3"] != nil && queries["grade3"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 3"
	}
	if queries["grade4"] != nil && queries["grade4"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 4"
	}

	//for classes (A, B, C, D, E, F) - default is true
	/* SQL data not ready yet!
	if queries["classA"] != nil && queries["classA"][0] == "false"{
		if where != "" {where += " and "}
		where += "class != 1"
	}*/

	//for gender (male, female) - default is true
	if queries["male"] != nil && queries["male"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "gender != 'M'"
	}
	if queries["female"] != nil && queries["female"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "gender != 'Z'"
	}

	return where
}

func parseUrl(r *http.Request) map[string][]string {
	str := r.URL.String()
	u, err := url.Parse(str)
	check(err)
	m, err := url.ParseQuery(u.RawQuery)
	check(err)
	return m
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func numberOf(what string, where string) int {
	var query string
	if where == "" {
		query = "select " + what + " from soc;"
	} else {
		query = "select " + what + " from soc where " + where + ";"
	}

	//debug
	fmt.Println("EXECUTING: " + query)

	con, err := sql.Open("mysql", sqlString)
	check(err)
	defer con.Close()

	rows, err := con.Query(query)
	check(err)
	var i int = 0

	for rows.Next() {
		i++
	}
	return i
}

func average(what, where string) float64 {
	var query string
	if where == "" {
		query = "select " + what + " from soc;"
	} else {
		query = "select " + what + " from soc where " + where + ";"
	}

	con, err := sql.Open("mysql", sqlString)
	check(err)
	defer con.Close()

	rows, err := con.Query(query)
	check(err)

	var sum float64 = 0
	var i int = 0

	for rows.Next() {
		i++
		var temp float64
		rows.Scan(&temp)
		sum += temp
	}

	if i == 0 {
		return -1
	}

	return float64(sum) / float64(i)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Facts Facts
	Stats Stats
}

type Facts struct {
	AllStudents     int
	CurrentStudents int
	AllMale         int
	CurrentMale     int
	AllFemale       int
	CurrentFemale   int
}

type Stats struct {
	AverageGrade        float64
	AverageAbsence      float64
	ExcusableAbsence    float64
	InexcusableAbsence  float64
	ExcusableStudents   int
	InexcusableStudents int
}
