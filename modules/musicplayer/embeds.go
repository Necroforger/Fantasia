package musicplayer

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dream"
)

// Emoji constants
const (
	EmojiStars = "✨"
	EmojiStar  = "⭐"
)

// EmbedQueue returns an embed object for a musicqueue
func EmbedQueue(q *SongQueue, index, beforeIndex, afterIndex int) *dream.Embed {
	embed := dream.NewEmbed()

	for i := index - beforeIndex; i < index+afterIndex; i++ {
		if i < 0 {
			i = 0
		}
		if i < len(q.Playlist) {
			var (
				songname string
				prefix   string
			)
			if i == index {
				songname = q.Playlist[i].Markdown()
			} else {
				songname = q.Playlist[i].String()
			}
			if q.Playlist[i].Rating != 0 {
				prefix = strings.Repeat(EmojiStar, q.Playlist[i].Rating)
			}
			embed.Description += fmt.Sprintf("%d. %s%s\n", i, prefix, songname)

		}
	}

	return embed
}
