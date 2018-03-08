package general

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"Fantasia/system"
)

// CmdCalc calculates the given input
func CmdCalc(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter an expression")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			ctx.ReplyError("Error evaluating expression")
		}
	}()

	expression, err := govaluate.NewEvaluableExpression(ctx.Args.After())
	result, err := expression.Evaluate(nil)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	ctx.ReplyNotify(fmt.Sprintf("**in:**\t`%s`\n\n**out:**\t%s\n", ctx.Args.After(), fmt.Sprint(result)))
}
