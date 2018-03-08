package general

import (
	"fmt"
	"strconv"
	"time"

	"Fantasia/system"
	"github.com/bwmarrin/discordgo"
	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"
)

// CmdServerInfo displays information about the current server
func CmdServerInfo(ctx *system.Context) {
	guild, err := ctx.Ses.Guild(ctx.Msg)
	if err != nil {
		ctx.ReplyError("Error obtaining guild information")
		return
	}

	var channelCount int
	var voiceChannelCount int

	for _, ch := range guild.Channels {
		if ch.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannelCount++
		} else {
			channelCount++
		}
	}

	embed := dream.NewEmbed().
		SetThumbnail(discordgo.EndpointGuildIcon(guild.ID, guild.Icon)).
		SetTitle(guild.Name).
		SetURL(discordgo.EndpointGuildIcon(guild.ID, guild.Icon)).
		SetDescription(fmt.Sprintf("*[%d members] [%d channels] [%d voice channels]*\n", guild.MemberCount, channelCount, voiceChannelCount)).
		AddField("Owner", "<@"+guild.OwnerID+">").
		AddField("Region", guild.Region).
		AddField("Online", strconv.Itoa(len(guild.Presences))).
		SetColor(system.StatusNotify)

	//////////////////////
	// Creation Time
	/////////////////////
	if creationtime, err := dream.CreationTime(guild.ID); err == nil {
		embed.AddField("Created", fmt.Sprint(creationtime.Format(time.RFC1123), "  (", humanize.Time(creationtime), ")"))
	}

	/////////////////////
	// List emojis
	////////////////////
	emojiList := ""
	for _, em := range guild.Emojis {
		emojiList += "<:" + em.Name + ":" + em.ID + ">"
	}
	if len(emojiList) > 0 {
		embed.AddField("Emojis ["+strconv.Itoa(len(guild.Emojis))+"]", emojiList)
	}

	/////////////////////////
	//    List Roles
	/////////////////////////
	roleList := ""
	for _, rl := range guild.Roles {
		roleList += "<@&" + rl.ID + ">  "
	}
	if len(roleList) > 0 {
		embed.AddField("Roles ["+strconv.Itoa(len(guild.Roles))+"]", roleList)
	}

	embed.InlineAllFields()
	embed.Truncate()
	ctx.ReplyEmbed(embed.MessageEmbed)
}
