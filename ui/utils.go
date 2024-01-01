package ui

import (
	"bufio"
	"fmt"
	"kvdb/domain"
	"log"
)

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
