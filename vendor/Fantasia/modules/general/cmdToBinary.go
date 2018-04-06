package general

import (
	"Fantasia/system"
	"strconv"
	"strings"
)

const binaryPadding = 8

// CmdToBinary converts text to binary
func CmdToBinary(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("You need to give text to convert to binary")
		return
	}

	results := []string{}
	for _, b := range []byte(ctx.Args.After()) {
		results = append(
			results,
			padStringLeft(
				strconv.FormatUint(
					uint64(b),
					2,
				),
				"0",
				8,
			),
		)
	}

	ctx.ReplySuccess(strings.Join(results, " "))
}

func padStringLeft(data string, pad string, npad int) string {
	if npad <= 0 {
		return data
	}
	if len(data) == npad {
		return data
	}
	return strings.Repeat(pad, (npad-(len(data)%npad))%npad) + data
}
