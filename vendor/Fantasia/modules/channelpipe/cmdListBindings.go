package channelpipe

import (
	"Fantasia/system"
	"fmt"

	"github.com/Necroforger/dgwidgets"
)

// CmdListBinding lists a guilds bindings
func (m *Module) CmdListBinding(ctx *system.Context) {
	var content string

	guild, err := ctx.Ses.Guild(ctx.Msg)
	if err != nil {
		ctx.ReplyError(err)
	}
	guildID := guild.ID
	content += "`Bindings in " + guild.Name + "`\n"

	m.bmu.RLock()
	defer m.bmu.RUnlock()

	channelName := func(cid string) string {
		c, err := ctx.Ses.DG.State.Channel(cid)
		if err != nil {
			return "unknown"
		}
		return c.Name
	}

	for _, v := range m.Bindings {
		if v.Source.GuildID == guildID {
			content += fmt.Sprintf("`%s(%s)` -> %s(%s)",
				channelName(v.Source.ChannelID), v.Source.ChannelID, channelName(v.Sink.ChannelID()), v.Sink.ChannelID())
		}
	}

	// Spawn a paginator if the message is too long
	if len(content) >= 2000 {
		paginator := dgwidgets.NewPaginator(ctx.Ses.DG, ctx.Msg.ChannelID)
		paginator.Add(dgwidgets.EmbedsFromString(content, 2000)...)
		paginator.Spawn()
		return
	}

	ctx.ReplyNotify(content)
}
