package general

import (
	"github.com/Necroforger/Fantasia/system"
	"strconv"
	"strings"
)

// CmdFromBinary converts the given binary to text
func CmdFromBinary(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("You need to supply a string of binary to convert to text")
		return
	}

	result := []byte{}
	for idx, binstr := range strings.Split(ctx.Args.After(), " ") {
		b, err := strconv.ParseUint(binstr, 2, 64)
		if err != nil {
			ctx.ReplyError("Error parsing uint from string ", idx)
			return
		}
		result = append(result, uint8(b))
	}

	ctx.ReplySuccess(string(result))
}
