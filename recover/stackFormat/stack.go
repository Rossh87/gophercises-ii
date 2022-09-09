package stackFormat

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func readByLine(s string, test func(s string) bool, format func(s string) string) string {
	out := strings.Builder{}

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		line := scanner.Text()
		if test(line) {
			line = format(line)
		}

		out.Write([]byte(line))
	}

	return out.String()
}

func parseStackPath(stack string) stackPath {
	fields := strings.Fields(stack)
	data := strings.Split(fields[0], ":")

	ln, err := strconv.ParseInt(data[1], 10, 32)

	if err != nil {
		panic(err)
	}

	return stackPath{
		int(ln),
		data[0],
		stack,
	}
}

type stackPath struct {
	lineNo   int
	filePath string
	rawPath  string
}

func formatStackPath(sp stackPath) string {
	return fmt.Sprintf("\n<a href=\"/source%s?lineNumber=%d\">%s</a>\n", sp.filePath, sp.lineNo, sp.rawPath)
}

func FormatStack(stack string) string {
	test := func(s string) bool {
		return strings.Trim(s, "\t ")[0] == '/'
	}

	format := func(s string) string {
		return formatStackPath(parseStackPath((s)))
	}

	return readByLine(stack, test, format)
}
