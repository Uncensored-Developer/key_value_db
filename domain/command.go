package domain

import "fmt"

const (
	SET string = "SET"
	GET string = "GET"
	DEL string = "DEL"
)

type CommandError struct {
	msg string
}

func (c *CommandError) Error() string {
	return fmt.Sprintf("(error) ERR %s", c.msg)
}

type Command struct {
	Keyword string
	key     string
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

func (c Command) Validate() (bool, error) {
	var errMsg string
	var keyword string

	switch c.Keyword {
	case SET:
		keyword = SET
		if c.key == "" {
			errMsg = fmt.Sprintf(
				"%s command expected 2 arguments but none was given (i.e no key & value)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		if c.Value == nil {
			errMsg = fmt.Sprintf("%s command expected 2 arguments but 1 was given (i.e no value)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		return true, nil
	case GET, DEL:
		keyword = GET
		if c.Keyword == DEL {
			keyword = DEL
		}
		if c.key == "" {
			errMsg = fmt.Sprintf("%s command expected 1 argument but none was given (i.e no key)", keyword)
			return false, &CommandError{msg: errMsg}
		}
		if c.Value != nil {
			errMsg = fmt.Sprintf("%s command expected 1 argument but 2 was given", keyword)
			return false, &CommandError{msg: errMsg}
		}
		return true, nil
	}
	return false, &CommandError{
		msg: fmt.Sprintf("unknown command %s", c.Keyword),
	}
}
