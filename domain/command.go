package domain

import "fmt"

const (
	SET     string = "SET"
	GET     string = "GET"
	DEL     string = "DEL"
	INCR    string = "INCR"
	INCRBY  string = "INCRBY"
	MULTI   string = "MULTI"
	DISCARD string = "DISCARD"
	EXEC    string = "EXEC"
)

type CommandError struct {
	msg string
}

func (c *CommandError) Error() string {
	return fmt.Sprintf("(error) ERR %s", c.msg)
}

type Command struct {
	Keyword string
	Key     string
	Value   any
}

func NewCommand(keyword string, args ...any) Command {
	var key string
	var value any

	if len(args) > 0 {
		key = fmt.Sprintf("%v", args[0])
	}

	if len(args) > 1 {
		value = args[1]
	}

	return Command{
		keyword,
		key,
		value,
	}
}

// Validate checks if the command is valid and returns a boolean value and an error.
//
// It validates the command based on its keyword and checks if the required arguments are present.
// It returns a boolean value indicating whether the command is valid or not, and an error if any.
func (c Command) Validate() (bool, error) {
	var errMsg string
	var keyword string

	switch c.Keyword {
	case SET, INCRBY:
		keyword = SET
		if c.Keyword == INCRBY {
			keyword = INCRBY
		}

		if c.Key == "" {
			errMsg = fmt.Sprintf(
				"%s command expected 2 arguments but none was given (i.e no Key & value)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		if c.Value == nil {
			errMsg = fmt.Sprintf("%s command expected 2 arguments but 1 was given (i.e no value)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		return true, nil
	case GET, DEL, INCR:
		keyword = GET
		if c.Keyword == DEL {
			keyword = DEL
		} else if c.Keyword == INCR {
			keyword = INCR
		}

		if c.Key == "" {
			errMsg = fmt.Sprintf("%s command expected 1 argument but none was given (i.e no Key)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		if c.Value != nil {
			errMsg = fmt.Sprintf("%s command expected 1 argument but 2 was given", keyword)
			return false, &CommandError{msg: errMsg}
		}
		return true, nil
	case MULTI, DISCARD, EXEC:
		keyword = MULTI
		if c.Keyword == DISCARD {
			keyword = DISCARD
		} else if c.Keyword == EXEC {
			keyword = EXEC
		}
		if c.Key != "" {
			errMsg = fmt.Sprintf("%s command expected no argument but was given", keyword)
			return false, &CommandError{msg: errMsg}
		}
		if c.Value != nil {
			errMsg = fmt.Sprintf("%s command expected no argument but was given", keyword)
			return false, &CommandError{msg: errMsg}
		}
		return true, nil
	}
	return false, &CommandError{
		msg: fmt.Sprintf("unknown command %s", c.Keyword),
	}
}

func (c Command) String() string {
	return fmt.Sprintf("{Keyword: %q, Key: %q, Value: %v}", c.Keyword, c.Key, c.Value)
}
