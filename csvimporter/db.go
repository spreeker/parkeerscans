package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"

	//"github.com/cheggaaa/pb"
	"io"
	"log"
	"os"

	"github.com/lib/pq"
)

//SQLImport import structure
type SQLImport struct {
	txn  *sql.Tx
	stmt *sql.Stmt
}

//AddRow Add a single row of data to the database
func (i *SQLImport) AddRow(columns ...interface{}) error {
	_, err := i.stmt.Exec(columns...)
	return err
}

//Commit the import to database
func (i *SQLImport) Commit() error {

	_, err := i.stmt.Exec()
	if err != nil {
		return err
	}

	// Statement might already be closed
	// therefore ignore errors
	_ = i.stmt.Close()

	return i.txn.Commit()
}

//CleanTargetTable the table we are importing to
func cleanTable(db *sql.DB, target string) {

	sql := fmt.Sprintf("TRUNCATE TABLE %s;", target)

	if _, err := db.Exec(sql); err != nil {
		panic(err)
	}
}

//dropTable drop it!
func dropTable(db *sql.DB, target string) {

	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", target)

	if _, err := db.Exec(sql); err != nil {
		panic(err)
	}

}

// add vakken / wegdelen to scans

//NewImport setup a new import struct
func NewImport(db *sql.DB, schema string, tableName string, columns []string) (*SQLImport, error) {

	txn, err := db.Begin()

	if err != nil {
		return nil, err
	}

	stmt, err := txn.Prepare(pq.CopyInSchema(schema, tableName, columns...))
	if err != nil {
		return nil, err
	}

	return &SQLImport{txn, stmt}, nil
}

//setup a database connection
func dbConnect(connStr string) (*sql.DB, error) {
	//connStr := connectStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	err = db.Ping()

	return db, err
}

// NewProgressBar initializes new progress bar based on size of file
//func NewProgressBar(file *os.File) *pb.ProgressBar {
//	fi, err := file.Stat()
//
//	total := int64(0)
//	if err == nil {
//		total = fi.Size()
//	}
//
//	bar := pb.New64(total)
//	bar.SetUnits(pb.U_BYTES)
//	return bar
//}

//LoadSingleCSV import single csv file into  database with progresbar
func LoadSingleCSV(filename string, pgTable *SQLImport) {

	//var bar *pb.ProgressBar
	// load file
	csvfile, err := os.Open(filename)

	defer func() {
		err = csvfile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//bar = NewProgressBar(csvfile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	//reader := csv.NewReader(io.TeeReader(csvfile, bar))
	reader := csv.NewReader(csvfile)
	//bar.Start()
	//stream contents of csv into postgres
	importCSV(pgTable, reader, filename)

	if err != nil {
		log.Println(filename)
	}
	//bar.Finish()
	//update DateMap
	//DateMap[filename] = startEnd
	log.Printf("Batch: %s", filename)
}

//CreateTables  table to put csv data in
func CreateTables(db *sql.DB, csvfile string) (string, string) {

	validTableName := regexp.MustCompile("_\\d\\d\\d\\d([-_a-z_0-9]*)")
	//validTableName := regexp.MustCompile("_2017([-_a-z_0-9]*)")

	tableName := validTableName.FindString(csvfile)

	if tableName == "" {
		panic(errors.New("filename no regexmatch"))
	}

	tableName = tableName[1:]
	tableName = strings.Replace(tableName, "-", "_", -1)

	targetTable := fmt.Sprintf("scans_%s", tableName)
	importTable := fmt.Sprintf("import_%s", tableName)

	log.Println("Tablename", targetTable)

	// makeTable(db, targetTable, true)
	makeTable(db, importTable, false)
	// we have partitioned table now.
	targetTable = "metingen_scan"
	return importTable, targetTable

}

func makeTable(db *sql.DB, tableName string, inherit bool) {

	dropTable(db, tableName)

	sql := ""
	if inherit {
		sql = fmt.Sprintf(`CREATE UNLOGGED TABLE %s () INHERITS (metingen_scan);`, tableName)
	} else {
		sql = fmt.Sprintf(`CREATE UNLOGGED TABLE %s (LIKE metingen_scan INCLUDING DEFAULTS INCLUDING CONSTRAINTS);`, tableName)
	}

	if _, err := db.Exec(sql); err != nil {
		panic(err)
	}
}

func importCSV(pgTable *SQLImport, reader *csv.Reader, filename string) {

	reader.FieldsPerRecord = 0
	reader.Comma = ';'
	fileErrors := 0
	rows := 0

	// Determine mapping of csv columns to db columns

	for {
		record, err := reader.Read()
		rows += 1
		if rows == 1 {
			// skip header
			continue
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			csvError.Printf("%s: %s \n", err, record)
			failed++
			fileErrors++
			continue
		}

		cols, countedErrors, err := NormalizeRow(&record, rows)
		fileErrors += countedErrors

		if err != nil {
			log.Println(err)
			csvError.Printf("%s: %s \n", err, record)
			failed++
			fileErrors++
			continue
		}
		// printCols(cols, dbColumns)

		if err := pgTable.AddRow(cols...); err != nil {
			log.Println("DB ERROR")
			printRecord(&record, cols)
			printCols(cols, dbColumns)
			log.Println(filename)
			fileErrors++
			s := []string{filename, "FATAL"}
			fileErrorsMap[strings.Join(s, " ")] = 1
			panic(err)
		}

		success++
	}

	fileErrorsMap[filename] = fileErrors
}
