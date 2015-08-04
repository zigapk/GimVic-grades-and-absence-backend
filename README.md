GimVic - grades and absence
===========================
Backend for Gimnazija Vič's grades and absence statistics.

Note: a MySQL database is needed. It's format must be exacly the same as shown in backup.sql file. 

Usage:
-------
1. compile it yourself using `go build` command (you may need to install some dependencies using `go get`)
2. run it - by default it listens on port 8080 (can be modified in main.go file)
3. navigate to some of following urls:

<table>
  <tr>
    <td>URL</td>
    <td>Data</td>
    <td>Format</td>
  </tr>
  <tr>
    <td>/data</td>
    <td>Facts and plain statistics about data. Including statistics for all an the currently selected data.</td>
    <td>JSON</td>
  </tr>
  <tr>
    <td>/graph</td>
    <td>Line chart data in <a href="https://developers.google.com/chart/interactive/docs/reference#dataparam"> Google Chart JSON Format.</a></td>
    <td>JSON</td>
  </tr>
  <tr>
    <td>/years</td>
    <td>List of all avalible years.</td>
    <td>CSV</td>
  </tr>
</table>

Optional parameters can be declared as queries in URL:

<table>
  <tr>
    <td>Name</td>
    <td>Possible values</td>
    <td>Default value</td>
    <td>Aplies to /data</td>
    <td>Aplies to /graph</td>
  </tr>
  <tr>
    <td>`gradeType`</td>
    <td>`final` - final student's grade (1,2,3,4,5), `average` - average student's grade</td>
    <td>`average`</td>
    <td>X</td>
    <td>X</td>
  </tr>
  <tr>
    <td>`absenceType`</td>
    <td>`excusable`, `inexcusable`</td>
    <td>`excusable`</td>
    <td></td>
    <td>X</td>
  </tr>
  <tr>
    <td>`grade1`, `grade2`, `grade3`, `grade4`</td>
    <td>`true`, `false`</td>
    <td>`true`</td>
    <td>X</td>
    <td>X</td>
  </tr>
  <tr>
    <td>`classA`, `classB`, ... , `classF`</td>
    <td>`true`, `false`</td>
    <td>`true`</td>
    <td>X</td>
    <td>X</td>
  </tr>
  <tr>
    <td>`male`, `female`</td>
    <td>`true`, `false`</td>
    <td>`true`</td>
    <td>X</td>
    <td>X</td>
  </tr>
  <tr>
    <td>years: `2013-14`, `2014-15`, ... - aplies to years, provided by /years</td>
    <td>`true`, `false`</td>
    <td>`true`</td>
    <td>X</td>
    <td>X</td>
  </tr>
  
  
