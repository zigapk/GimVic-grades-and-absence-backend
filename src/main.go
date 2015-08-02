package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"math"
	"sort"
)

var tableName string = "data"
var sqlString string = "root:root@/soc"

var precision int = 1

func main() {
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/years", yearsHandler)
	http.ListenAndServe(":8080", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	//temp variables
	var tempExcusableAbsence float64 = average("excusable", where)
	var tempInexcusableAbsence float64 = average("inexcusable", where)
	var tempExcusableStudents int = numberOf(where+optionalAnd+"excusable != 0")
	var tempInexcusableStudents int = numberOf(where+optionalAnd+"inexcusable != 0")
	var tempCurrentStudents int = numberOf(where)
	var tempExcusableStudentsPercent float64 = 0
	var tempInexcusableStudentsPercent float64 = 0
	if tempCurrentStudents != 0 {
		tempExcusableStudentsPercent = float64(tempExcusableStudents)/float64(tempCurrentStudents)*100
		tempInexcusableStudentsPercent = float64(tempInexcusableStudents)/float64(tempCurrentStudents)*100
	}

	response := &DataResponse{
		Facts: Facts{
			AllStudents:     numberOf(""),
			CurrentStudents: tempCurrentStudents,
			AllMale:         numberOf("gender = 'M'"),
			CurrentMale:     numberOf(where+optionalAnd+"gender = 'M'"),
			AllFemale:       numberOf("gender = 'Z'"),
			CurrentFemale:   numberOf(where+optionalAnd+"gender = 'Z'"),
		},
		Stats: Stats{
			AverageGrade:        RoundOn(average(gradeType, where), precision),
			ExcusableAbsence:    RoundOn(tempExcusableAbsence, precision),
			InexcusableAbsence:  RoundOn(tempInexcusableAbsence, precision),
			AverageAbsence:      RoundOn(tempExcusableAbsence + tempInexcusableAbsence, precision),
			ExcusableStudents:   tempExcusableStudents,
			InexcusableStudents: tempInexcusableStudents,
			ExcusableStudentsPercent: RoundOn(tempExcusableStudentsPercent, precision),
			InexcusableStudentsPercent: RoundOn(tempInexcusableStudentsPercent, precision),
		},
	}
	responseStr, err := json.Marshal(response)
	check(err)
	fmt.Fprint(w, string(responseStr))
}

func yearsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	years := differentYears()
	var result string = ""
	for i, year := range years {
		if i!=0 {result += ","}
		result += year
	}
	fmt.Fprint(w, result)
}

func generateWhere(queries map[string][]string) string {
	var where string = ""

	//for years
	years := differentYears()
	for _, year := range years {
		if queries[year] != nil && queries[year][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "year != '" + year + "'"
	}
	}

	//for grades (1., 2., 3., 4. grade) - default is true
	if queries["grade1"] != nil && queries["grade1"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "grade != 1"
	}
	if queries["grade2"] != nil && queries["grade2"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "grade != 2"
	}
	if queries["grade3"] != nil && queries["grade3"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "grade != 3"
	}
	if queries["grade4"] != nil && queries["grade4"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "grade != 4"
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

func differentYears() []string{
	//count different years that appear in table
	var query string = "select distinct year from " + tableName + ";"

	con, err := sql.Open("mysql", sqlString)
	check(err)
	defer con.Close()

	rows, err := con.Query(query)
	check(err)

	var result []string

	for rows.Next() {
		var temp string
		rows.Scan(&temp)
		result = append(result, temp)
	}
	sort.Strings(result)
	return result
}
func parseUrl(r *http.Request) map[string][]string {
	str := r.URL.String()
	u, err := url.Parse(str)
	check(err)
	m, err := url.ParseQuery(u.RawQuery)
	check(err)
	return m
}

func numberOf(where string) int {
	var query string
	if where == "" {
		query = "select * from " + tableName + ";"
	} else {
		query = "select * from " + tableName + " where " + where + ";"
	}

	//debug
	//fmt.Println("EXECUTING: " + query)

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
		query = "select " + what + " from " + tableName + ";"
	} else {
		query = "select " + what + " from " + tableName + " where " + where + ";"
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
		return 0
	}

	return float64(sum) / float64(i)
}

func Round(f float64) float64 {
    return math.Floor(f + .5)
}

func RoundOn(f float64, places int) (float64) {
    shift := math.Pow(10, float64(places))
    return Round(f * shift) / shift;    
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type DataResponse struct {
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
	ExcusableStudentsPercent float64
	InexcusableStudentsPercent float64
}

type YearsResponse struct {
	Years []string
}