package main

import "strings"

func (s cell) String() string {
	var res []byte
	if s.y.crossable {
		res = []byte("   ")
	} else {
		res = []byte("___")
	}

	if !s.x.crossable {
		res[2] = '|'
	}

	return string(res)
}

func (r row) Array() []byte {
	var line []byte
	for _, c := range r {
		line = append(line, []byte(c.String())...)
	}
	return line
}

func (r row) String() string {
	return string(r.Array())
}

func (m Maze) String() string {
	var lines []string
	for _, row := range m.area {
		lines = append(lines, row.String()[m.edge:])
	}
	// prev := lines[0]
	// var res = []string{string(prev)}
	// for _, row := range
	return strings.Join(lines, "\n")
}
