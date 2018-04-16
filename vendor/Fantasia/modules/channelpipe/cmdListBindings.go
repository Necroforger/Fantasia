package channelpipe

import (
	"Fantasia/system"

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

	for _, v := range m.Bindings {
		if v.Source.GuildID == guildID {
			content += "`" + v.Source.ChannelID + "` -> `" + v.Sink.ChannelID() + "`\n"
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
