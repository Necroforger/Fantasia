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

// EmbedQueueFilter ...
func EmbedQueueFilter(q *SongQueue, index, beforeIndex, afterIndex int, filterFunc func(*SongQueue, *Song) bool) *dream.Embed {
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
			if filterFunc(q, q.Playlist[i]) {
				if i == q.Index {
					songname = q.Playlist[i].Markdown()
				} else if i == index {
					songname = "`" + q.Playlist[i].String() + "`"
				} else {
					songname = q.Playlist[i].String()
				}
				if q.Playlist[i].Rating != 0 {
					prefix = strings.Repeat(EmojiStar, q.Playlist[i].Rating)
				}
				embed.Description += fmt.Sprintf("%d. %s%s\n", i, prefix, songname)
			}
		}
	}

	return embed
}

// EmbedQueue returns an embed object for a musicqueue
func EmbedQueue(q *SongQueue, index, beforeIndex, afterIndex int) *dream.Embed {
	return EmbedQueueFilter(q, index, beforeIndex, afterIndex, func(q *SongQueue, s *Song) bool { return true })
}
