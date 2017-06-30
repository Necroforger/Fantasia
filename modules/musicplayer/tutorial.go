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

	QuickGuide := dream.NewEmbed().
		SetTitle("Quick guide").
		SetDescription("Get started quickly\n" +
			"`>m queue https://www.youtube.com/watch?v=dm2YwytOyH0`\n" +
			"`>m play`").
		SetImage("http://i.imgur.com/7MLMNmM.gif").
		SetColor(system.StatusNotify).
		MessageEmbed

	QueueHelp := dream.NewEmbed().
		SetTitle("Queueing videos").
		SetImage("http://i.imgur.com/Bs6drYI.gif").
		SetDescription("`queue [ URL | playlist URL | search query | song index]`\n"+
			"If no arguments are provided, list the playlist").
		AddField("URL",
			"Use youtube-dl to queue the requested url. can be a link to soundcloud, youtube, facebook, or any other site supported by youtube-dl").
		AddField("Playlist URL",
			"If you provide a playlist URL, it will add every video in the playlist to the queue").
		AddField("Search Query", "If the argument is not a URL it will be treated as a youtube search query.\n"+
			"The bot will search youtube with the query and queue the first result found").
		AddField("Song Index", "If the argument is a song index integer it will list the songs around the given song index.").
		SetColor(system.StatusNotify).
		InlineAllFields().
		MessageEmbed

	SaveAndLoad := dream.NewEmbed().
		SetTitle("Save and load playlists").
		SetDescription("Use the save and load commands to save and load playlists").
		SetImage("http://i.imgur.com/CJ3tQvj.gif").
		SetColor(system.StatusNotify).
		MessageEmbed

	Navigating := dream.NewEmbed().
		SetTitle("Navigating the queue").
		SetDescription("Navigate the queue using the `next`, `previous`, and multipurpose `go` command" +
			"Any command that accepts a song index as an argument (except the queue command) can use one of the following strings " +
			"in its place\n\n" +
			"`start, beginning : select index 0`\n" +
			"`end, last : select the last index`\n" +
			"`middle, center: select the middle index`\n" +
			"`random, rand : select a random index`").
		SetImage("http://i.imgur.com/KaaRGBb.gif").
		SetColor(system.StatusNotify).
		MessageEmbed

	paginator.Add(QuickGuide, QueueHelp, Navigating, SaveAndLoad)
	paginator.SetPageFooters()
	paginator.NavigationTimeout = time.Minute * 5
	paginator.Spawn()
}
