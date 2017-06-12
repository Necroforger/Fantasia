package musicplayer

import (
	"fmt"

	"github.com/Necroforger/dream"
)

// EmbedQueue returns an embed object for a musicqueue
func EmbedQueue(q *SongQueue, index, beforeIndex, afterIndex int) *dream.Embed {
	embed := dream.NewEmbed()

	for i := index - beforeIndex; i < index+afterIndex; i++ {
		if i < 0 {
			i = 0
		}
		if i < len(q.Playlist) {
			if i == index {
				embed.Description += fmt.Sprintf("%d. [%s](%s)\n", i, q.Playlist[i].String(), q.Playlist[i].URL)
			} else {
				embed.Description += fmt.Sprintf("%d. %s\n", i, q.Playlist[i].String())
			}
		}
	}

	return embed
}
