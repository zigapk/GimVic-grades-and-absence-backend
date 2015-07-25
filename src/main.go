package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var sqlString string = "root:root@/soc"

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "numberOf ucencev: ", numberOf("gender", ""), "\n")
	fmt.Fprint(w, "numberOf moskih: ", numberOf("gender", "gender = 'M'"), "\n")
	fmt.Fprint(w, "numberOf zensk: ", numberOf("gender", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecni uspeh: ", average("final", ""), "\n")
	fmt.Fprint(w, "Povprecni uspeh moskih: ", average("final", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecni uspeh zensk: ", average("final", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "average ocena: ", average("average", ""), "\n")
	fmt.Fprint(w, "average ocena moskih: ", average("average", "gender = 'M'"), "\n")
	fmt.Fprint(w, "average ocena zensk: ", average("average", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur: ", average("opravicene", ""), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur pri moskih: ", average("opravicene", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur pri zenskah: ", average("opravicene", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur v 1. letniku: ", average("opravicene", "class = 1"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur v 2. letniku: ", average("opravicene", "class = 2"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur v 3. letniku: ", average("opravicene", "class = 3"), "\n")
	fmt.Fprint(w, "Povprecno numberOf opravicenih ur v 4. letniku: ", average("opravicene", "class = 4"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur: ", average("neopravicene", ""), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur pri moskih: ", average("neopravicene", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur pri zenskah: ", average("neopravicene", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur v 1. letniku: ", average("neopravicene", "class = '1'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur v 2. letniku: ", average("neopravicene", "class = '2'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur v 3. letniku: ", average("neopravicene", "class = '3'"), "\n")
	fmt.Fprint(w, "Povprecno numberOf neopravicenih ur v 4. letniku: ", average("neopravicene", "class = '4'"), "\n")
	fmt.Fprint(w, "numberOf dijakov z neopravicenimi urami: ", numberOf("neopravicene", "neopravicene != 0"), "\n")
	fmt.Fprint(w, "numberOf moskih z neopravicenimi urami: ", numberOf("neopravicene", "neopravicene != 0 AND gender='M'"), "\n")
	fmt.Fprint(w, "numberOf zensk z neopravicenimi urami: ", numberOf("neopravicene", "neopravicene != 0 AND gender='Z'"), "\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (5): ", average("opravicene", "final = 5"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (4): ", average("opravicene", "final = 4"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (3): ", average("opravicene", "final = 3"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (2): ", average("opravicene", "final = 2"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (1): ", average("opravicene", "final = 1"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (0): ", average("opravicene", "final = 0"), "\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (5): ", average("neopravicene", "final = 5"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (4): ", average("neopravicene", "final = 4"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (3): ", average("neopravicene", "final = 3"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (2): ", average("neopravicene", "final = 2"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (1): ", average("neopravicene", "final = 1"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (0): ", average("neopravicene", "final = 0"), "\n")
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