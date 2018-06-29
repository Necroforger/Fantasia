package guildconfig

import (
	"fmt"
	"log"
	"strings"

	"github.com/Necroforger/Fantasia/models"
	"github.com/Necroforger/Fantasia/system"
)

// Module ...
type Module struct {
}

// Build builds the module
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter

	t, err := system.NewSubCommandRouter(`^config(\s|$)`, "config")
	if err != nil {
		log.Println(err)
		return
	}
	t.Router.Prefix = "^"
	r.AddSubrouter(t)

	t.CommandRoute = &system.CommandRoute{
		Name:    "config",
		Desc:    "configures guild settings",
		Handler: Auth(CmdConfig),
	}

	k := t.Router
	k.On("prefix", Auth(CmdPrefix)).Set("", "sets the guild command prefix")
	k.On("admins", Auth(CmdAdmins)).Set("", "sets the admin list")
}

const flagDefault = "--default"

const help = `**Config:**
To print the current value of a setting use: **"config [command]"**
To update the value of a setting use: **"config [command] [value]"**

To reset a flag to its default setting, use **` + flagDefault + `** for the value

**Prefix** : Controls the guild specific prefix for this bot

**Admins** : Admins can edit the config and use any command
         Pass a comma separated list of user IDs with no
         spaces to update this field
`

// CmdConfig allows setting of config
func CmdConfig(ctx *system.Context) {
	ctx.ReplyNotify(help)
}

// SetString sets a string value
func SetString(ctx *system.Context, gconfig *models.Guild, name, defaultVal string, value *string) {
	if ctx.Args.After() != "" {
		if ctx.Args.After() == flagDefault {
			*value = defaultVal
		} else {
			*value = ctx.Args.After()
		}

		err := ctx.System.DB.SaveGuild(ctx.Msg.GuildID, gconfig)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
	}

	ctx.ReplyNotify(fmt.Sprintf("%s: `%s`", name, *value))
}

// SetStrings sets a []string variable
func SetStrings(ctx *system.Context, gconfig *models.Guild, name string, defaultVal []string, value *[]string) {
	if ctx.Args.After() != "" {
		if ctx.Args.After() == flagDefault {
			*value = defaultVal
		} else {
			*value = strings.Split(ctx.Args.After(), ",")
		}

		err := ctx.System.DB.SaveGuild(ctx.Msg.GuildID, gconfig)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
	}

	ctx.ReplyNotify(fmt.Sprintf("%s:\n%s", name, strings.Join(gconfig.Admins, ",")))
}

// CmdPrefix sets a guild command prefix
func CmdPrefix(ctx *system.Context) {
	gconfig := ctx.Get("gconfig").(*models.Guild)
	SetString(ctx, gconfig, "Prefix", "", &gconfig.Prefix)
}

// CmdAdmins sets the list of admins
func CmdAdmins(ctx *system.Context) {
	gconfig := ctx.Get("gconfig").(*models.Guild)
	SetStrings(ctx, gconfig, "Admins", []string{}, &gconfig.Admins)
}

// Auth is authentication middleware
func Auth(fn func(ctx *system.Context)) func(ctx *system.Context) {
	return func(ctx *system.Context) {
		gconfig, err := ctx.System.DB.CreateGuildIfNotExists(ctx.Msg.GuildID)
		if err != nil {
			ctx.ReplyError("Error getting guild configuration: ", err)
			return
		}

		isAdmin, err := ctx.IsAdmin()
		if err != nil {
			ctx.ReplyError("Error checking administrator status: ", err)
			return
		}

		if !isAdmin {
			ctx.ReplyError("You need to be an administrator or own the guild to configure guild settings")
			return
		}

		ctx.Set("gconfig", gconfig)
		fn(ctx)
	}
}
