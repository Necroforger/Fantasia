package general

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Necroforger/dream"

	"github.com/Necroforger/Fantasia/system"
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
func PiratebayScrape(query string) ([]PiratebayItem, error) {
	items := []PiratebayItem{}

	searchURL := func(query string) string {
		return fmt.Sprintf("https://thepiratebay.org/search/%s/0/99/0", url.QueryEscape(query))
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
						item.URL = "https://thepiratebay.se" + val
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

	res, err := PiratebayScrape(ctx.Args.After())
	if err != nil {
		ctx.ReplyError(err)
	}

	cut := func(str string, limit int) string {
		if len(str) > limit {
			return str[:limit] + "..."
		}
		return str
	}

	if len(res) > 0 {
		r := res[0]

		PiratebayInfo(r.URL, &r)
		ctx.ReplyEmbed(dream.NewEmbed().
			SetTitle(r.Title).
			SetURL(r.URL).
			SetDescription("[<:iconmagnet:342227624770142208>magnet]("+r.Magnet+")").
			AddField("Uploader", "`"+r.Uploader+"`").
			AddField("Seeders", "`"+strconv.Itoa(r.Seeders)+"`").
			AddField("Leechers", "`"+strconv.Itoa(r.Leechers)+"`").
			AddField("Size", "`"+r.Size+"`").
			AddField("Uploaded", "`"+r.Uploaded+"`").
			AddField("Description", "```"+cut(r.Description, 1024)+"```").
			SetColor(system.StatusNotify).
			InlineAllFields().
			MessageEmbed)
		return
	}

	ctx.ReplyError("No results found")
}
