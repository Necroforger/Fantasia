package general

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Necroforger/dream"

	"Fantasia/system"

	"github.com/PuerkitoBio/goquery"
)

// PiratebayItem ...
type PiratebayItem struct {
	Title       string
	URL         string
	Magnet      string
	Seeders     int
	Leechers    int
	Category    string
	Size        string
	Uploader    string
	Uploaded    string
	Description string
	Subcategory string
}

// PiratebayScrape does a piratebay search for the given query
func PiratebayScrape(query string, page int) ([]PiratebayItem, error) {
	items := []PiratebayItem{}

	searchURL := func(query string) string {
		return fmt.Sprintf("https://thepiratebay.org/search/%s/%d", url.QueryEscape(query), page)
	}

	doc, err := goquery.NewDocument(searchURL(query))
	if err != nil {
		return items, err
	}

	doc.Find("#searchResult").
		Find("tbody").
		Find("tr").
		Not(".alt").
		Each(func(_ int, row *goquery.Selection) {
			item := PiratebayItem{}

			td := row.Find("td")
			// Parse category and subcategory
			if len(td.Nodes) > 0 {
				n := td.Nodes[0]
				if c := goquery.NewDocumentFromNode(n).Find("a").Nodes; len(c) > 0 {
					item.Category = goquery.NewDocumentFromNode(c[0]).Text()
					if len(c) > 1 {
						item.Subcategory = goquery.NewDocumentFromNode(c[1]).Text()
					}
				}
			}

			// parse title and url
			if len(td.Nodes) > 1 {
				n := td.Nodes[1]
				if c := goquery.NewDocumentFromNode(n).Find("a").Nodes; len(c) > 0 {
					// URL and Title
					item.Title = goquery.NewDocumentFromNode(c[0]).Text()
					if val, ok := goquery.NewDocumentFromNode(c[0]).Attr("href"); ok {
						item.URL = "https://thepiratebay.org" + val
					}

					// Magnet link
					if len(c) > 1 {
						item.Magnet = goquery.NewDocumentFromNode(c[1]).AttrOr("href", "")
					}

					// Username
					if len(c) > 2 {
						item.Uploader = goquery.NewDocumentFromNode(n).Find(".detDesc").Find("a").Text()
					}

				}
			}
			// Parse seeders
			if len(td.Nodes) > 2 {
				if n, err := strconv.Atoi(goquery.NewDocumentFromNode(td.Nodes[2]).Text()); err == nil {
					item.Seeders = n
				}
			}
			// Parse leechers
			if len(td.Nodes) > 3 {
				if n, err := strconv.Atoi(goquery.NewDocumentFromNode(td.Nodes[3]).Text()); err == nil {
					item.Leechers = n
				}
			}
			items = append(items, item)
		})

	return items, nil
}

// PiratebayInfo gives extra information to a torrent
func PiratebayInfo(URL string, item *PiratebayItem) error {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return err
	}

	item.Description = doc.Find(".nfo").Text()
	item.Size = doc.Find("dt:contains('Size:')").Next().Text()
	item.Uploaded = doc.Find("dt:contains('Uploaded:')").Next().Text()

	return nil
}

// CmdPirateBay searches pirate bay for the given query
func CmdPirateBay(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter a search query")
		return
	}

	ctx.Ses.DG.ChannelTyping(ctx.Msg.ChannelID)

	res, err := PiratebayScrape(ctx.Args.After(), 0)
	if err != nil {
		ctx.ReplyError(err)
	}

	cut := func(str string, limit int) string {
		if len(str) > limit {
			return str[:limit] + "..."
		}
		return str
	}

	embedItem := func(r PiratebayItem) *dream.Embed {
		return dream.NewEmbed().
			SetTitle(r.Title).
			SetURL(r.URL).
			SetDescription(r.Magnet).
			AddField("Uploader", "`"+r.Uploader+"`").
			AddField("Seeders", "`"+strconv.Itoa(r.Seeders)+"`").
			AddField("Leechers", "`"+strconv.Itoa(r.Leechers)+"`").
			AddField("Size", "`"+r.Size+"`").
			AddField("Uploaded", "`"+r.Uploaded+"`").
			AddField("Description", "```"+cut(r.Description, 600)+"```").
			SetColor(system.StatusNotify).
			InlineAllFields()
	}

	if len(res) == 0 {
		ctx.ReplyError("No results found")
		return
	}

	if len(res) == 1 {
		r := res[0]
		PiratebayInfo(r.URL, &r)
		ctx.ReplyEmbed(embedItem(r).MessageEmbed)
		return
	}

	embed := dream.NewEmbed().
		SetTitle(fmt.Sprintf("Search results [%d]", len(res))).
		SetColor(system.StatusNotify).SetFooter(strings.Repeat("_", 156))

	embed.Description = "***Enter an index to select a torrent or type `cancel` to finish***\n\n"

	for i, v := range res {
		embed.Description += fmt.Sprintf("`[%d] (%s / %s): `__%s__ `\n (SE: %d | LE: %d)`\n", i, v.Category, v.Subcategory, v.Title, v.Seeders, v.Leechers)
	}

	msg, err := ctx.ReplyEmbed(embed.MessageEmbed)
	if err != nil {
		ctx.ReplyError("Error sending results")
		return
	}

	for {
		m := ctx.Ses.NextMessageCreate()
		if m.Author.ID != ctx.Msg.Author.ID {
			continue
		}

		if m.Content == "cancel" {
			return
		}
		n, err := strconv.Atoi(m.Content)
		if err != nil {
			return
		}
		ctx.Ses.DG.ChannelMessageDelete(m.ChannelID, m.ID)

		if n < 0 || n >= len(res) {
			ctx.ReplyError("Index out of bounds")
		}
		err = PiratebayInfo(res[n].URL, &res[n])
		if err != nil {
			ctx.ReplyError("Error fetching torrent information")
			return
		}

		_, err = ctx.Ses.DG.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, embedItem(res[n]).MessageEmbed)
		if err != nil {
			ctx.ReplyError(err)
		}
	}
}
