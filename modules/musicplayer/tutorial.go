package musicplayer

import (
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/widgets"
	"github.com/Necroforger/dream"
)

// CmdTutorial sends a tutorial with instructions on how to use the musicplayer
func (m *Module) CmdTutorial(ctx *system.Context) {
	paginator := widgets.NewPaginator(ctx.Ses, ctx.Msg.ChannelID)

	paginator.
		Add(dream.NewEmbed().SetTitle("Page 1").
			SetDescription("Paginator test").MessageEmbed,
			dream.NewEmbed().SetTitle("Page 2").
				SetDescription("This is page two's content").
				SetImage("https://pbs.twimg.com/profile_images/867187199804133377/Qwnf-x4j_400x400.jpg").MessageEmbed,
			dream.NewEmbed().SetTitle("Page 3").
				SetDescription("Nothing to see here").MessageEmbed)

	paginator.SetPageFooters()
	paginator.NavigationTimeout = time.Minute * 5
	paginator.Spawn()
}
