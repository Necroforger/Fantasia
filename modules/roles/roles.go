package roles

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

// Module ...
type Module struct{}

// Build ...
func (m *Module) Build(s *system.System) {
	r, _ := system.NewSubCommandRouter("^role", "role")
	r.Set("role", "subrouter for role commands. example useage: `role color [hex]`")
	s.CommandRouter.AddSubrouter(r)

	r.Router.On("color|colour", m.Color).Set("", "Changes your role colour to the supplied hex code")
}

// Color ...
func (m *Module) Color(ctx *system.Context) {
	rolename := ctx.Msg.Author.ID
	b := ctx.Ses

	guild, err := b.Guild(ctx.Msg)
	if err != nil {
		b.SendEmbed(ctx.Msg, "Error obtaining user's guild")
		return
	}

	if ctx.Args.After() == "" {
		err = b.GuildRoleDeleteByName(guild, rolename)
		if err != nil {
			b.SendEmbed(ctx.Msg, "Failed to remove your coloured role")
		} else {
			b.SendEmbed(ctx.Msg, "Your role colour has been reset")
		}
		return
	}

	rolecolor, err := strconv.ParseInt(ctx.Args.After(), 16, 64)
	if err != nil {
		ctx.ReplyError("Error parsing supplied hex: " + ctx.Args.After())
		return
	}

	// Find or create the coloured role
	role, err := editOrCreateRoleIfNotExist(b, guild, rolename, int(rolecolor))
	if err != nil {
		ctx.ReplyError("Error editing or creating role if it does not exist")
		return
	}

	// Attempt to add the role to the member if it doesn't already exist
	err = addRoleToMemberIfNotExist(b, guild.ID, ctx.Msg.Author.ID, role.ID)
	if err != nil {
		ctx.ReplyError("Error adding role to member if not exist")
		return
	}

	// Reposition the role to the top of the role list so that it becomes
	// The user's primary colour.
	err = moveRoleToTop(b, guild.ID, role.ID)
	if err != nil {
		b.SendEmbed(ctx.Msg, "Error repositioning coloured role beneath bot role: "+fmt.Sprint(err))
		return
	}

	b.SendEmbed(ctx.Msg, dream.NewEmbed().
		SetTitle(ctx.Msg.Author.Username+": Role colour changed to ["+ctx.Args.After()+"]").
		SetColor(int(rolecolor)),
	)

}

func editOrCreateRoleIfNotExist(b *dream.Bot, guild *discordgo.Guild, rolename string, roleColor int) (*discordgo.Role, error) {
	guildRoles, err := b.GuildRoles(guild)
	if err != nil {
		return nil, err
	}

	for _, v := range guildRoles {
		if v.Name == rolename {
			_, err = b.GuildRoleEdit(guild.ID, v.ID, dream.RoleSettings{
				Name:  v.Name,
				Color: roleColor,
			})
			if err != nil {
				return nil, err
			}
			return v, nil
		}
	}

	role, err := b.GuildRoleCreate(guild.ID, dream.RoleSettings{
		Name:  rolename,
		Color: roleColor,
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func addRoleToMemberIfNotExist(b *dream.Bot, guildID, userID, roleID string) error {
	memberRoles, err := b.GuildMemberRoles(guildID, userID)
	if err != nil {
		return err
	}

	// Return if the member has the given role
	for _, v := range memberRoles {
		if v.ID == roleID {
			return nil
		}
	}

	err = b.DG.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func moveRoleToTop(b *dream.Bot, guildID, roleID string) error {
	guildRoles, err := b.GuildRoles(guildID)
	if err != nil {
		return err
	}

	// Find the highest client role.
	clientHighest, err := highestMemberRolePosition(b, guildID, b.DG.State.User.ID)
	if err != nil {
		return err
	}

	roles := dream.Roles(guildRoles)

	// Sort roles and set positions accordingly
	sort.Sort(roles)

	// Move colour role to the one below the client's highest role
	err = roles.MoveByID(roleID, clientHighest-1)
	if err != nil {
		return err
	}

	// Update the role positions to reflect their positions in the slice
	roles.UpdatePositions()

	_, err = b.DG.GuildRoleReorder(guildID, roles)
	return err
}

func highestMemberRolePosition(b *dream.Bot, guildID, userID string) (int, error) {
	memberRoles, err := b.GuildMemberRoles(guildID, userID)
	if err != nil {
		return -1, err
	}

	if len(memberRoles) == 0 {
		return -1, errors.New("Member has no roles")
	}

	roles := dream.Roles(memberRoles)
	sort.Sort(roles)
	return roles[len(roles)-1].Position, nil
}
