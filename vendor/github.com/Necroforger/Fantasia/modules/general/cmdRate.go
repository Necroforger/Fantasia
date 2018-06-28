package general

import (
	"hash/fnv"
	"math/rand"

	"github.com/Necroforger/Fantasia/system"
)

// CmdRate rates the supplied query as an integer between zero and ten
func CmdRate(ctx *system.Context) {
	var (
		result = 0
		q      = ctx.Args.After()
	)

	if q == "" {
		ctx.ReplyError("Please give me something to rate")
		return
	}

	h := fnv.New32a()
	h.Write([]byte(q))
	generator := rand.New(rand.NewSource(int64(h.Sum32())))
	result = int(generator.Float64() * 11)

	ctx.ReplyNotify("I rate **", q, "**: *", result, " / 10*")
}
