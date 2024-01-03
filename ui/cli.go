package ui

import (
	"bufio"
	"fmt"
	"kvdb/domain"
	"os"
)

func RunCLI(db domain.KeyValueDB) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command input:", err)
			return
		}

		cmd, err := getCommand(input)
		if err != nil {
			fmt.Println(err)
		} else {
			result := db.Execute(cmd)
			writer := bufio.NewWriter(os.Stdout)
			PrintDbResult(writer, result)
		}
	}
}
