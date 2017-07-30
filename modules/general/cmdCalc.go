package general

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/Necroforger/Fantasia/system"
)

// CmdCalc calculates the given input
func CmdCalc(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter an expression")
		return
	}

	expression, err := govaluate.NewEvaluableExpression(ctx.Args.After())
	result, err := expression.Evaluate(make(map[string]interface{}, 0))
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	ctx.ReplyNotify(fmt.Sprintf("**in:**\t`%s`\n\n**out:**\t%s\n", ctx.Args.After(), fmt.Sprint(result)))
}
