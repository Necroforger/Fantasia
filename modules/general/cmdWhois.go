package general

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"
)

// CmdWhois returns information about a user
func CmdWhois(ctx *system.Context) {
	var (
		m        = ctx.Msg
		b        = ctx.Ses
		guild    *discordgo.Guild
		user     *discordgo.User
		presence *discordgo.Presence
		member   *discordgo.Member
		embed    = dream.NewEmbed()
		err      error
	)

	if len(m.Mentions) > 0 {
		user = m.Mentions[0]
	} else if ctx.Args.After() != "" {
		user, err = b.DG.User(ctx.Args.After())
		if err != nil {
			ctx.ReplyError("Error obtaining user from ID")
			return
		}
	} else {
		user = m.Author
	}

	guild, err = b.Guild(m)
	if err != nil {
		ctx.ReplyError(err)
	}

	embed.SetTitle(user.Username)
	if accountcreated, err := dream.CreationTime(user.ID); err == nil {
		embed.AddField("Account created", fmt.Sprintf("`%s` (%s)", accountcreated.Format("2006-01-02"), humanize.Time(accountcreated)))
	}
	embed.AddField("ID", fmt.Sprint(user.ID))
	embed.AddField("Discriminator", user.Discriminator)
	embed.SetThumbnail(user.AvatarURL("2048"))

	if guild != nil {
		if presence, err = b.GuildPresence(guild, user.ID); err == nil {
			if presence.Game != nil {
				embed.AddField("Playing", presence.Game.Name)
			}
			embed.AddField("Status", string(presence.Status))
			embed.SetColor(dream.StatusColor(presence.Status))
		}
		if member, err = b.GuildMember(guild, user.ID); err == nil {
			if joindate, err := discordgo.Timestamp(member.JoinedAt).Parse(); err == nil {
				embed.AddField("Joined guild", fmt.Sprintf("`%s` (%s)", joindate.Format("2006-01-02"), humanize.Time(joindate)))
			}
			if member.Nick != "" {
				embed.AddField("Nickname", member.Nick)
			}
			if len(member.Roles) != 0 {
				if roles, err := b.GuildMemberRoles(member); err == nil {
					var roletext string
					for _, v := range roles {
						roletext += "<@&" + v.ID + "> "
					}
					embed.AddField("Roles", roletext)
				}
			}
		}
	}
	embed.InlineAllFields()
	embed.Truncate()
	_, err = ctx.ReplyEmbed(embed.MessageEmbed)
	if err != nil {
		ctx.ReplyError(err)
		var out bytes.Buffer
		toml.NewEncoder(&out).Encode(embed)
		ctx.ReplyWarning(string(out.Bytes()))
	}
}
