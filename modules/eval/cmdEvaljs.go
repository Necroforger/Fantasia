package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dgwidgets"
	"github.com/Necroforger/discordgo"
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
	for _, v := range ctx.System.Config.Admins {
		if v == ctx.Msg.Author.ID {
			evalJSSetFunctions(ctx, vm)
			break
		}
	}

	if len(script) != 0 {
		sendResult(ctx, ctx.Msg.ChannelID, evalJSEmbed(vm, script, time.Second*1).MessageEmbed)
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
		chunklen := 1024
		embed := evalJSEmbed(vm, msg.Content, time.Second*1)
		if len(embed.Description) < chunklen {
			b.SendEmbed(msg.ChannelID, embed)
		} else {
			sendResult(ctx, msg.ChannelID, embed.MessageEmbed)
		}
	}
}

func sendResult(ctx *system.Context, channelID string, embed *discordgo.MessageEmbed) {

	if len(embed.Description) <= 1024 {
		ctx.ReplyEmbed(embed)
	}

	p := dgwidgets.NewPaginator(ctx.Ses.DG, channelID)
	p.Add(dgwidgets.EmbedsFromString(embed.Description, 1024)...)

	p.Widget.Handle("💾", func(w *dgwidgets.Widget, r *discordgo.MessageReaction) {
		if r.UserID != ctx.Msg.Author.ID {
			return
		}
		if userchan, err := w.Ses.UserChannelCreate(r.UserID); err == nil {
			var content string
			extension := "txt"
			for _, v := range p.Pages {
				content += v.Description
			}
			var js interface{}
			if err := json.Unmarshal([]byte(content), &js); err == nil {
				if t, err := json.MarshalIndent(js, "", "\t"); err == nil {
					content = string(t)
					extension = "json"
				}
			}
			w.Ses.ChannelFileSend(userchan.ID, "content."+extension, bytes.NewReader([]byte(content)))
		}
	})

	for _, v := range p.Pages {
		v.Color = embed.Color
	}

	p.SetPageFooters()
	p.Widget.Timeout = time.Minute * 2
	p.ColourWhenDone = system.StatusWarning

	go p.Spawn()
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
