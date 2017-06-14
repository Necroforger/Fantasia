package musicplayer

import (
	"strconv"
	"strings"

	"github.com/Necroforger/Fantasia/system"
)

// createIDList creates an ID list from the supplied arguments.
// Used for dealing with playlist queues
func createIDList(args []string) []int {
	ids := []int{}
	for _, arg := range args {

		// Check for range of numbers
		if strings.Contains(arg, "-") {
			if nums := strings.Split(arg, "-"); len(nums) > 1 && nums[0] != "" && nums[1] != "" {
				if n1, err := strconv.Atoi(nums[0]); err == nil {
					if n2, err := strconv.Atoi(nums[1]); err == nil {
						for i := n1; i <= n2; i++ {
							ids = append(ids, i)
						}
					}
				}
			}
		} else if num, err := strconv.Atoi(arg); err == nil {
			ids = append(ids, num)
		}

	}
	return ids
}

func guildIDFromContext(ctx *system.Context) (string, error) {
	var guildID string

	vs, err := ctx.Ses.UserVoiceState(ctx.Msg.Author.ID)
	if err != nil {
		guildID, err = ctx.Ses.GuildID(ctx.Msg)
		if err != nil {
			return "", err
		}
	} else {
		guildID = vs.GuildID
	}

	return guildID, nil
}

//ProgressBar generates a progressbar given a value, end point, and scale
func ProgressBar(value, end, scale int) string {
	if end == 0 {
		return "[" + strings.Repeat("-", scale) + "]"
	}
	if value >= end {
		return "[" + strings.Repeat("#", scale) + "]"
	}

	num := (float64(value) / float64(end)) * float64(scale)
	numrem := (1 - (float64(value) / float64(end))) * float64(scale)

	return "[" + strings.Repeat("#", int(num)) + strings.Repeat("=", int(numrem)) + "]"
}
