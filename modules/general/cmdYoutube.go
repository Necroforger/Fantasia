package general

import (
	"fmt"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/youtubeapi"
)

// CmdYoutube searches youtube for a video
func CmdYoutube(ctx *system.Context) {
	videoURLS := []string{}

	if ctx.System.Config.GoogleAPIKey != "" {
		res, err := youtubeapi.New(ctx.System.Config.GoogleAPIKey).Search(ctx.Args.After(), 10)
		if err != nil {
			ctx.ReplyError(err)
			return
		}

		for _, v := range res.Items {
			if v.ID.Kind != "youtube#video" {
				continue
			}
			videoURLS = append(videoURLS, fmt.Sprint("http://youtube.com/watch?v=", v.ID.VideoID))
		}
	} else {
		res, err := youtubeapi.ScrapeSearch(ctx.Args.After(), 10)
		if err != nil {
			ctx.ReplyError(err)
			return
		}
		videoURLS = append(videoURLS, res...)
	}

	if len(videoURLS) > 0 {
		ctx.Reply(videoURLS[0])
	}
}
