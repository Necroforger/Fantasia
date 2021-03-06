package general

import (
	"math/rand"
	"time"

	"github.com/Necroforger/Fantasia/system"

	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r := s.CommandRouter
	r.On("avatar", CmdAvatar).Set("", "Retrieves your avatar or the avatar of the user with the given ID.\n`avatar [userid]`")
	r.On("whois", CmdWhois).Set("", "Retrieves information abour a user. If called without any arguments, it will retrieve your user info\n`whois [userid | user mention]`")
	r.On("serverinfo", CmdServerInfo).Set("", "Retrieves information about the current server")
	r.On("quote", CmdQuote).Set("", "Quotes a user by message id.\n'quote [messageID]`")
	r.On("emojify", m.emojifyCommand).Set("", "Emojifies the given text")
	r.On("ping", Ping).Set("", "responds with the amount of time taken to send and retrieve a message")
	r.On("snowflake", Snowflake).Set("", "gives the creation date of a discord ID")
	r.On("youtube", CmdYoutube).Set("", "Searches for the youtube video with the give title.\n`youtube video name`")
	r.On("hex", CmdHexDisplay).Set("", "Returns an image representation of the given hex code. example: `hex ff00ff`")
	r.On("hexcheck", CmdHexCheck).Set("", "Returns the hex value of the center pixel of the given image")
	r.On("remind", CmdRemind).Set("", "Reminds you about something after a duration set in seconds.\n`remind 10 hello`")
	r.On("calc", CmdCalc).Set("", "Calculates the given expression.\n`calc 10 + 10`")
	r.On("tpb", CmdPirateBay).Set("", "Searches the pirate bay for the given query\n`tpb music`")
	r.On("gog", CmdGog).Set("", "Searches [GOG](http://gog.com/) for the given query")
	r.On("invite", CmdInviteURL).Set("", "Obtain an OAUTH2 invite URL for the bot")

	// Conversion
	r.On("tobinary", CmdToBinary).Set("", "convert the given text to binary")
	r.On("frombinary", CmdFromBinary).Set("", "converts the given text from binary to utf-8 text")

	// Random
	r.On("rate", CmdRate).Set("", "Rates the supplied thing on a scale of 1-10")
	r.On("8ball", CmdEightBall).Set("", "Query the magic eightball")
	r.On("choose", CmdChoose).Set("", "Choose between a list of options. Enter your choices as a comma separated list. ex. `choose Orchestral arrangements, Piano arrangements`")

	rand.Seed(time.Now().UnixNano())
}

// Ping returns the time taken to send a message and recieve back the discord event
func Ping(ctx *system.Context) {
	embed := dream.NewEmbed().
		SetDescription("pinging").
		SetColor(system.StatusWarning).
		MessageEmbed

	succEmbed := dream.NewEmbed().
		SetColor(system.StatusSuccess).
		MessageEmbed

	start := time.Now()
	m, err := ctx.Ses.DG.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	if err == nil {
		succEmbed.Description = time.Since(start).String()
		ctx.Ses.DG.ChannelMessageEditEmbed(m.ChannelID, m.ID, succEmbed)
	}
}

// Snowflake gives the creation date of a discord ID
func Snowflake(ctx *system.Context) {
	t, err := dream.CreationTime(ctx.Args.After())
	if err == nil {

		ctx.ReplyEmbed(dream.NewEmbed().
			SetTitle(ctx.Args.After()).
			SetDescription(
				t.Format(
					time.RFC1123) + "\n" +
					humanize.Time(t),
			).
			SetColor(system.StatusNotify).
			MessageEmbed)
	}
}
