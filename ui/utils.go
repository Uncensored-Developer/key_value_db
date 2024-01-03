package ui

import (
	"bufio"
	"errors"
	"fmt"
	"kvdb/domain"
	"log"
	"strings"
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
	switch res := result.(type) {
	case []domain.DBResult:
		for i, dbResult := range res {
			_, err := fmt.Fprintf(writer, "%d) %v\n", i+1, dbResult.SimpleMsg())
			if err != nil {
				log.Fatalln(err)
			}
		}
	case domain.DBResult:
		_, err := fmt.Fprintf(writer, "%v\n", res.SimpleMsg())
		if err != nil {
			log.Fatalln(err)
		}
	default:
		_, err := fmt.Fprintf(writer, "%v\n", res)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// getCommand parses the input string and returns a domain.Command and an error.
//
// It takes an input string and trims any leading or trailing whitespace.
// The input string is then split into individual words using space as the delimiter.
// The function iterates over the words, checking for opening and closing quotes to handle quoted arguments correctly.
// If a closing quote is missing, an error is returned.
// The keyword variable is set to the first word in uppercase.
// The remaining words are appended to the args slice, either as individual arguments or as a single argument if enclosed in quotes.
// Finally, a domain.Command is created with the keyword and args, and returned along with any error encountered.
func getCommand(input string) (domain.Command, error) {
	var keyword string
	var args []any
	input = strings.TrimSpace(input)

	// Split the input into there individual words
	words := strings.Split(input, " ")
	foundOpeningQuote := false
	foundClosingQuote := false
	var w string
	for i, word := range words {
		if i != 0 {
			if word != "" {

				if foundOpeningQuote {
					if strings.HasSuffix(word, "\"") {
						foundClosingQuote = true
						w += word[:len(word)-1]
					} else {
						w += word + " "
					}
				} else {
					if strings.HasPrefix(word, "\"") {
						foundOpeningQuote = true
						w += word[1:] + " "
					} else {
						args = append(args, strings.TrimSpace(word))
					}
				}
				if foundOpeningQuote && foundClosingQuote {
					args = append(args, strings.TrimSpace(w))
					w = ""
					foundOpeningQuote = false
					foundClosingQuote = false
				}

			} else {
				if foundOpeningQuote && !foundClosingQuote {
					w += " "
				}
			}
		} else {
			keyword = strings.ToUpper(word)
		}
	}

	if foundOpeningQuote && !foundClosingQuote {
		return domain.Command{}, errors.New("(error) ERR Syntax error: arguments has no closing quote")
	}

	if len(args) > 2 {
		return domain.Command{}, errors.New("(error) ERR Syntax error")
	}
	return domain.NewCommand(keyword, args...), nil
}
