package eval

import (
	"errors"
	"fmt"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
	"github.com/robertkrimen/otto"
)

var errVMTimeout = errors.New("error: code execution halted for taking too long")

// EvalJS ...
func (m *Module) EvalJS(ctx *system.Context) {
	script := ctx.Args.After()
	vm := otto.New()
	b := ctx.Ses

	// Dangerous: Gives full access to Command router, discordgo session, tokens etc...
	//evalJSSetFunctions(ctx, vm)

	if len(script) != 0 {
		ctx.ReplyEmbed(evalJSEmbed(vm, script, time.Second*1).MessageEmbed)
		return
	}

	b.SendEmbed(ctx.Msg, dream.NewEmbed().
		SetTitle("Entered javascript interpreter").
		SetDescription("Type `exit` to leave"))
	for {

		msg := b.NextMessageCreate()
		fmt.Println(msg.Content)

		if msg.Author.ID != ctx.Msg.Author.ID {
			continue
		}

		if msg.Content == "exit" {
			ctx.ReplyStatus(system.StatusNotify, "Left javascript interpreter")
			return
		}

		b.SendEmbed(msg.ChannelID, evalJSEmbed(vm, msg.Content, time.Second*1))
	}
}

func evalJSEmbed(vm *otto.Otto, script string, timeout time.Duration) *dream.Embed {
	embed := dream.NewEmbed()

	res, err := evalJS(vm, script, timeout)

	if err != nil {
		embed.Description = err.Error()
		embed.Color = system.StatusError
	} else {
		embed.Description = res
		embed.Color = system.StatusSuccess
	}

	return embed
}

func evalJS(vm *otto.Otto, script string, timeout time.Duration) (result string, err error) {

	resChan := make(chan string)
	errChan := make(chan error)
	timeoutChan := make(chan error)
	vm.Interrupt = make(chan func(), 1)

	result = "error: timed out"

	go func() {
		defer func() {
			if v := recover(); v != nil {
				if v == errVMTimeout {
					return
				}
				panic(v)
			}
		}()

		res, er := vm.Run(script)
		if er != nil {
			errChan <- er
		} else {
			resChan <- res.String()
		}

	}()

	go func() {
		time.Sleep(timeout)
		timeoutChan <- errVMTimeout
	}()

	select {
	case result = <-resChan:
	case err = <-timeoutChan:
		vm.Interrupt <- func() {
			panic(errVMTimeout)
		}
	case err = <-errChan:
	}

	return
}

func evalJSSetFunctions(ctx *system.Context, vm *otto.Otto) {
	vm.Set("ctx", ctx)
	vm.Run(`function view(data) { return JSON.stringify(data, "", "\t"); }`)
}
