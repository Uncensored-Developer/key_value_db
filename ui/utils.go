package ui

import (
	"bufio"
	"fmt"
	"kvdb/domain"
	"log"
)

// PrintDbResult prints the given database result(s) to the provided writer.
//
// The function takes a writer (*bufio.Writer) and a result (any) as parameters.
// The result can be either a slice of DBResult objects ([]domain.DBResult) or a single DBResult object (domain.DBResult).
// It writes the result(s) to the writer in a formatted manner.
// If the result is a slice, it iterates over each DBResult object and writes its SimpleMsg() value to the writer.
// If the result is a single object, it writes the SimpleMsg() value of that object to the writer.
// The function returns nothing.
func PrintDbResult(writer *bufio.Writer, result any) {
	if dbResults, ok := result.([]domain.DBResult); ok {
		for i, dbResult := range dbResults {
			_, err := fmt.Fprintf(writer, "%d) %v\n", i+1, dbResult.SimpleMsg())
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		dbResult := result.(domain.DBResult)
		_, err := fmt.Fprintf(writer, "%v\n", dbResult.SimpleMsg())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
