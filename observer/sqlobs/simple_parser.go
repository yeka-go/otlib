package otsql

import (
	"errors"
	"strings"
)

func Parse(q string) (op, table string, err error) {
	q = strings.Replace(q, "\n", " ", -1)
	p := strings.Split(strings.ToUpper(strings.TrimRight(q, " ")), " ")
	if len(p) < 1 {
		err = errors.New("unknown query")
		return
	}

	switch p[0] {
	case "SELECT":
		op = p[0]

	case "UPDATE":
		op = p[0]

	case "DELETE":
		op = p[0]

	case "INSERT":
		op = p[0]

	default:
		err = errors.New("unknown operation")
		return
	}

	return
}
