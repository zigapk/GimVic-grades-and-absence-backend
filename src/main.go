package main

import (
    "fmt"
    "net/http"
    "database/sql"
    "net/url"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
)

var sqlString string = "root:root@/soc"

func handler(w http.ResponseWriter, r *http.Request) {
    queries := parseUrl(r)
    
    var where string = ""
    var i int = 0
    for name, values := range queries {
    	if i != 0 {
    		where += " and "
    	}
    	where += name + "=" + values[0]
    	i++
    }
    var optionalAnd string = ""
    if where != "" {optionalAnd = " and "}

    response := &Response {
    	Facts: Facts {
    		AllStudents: numberOf("gender", ""),
    		CurrentStudents: numberOf("gender", where),
    		AllMale: numberOf("gender", "gender = 'M'"),
    		CurrentMale: numberOf("gender", where + optionalAnd + "gender = 'M'"),
    		AllFemale: numberOf("gender", "gender = 'Z'"),
    		CurrentFemale: numberOf("gender", where + optionalAnd + "gender = 'Z'"),
    	},
    }
    responseStr, err := json.Marshal(response)
    check(err)
    fmt.Fprint(w, string(responseStr))
}
 func parseUrl(r *http.Request) map[string][]string{
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
	AllStudents int
	CurrentStudents int
	AllMale int
	CurrentMale int
	AllFemale int
	CurrentFemale int

}

type Stats struct {
	AverageGrade float64
	AverageAbsence float64
	ExcusableAbsence float64
	InexcusableAbsence float64
	ExcusableStudents int
	InexcusableStudents int
}