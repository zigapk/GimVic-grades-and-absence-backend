package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

var tableName string = "data"

//since 0.2 provided through sql flag
var sqlString string

var precision int = 1

func main() {

	sqlFlag := flag.String("sql", "root:root@/soc", "MySQL credidential in Go database format")
	portFlag := flag.Int("p", 8080, "Port to listen to")
	flag.Parse()

	sqlString = *sqlFlag

	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/years", yearsHandler)
	http.HandleFunc("/graph", graphHandler)
	http.ListenAndServe(":"+strconv.Itoa(*portFlag), nil)
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
	var tempExcusableStudents int = numberOf(where + optionalAnd + "excusable != 0")
	var tempInexcusableStudents int = numberOf(where + optionalAnd + "inexcusable != 0")
	var tempCurrentStudents int = numberOf(where)
	var tempExcusableStudentsPercent float64 = 0
	var tempInexcusableStudentsPercent float64 = 0
	if tempCurrentStudents != 0 {
		tempExcusableStudentsPercent = float64(tempExcusableStudents) / float64(tempCurrentStudents) * 100
		tempInexcusableStudentsPercent = float64(tempInexcusableStudents) / float64(tempCurrentStudents) * 100
	}

	response := &DataResponse{
		Facts: Facts{
			AllStudents:     numberOf(""),
			CurrentStudents: tempCurrentStudents,
			AllMale:         numberOf("gender = 'M'"),
			CurrentMale:     numberOf(where + optionalAnd + "gender = 'M'"),
			AllFemale:       numberOf("gender = 'Z'"),
			CurrentFemale:   numberOf(where + optionalAnd + "gender = 'Z'"),
		},
		Stats: Stats{
			AverageGrade:               RoundOn(average(gradeType, where), precision),
			ExcusableAbsence:           RoundOn(tempExcusableAbsence, precision),
			InexcusableAbsence:         RoundOn(tempInexcusableAbsence, precision),
			AverageAbsence:             RoundOn(tempExcusableAbsence+tempInexcusableAbsence, precision),
			ExcusableStudents:          tempExcusableStudents,
			InexcusableStudents:        tempInexcusableStudents,
			ExcusableStudentsPercent:   RoundOn(tempExcusableStudentsPercent, precision),
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
		if i != 0 {
			result += ","
		}
		result += year
	}
	fmt.Fprint(w, result)
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//data := getGraphData("excusable", "average", 0.1, 0.05, 51)

	queries := parseUrl(r)
	where := generateWhere(queries)

	var gradeType string = "average"
	if queries["gradeType"] != nil {
		gradeType = queries["gradeType"][0]
	}

	var absenceType string = "excusable"
	if queries["absenceType"] != nil {
		absenceType = queries["absenceType"][0]
	}

	var labels [3]string
	if gradeType == "final" {
		labels[0] = "Končni uspeh"
	} else {
		labels[0] = "Povprečna ocena"
	}
	if absenceType == "inexcusable" {
		labels[1] = "Št. neopravičenih ur v izbranem vzorcu"
		labels[2] = "Št. neopravičenih ur na celi šoli"
	} else {
		labels[1] = "Št. opravičenih ur v izbranem vzorcu"
		labels[2] = "Št. opravičenih ur na celi šoli"
	}

	result := getGraphJson(absenceType, gradeType, where, 1, 0.5, 6, labels)
	fmt.Fprint(w, result)
}

func getGraphJson(what, ranking, where string, step, diff float64, to int, labels [3]string) string {
	var optionalAnd string = ""
	if where != "" {
		optionalAnd = " and "
	}

	result := &GraphResponse{}
	result.Cols = append(result.Cols, Col{Type: "number", Label: labels[0]})
	result.Cols = append(result.Cols, Col{Type: "number", Label: labels[1]})
	result.Cols = append(result.Cols, Col{Type: "number", Label: labels[2]})

	for i := 0; i < to; i++ {
		lower := float64(i)*step - diff
		upper := float64(i)*step + diff
		avgCurrent := RoundOn(average(what, ranking+" >= "+strconv.FormatFloat(lower, 'f', -1, 64)+" and "+ranking+" < "+strconv.FormatFloat(upper, 'f', -1, 64)+optionalAnd+where), 2)
		avgAll := RoundOn(average(what, ranking+" >= "+strconv.FormatFloat(lower, 'f', -1, 64)+" and "+ranking+" < "+strconv.FormatFloat(upper, 'f', -1, 64)), 2)

		base := float64(float64(i) * step)
		var finalCurrent interface{} = nil
		if avgCurrent != 0 {
			finalCurrent = avgCurrent
		}
		var finalAll interface{} = nil
		if avgAll != 0 {
			finalAll = avgAll
		}

		result.Rows = append(result.Rows, Row{VStructs: []VStruct{VStruct{Value: base}, VStruct{Value: finalCurrent}, VStruct{Value: finalAll}}})
	}

	responseStr, err := json.Marshal(result)
	check(err)
	return string(responseStr)
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
	if queries["classA"] != nil && queries["classA"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'A'"
	}
	if queries["classB"] != nil && queries["classB"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'B'"
	}
	if queries["classC"] != nil && queries["classC"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'C'"
	}
	if queries["classD"] != nil && queries["classD"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'D'"
	}
	if queries["classE"] != nil && queries["classE"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'E'"
	}
	if queries["classF"] != nil && queries["classF"][0] == "false" {
		if where != "" {
			where += " and "
		}
		where += "class != 'F'"
	}

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

func differentYears() []string {
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

	//debug
	//fmt.Println("EXECUTING: " + query)

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

func RoundOn(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
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
	AverageGrade               float64
	AverageAbsence             float64
	ExcusableAbsence           float64
	InexcusableAbsence         float64
	ExcusableStudents          int
	InexcusableStudents        int
	ExcusableStudentsPercent   float64
	InexcusableStudentsPercent float64
}

type YearsResponse struct {
	Years []string
}

type GraphResponse struct {
	Cols []Col "json:\"cols\""
	Rows []Row "json:\"rows\""
}

type Col struct {
	Id    string "json:\"id\""
	Label string "json:\"label\""
	Type  string "json:\"type\""
}

type Row struct {
	VStructs []VStruct "json:\"c\""
}

type VStruct struct {
	Value interface{} "json:\"v\""
}
