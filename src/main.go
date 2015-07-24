package main

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var sqlString string = "root:root@/soc"

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Stevilo ucencev: ", stevilo("gender", ""), "\n")
	fmt.Fprint(w, "Stevilo moskih: ", stevilo("gender", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Stevilo zensk: ", stevilo("gender", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecni uspeh: ", povprecna("final", ""), "\n")
	fmt.Fprint(w, "Povprecni uspeh moskih: ", povprecna("final", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecni uspeh zensk: ", povprecna("final", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecna ocena: ", povprecna("average", ""), "\n")
	fmt.Fprint(w, "Povprecna ocena moskih: ", povprecna("average", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecna ocena zensk: ", povprecna("average", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur: ", povprecna("opravicene", ""), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur pri moskih: ", povprecna("opravicene", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur pri zenskah: ", povprecna("opravicene", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur v 1. letniku: ", povprecna("opravicene", "class = 1"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur v 2. letniku: ", povprecna("opravicene", "class = 2"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur v 3. letniku: ", povprecna("opravicene", "class = 3"), "\n")
	fmt.Fprint(w, "Povprecno stevilo opravicenih ur v 4. letniku: ", povprecna("opravicene", "class = 4"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur: ", povprecna("neopravicene", ""), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur pri moskih: ", povprecna("neopravicene", "gender = 'M'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur pri zenskah: ", povprecna("neopravicene", "gender = 'Z'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur v 1. letniku: ", povprecna("neopravicene", "class = '1'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur v 2. letniku: ", povprecna("neopravicene", "class = '2'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur v 3. letniku: ", povprecna("neopravicene", "class = '3'"), "\n")
	fmt.Fprint(w, "Povprecno stevilo neopravicenih ur v 4. letniku: ", povprecna("neopravicene", "class = '4'"), "\n")
	fmt.Fprint(w, "Stevilo dijakov z neopravicenimi urami: ", stevilo("neopravicene", "neopravicene != 0"), "\n")
	fmt.Fprint(w, "Stevilo moskih z neopravicenimi urami: ", stevilo("neopravicene", "neopravicene != 0 AND gender='M'"), "\n")
	fmt.Fprint(w, "Stevilo zensk z neopravicenimi urami: ", stevilo("neopravicene", "neopravicene != 0 AND gender='Z'"), "\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (5): ", povprecna("opravicene", "final = 5"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (4): ", povprecna("opravicene", "final = 4"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (3): ", povprecna("opravicene", "final = 3"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (2): ", povprecna("opravicene", "final = 2"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (1): ", povprecna("opravicene", "final = 1"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo opravicenih ur dijakov (0): ", povprecna("opravicene", "final = 0"), "\n")
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (5): ", povprecna("neopravicene", "final = 5"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (4): ", povprecna("neopravicene", "final = 4"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (3): ", povprecna("neopravicene", "final = 3"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (2): ", povprecna("neopravicene", "final = 2"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (1): ", povprecna("neopravicene", "final = 1"), "\n")
	fmt.Fprint(w, "Povprecno stevlilo neopravicenih ur dijakov (0): ", povprecna("neopravicene", "final = 0"), "\n")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

func stevilo(what string, where string) int {
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

func povprecna(what, where string) float64 {
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