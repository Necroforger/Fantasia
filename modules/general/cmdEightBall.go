package general

import (
	"math/rand"

	"github.com/Necroforger/Fantasia/system"
)

// CmdEightballAnswers stores answers for eightball
var CmdEightballAnswers = []string{
	"it is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes, definitely",
	"You may rely on it",
	"As I see it, yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

// CmdEightBall returns a variety of eightball answers.
func CmdEightBall(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter a query")
		return
	}
	ctx.ReplyNotify(CmdEightballAnswers[int(rand.Float64()*float64(len(CmdEightballAnswers)))])
}
