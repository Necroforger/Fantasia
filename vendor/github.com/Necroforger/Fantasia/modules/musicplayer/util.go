package musicplayer

import (
	"strconv"
	"strings"

	"github.com/Necroforger/Fantasia/system"
)

// getIndexes creates an ID list from the supplied arguments.
// Used for dealing with playlist queues
func getIndexes(args []string, radio *Radio) []int {
	ids := []int{}
	for _, arg := range args {

		// Check for range of numbers
		if strings.Contains(arg, "-") {
			if nums := strings.Split(arg, "-"); len(nums) > 1 && nums[0] != "" && nums[1] != "" {
				if n1, err := getIndex(nums[0], radio); err == nil {
					if n2, err := getIndex(nums[1], radio); err == nil {
						for i := n1; i <= n2; i++ {
							ids = append(ids, i)
						}
					}
				}
			}
		} else if num, err := getIndex(arg, radio); err == nil {
			ids = append(ids, num)
		}

	}
	return ids
}

func getIndex(index string, radio *Radio) (int, error) {
	switch index {
	case "start", "beginning":
		return 0, nil
	case "end", "last":
		return len(radio.Queue.Playlist) - 1, nil
	case "mid", "center", "middle":
		return len(radio.Queue.Playlist) / 2, nil
	case "rand", "random":
		return int(rng.Float64() * float64(len(radio.Queue.Playlist)-1)), nil
	case "current", "playing":
		return radio.Queue.Index, nil
	case "next":
		return radio.Queue.Index + 1, nil
	case "prev", "previous":
		return radio.Queue.Index - 1, nil
	default:
		return strconv.Atoi(index)
	}

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

// ProgressBar generates a text progress bar.
//    value: Current progress
//    end  : End value
//    scale: Size of the progress bar
func ProgressBar(value, end, scale int) string {
	const (
		spaceChar = "-"
		fillChar  = "#"
	)
	if end == 0 {
		return "[" + strings.Repeat(spaceChar, scale) + "]"
	}
	if value >= end {
		return "[" + strings.Repeat(fillChar, scale) + "]"
	}

	num := (float64(value) / float64(end)) * float64(scale)
	numrem := (1 - (float64(value) / float64(end))) * float64(scale)

	return "[" + strings.Repeat(fillChar, int(num)) + strings.Repeat(spaceChar, int(numrem)) + "]"
}
