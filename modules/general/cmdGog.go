package general

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	humanize "github.com/dustin/go-humanize"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/Fantasia/util"
)

// GogResult is the json returned from querying the GOG search endpoint
type GogResult struct {
	Products         []GogProduct `json:"products"`
	Ts               interface{}  `json:"ts"`
	Page             int          `json:"page"`
	TotalPages       int          `json:"totalPages"`
	TotalResults     string       `json:"totalResults"`
	TotalGamesFound  int          `json:"totalGamesFound"`
	TotalMoviesFound int          `json:"totalMoviesFound"`
}

// GogProduct ...
type GogProduct struct {
	CustomAttributes []interface{} `json:"customAttributes"`
	Developer        string        `json:"developer"`
	Publisher        string        `json:"publisher"`
	Price            struct {
		Amount                     string `json:"amount"`
		BaseAmount                 string `json:"baseAmount"`
		FinalAmount                string `json:"finalAmount"`
		IsDiscounted               bool   `json:"isDiscounted"`
		DiscountPercentage         int    `json:"discountPercentage"`
		DiscountDifference         string `json:"discountDifference"`
		Symbol                     string `json:"symbol"`
		IsFree                     bool   `json:"isFree"`
		Discount                   int    `json:"discount"`
		IsBonusStoreCreditIncluded bool   `json:"isBonusStoreCreditIncluded"`
		BonusStoreCreditAmount     string `json:"bonusStoreCreditAmount"`
	} `json:"price"`
	IsDiscounted    bool `json:"isDiscounted"`
	IsInDevelopment bool `json:"isInDevelopment"`
	ID              int  `json:"id"`
	ReleaseDate     int  `json:"releaseDate"`
	Availability    struct {
		IsAvailable          bool `json:"isAvailable"`
		IsAvailableInAccount bool `json:"isAvailableInAccount"`
	} `json:"availability"`
	SalesVisibility struct {
		IsActive   bool `json:"isActive"`
		FromObject struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"fromObject"`
		From     int `json:"from"`
		ToObject struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"toObject"`
		To int `json:"to"`
	} `json:"salesVisibility"`
	Buyable    bool   `json:"buyable"`
	Title      string `json:"title"`
	Image      string `json:"image"`
	URL        string `json:"url"`
	SupportURL string `json:"supportUrl"`
	ForumURL   string `json:"forumUrl"`
	WorksOn    struct {
		Windows bool `json:"Windows"`
		Mac     bool `json:"Mac"`
		Linux   bool `json:"Linux"`
	} `json:"worksOn"`
	Category         string `json:"category"`
	OriginalCategory string `json:"originalCategory"`
	Rating           int    `json:"rating"`
	Type             int    `json:"type"`
	IsComingSoon     bool   `json:"isComingSoon"`
	IsPriceVisible   bool   `json:"isPriceVisible"`
	IsMovie          bool   `json:"isMovie"`
	IsGame           bool   `json:"isGame"`
	Slug             string `json:"slug"`
}

// GogResultExpanded ...
type GogResultExpanded struct {
	ID                         int    `json:"id"`
	Title                      string `json:"title"`
	PurchaseLink               string `json:"purchase_link"`
	Slug                       string `json:"slug"`
	ContentSystemCompatibility struct {
		Windows bool `json:"windows"`
		Osx     bool `json:"osx"`
		Linux   bool `json:"linux"`
	} `json:"content_system_compatibility"`
	Languages struct {
		En string `json:"en"`
		Jp string `json:"jp"`
	} `json:"languages"`
	Links struct {
		PurchaseLink string `json:"purchase_link"`
		ProductCard  string `json:"product_card"`
		Support      string `json:"support"`
		Forum        string `json:"forum"`
	} `json:"links"`
	InDevelopment struct {
		Active bool        `json:"active"`
		Until  interface{} `json:"until"`
	} `json:"in_development"`
	IsSecret    bool   `json:"is_secret"`
	GameType    string `json:"game_type"`
	IsPreOrder  bool   `json:"is_pre_order"`
	ReleaseDate string `json:"release_date"`
	Images      struct {
		Background          string `json:"background"`
		Logo                string `json:"logo"`
		Logo2X              string `json:"logo2x"`
		Icon                string `json:"icon"`
		SidebarIcon         string `json:"sidebarIcon"`
		SidebarIcon2X       string `json:"sidebarIcon2x"`
		MenuNotificationAv  string `json:"menuNotificationAv"`
		MenuNotificationAv2 string `json:"menuNotificationAv2"`
	} `json:"images"`
	Downloads struct {
		Installers []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Os           string `json:"os"`
			Language     string `json:"language"`
			LanguageFull string `json:"language_full"`
			Version      string `json:"version"`
			TotalSize    int64  `json:"total_size"`
			Files        []struct {
				ID       string `json:"id"`
				Size     int    `json:"size"`
				Downlink string `json:"downlink"`
			} `json:"files"`
		} `json:"installers"`
		Patches       []interface{} `json:"patches"`
		LanguagePacks []interface{} `json:"language_packs"`
		BonusContent  []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Count     int    `json:"count"`
			TotalSize int    `json:"total_size"`
			Files     []struct {
				ID       int    `json:"id"`
				Size     int    `json:"size"`
				Downlink string `json:"downlink"`
			} `json:"files"`
		} `json:"bonus_content"`
	} `json:"downloads"`
	ExpandedDlcs []struct {
		ID                         int    `json:"id"`
		Title                      string `json:"title"`
		PurchaseLink               string `json:"purchase_link"`
		Slug                       string `json:"slug"`
		ContentSystemCompatibility struct {
			Windows bool `json:"windows"`
			Osx     bool `json:"osx"`
			Linux   bool `json:"linux"`
		} `json:"content_system_compatibility"`
		Languages struct {
			En string `json:"en"`
			Jp string `json:"jp"`
		} `json:"languages"`
		Links struct {
			PurchaseLink string `json:"purchase_link"`
			ProductCard  string `json:"product_card"`
			Support      string `json:"support"`
			Forum        string `json:"forum"`
		} `json:"links"`
		InDevelopment struct {
			Active bool        `json:"active"`
			Until  interface{} `json:"until"`
		} `json:"in_development"`
		IsSecret    bool   `json:"is_secret"`
		GameType    string `json:"game_type"`
		IsPreOrder  bool   `json:"is_pre_order"`
		ReleaseDate string `json:"release_date"`
		Images      struct {
			Background          string `json:"background"`
			Logo                string `json:"logo"`
			Logo2X              string `json:"logo2x"`
			Icon                string `json:"icon"`
			SidebarIcon         string `json:"sidebarIcon"`
			SidebarIcon2X       string `json:"sidebarIcon2x"`
			MenuNotificationAv  string `json:"menuNotificationAv"`
			MenuNotificationAv2 string `json:"menuNotificationAv2"`
		} `json:"images"`
		Downloads struct {
			Installers []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				Os           string `json:"os"`
				Language     string `json:"language"`
				LanguageFull string `json:"language_full"`
				Version      string `json:"version"`
				TotalSize    int    `json:"total_size"`
				Files        []struct {
					ID       string `json:"id"`
					Size     int    `json:"size"`
					Downlink string `json:"downlink"`
				} `json:"files"`
			} `json:"installers"`
			Patches       []interface{} `json:"patches"`
			LanguagePacks []interface{} `json:"language_packs"`
			BonusContent  []interface{} `json:"bonus_content"`
		} `json:"downloads"`
	} `json:"expanded_dlcs"`
	Description struct {
		Lead             string `json:"lead"`
		Full             string `json:"full"`
		WhatsCoolAboutIt string `json:"whats_cool_about_it"`
	} `json:"description"`
	Screenshots []struct {
		ImageID              string `json:"image_id"`
		FormatterTemplateURL string `json:"formatter_template_url"`
		FormattedImages      []struct {
			FormatterName string `json:"formatter_name"`
			ImageURL      string `json:"image_url"`
		} `json:"formatted_images"`
	} `json:"screenshots"`
	Videos []struct {
		VideoURL     string `json:"video_url"`
		ThumbnailURL string `json:"thumbnail_url"`
		Provider     string `json:"provider"`
	} `json:"videos"`
	RelatedProducts []struct {
		ID                         int    `json:"id"`
		Title                      string `json:"title"`
		PurchaseLink               string `json:"purchase_link"`
		Slug                       string `json:"slug"`
		ContentSystemCompatibility struct {
			Windows bool `json:"windows"`
			Osx     bool `json:"osx"`
			Linux   bool `json:"linux"`
		} `json:"content_system_compatibility"`
		Languages struct {
			En string `json:"en"`
			Jp string `json:"jp"`
		} `json:"languages"`
		Links struct {
			PurchaseLink string `json:"purchase_link"`
			ProductCard  string `json:"product_card"`
			Support      string `json:"support"`
			Forum        string `json:"forum"`
		} `json:"links"`
		InDevelopment struct {
			Active bool        `json:"active"`
			Until  interface{} `json:"until"`
		} `json:"in_development"`
		IsSecret    bool   `json:"is_secret"`
		GameType    string `json:"game_type"`
		IsPreOrder  bool   `json:"is_pre_order"`
		ReleaseDate string `json:"release_date"`
		Images      struct {
			Background          string `json:"background"`
			Logo                string `json:"logo"`
			Logo2X              string `json:"logo2x"`
			Icon                string `json:"icon"`
			SidebarIcon         string `json:"sidebarIcon"`
			SidebarIcon2X       string `json:"sidebarIcon2x"`
			MenuNotificationAv  string `json:"menuNotificationAv"`
			MenuNotificationAv2 string `json:"menuNotificationAv2"`
		} `json:"images"`
	} `json:"related_products"`
	Changelog string `json:"changelog"`
}

// CmdGog searches gog for the given query
func CmdGog(ctx *system.Context) {
	if ctx.Args.After() == "" {
		ctx.ReplyError("Please enter a search query")
		return
	}

	// =============== Helper functions ===================

	// getUnmarshal query's the given URL and attempts to unmarshal its
	// json response.
	getUnmarshal := func(URL string, data interface{}) (err error) {
		res, err := http.Get(URL)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(data)
		if err != nil {
			return err
		}
		return nil
	}

	getPrice := func(p *GogProduct) string {
		var pricing string
		if p.Price.IsFree {
			pricing = "FREE"
		} else {
			if p.IsDiscounted {
				pricing = "~~$" + p.Price.BaseAmount + "~~" + " | $" + p.Price.FinalAmount
			} else {
				pricing = "$" + p.Price.FinalAmount
			}
		}
		return pricing
	}

	cutstring := func(txt string, maxlen int) string {
		if len(txt) > maxlen {
			return txt[:maxlen] + "..."
		}
		return txt
	}
	// ==========================================================

	var data GogResult
	err := getUnmarshal(fmt.Sprintf("https://www.gog.com/games/ajax/filtered?limit=10&search=%s", url.QueryEscape(ctx.Args.After())), &data)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	if data.Products == nil {
		ctx.ReplyError("Error obtaining gog information")
		return
	}
	if len(data.Products) == 0 {
		ctx.ReplyWarning("No results found")
		return
	}

	var text string
	for i, p := range data.Products {
		text += fmt.Sprintf("[%d] [%s](%s) %s", i, p.Title, "http://gog.com/"+p.URL, getPrice(&p)) + "\n"
	}

	var (
		n   int
		msg *discordgo.Message
	)
	// If there is one result, select the item automatically.
	if len(data.Products) == 1 {
		n = 0
	} else {
		msg, err = ctx.ReplyEmbed(
			dream.NewEmbed().
				SetTitle(fmt.Sprintf("Gog search, [%d] results", len(data.Products))).
				SetColor(system.StatusNotify).
				SetDescription("Type an index to select a search result or `cancel` to cancel the search\n\n" + text).
				MessageEmbed,
		)
		if err != nil {
			return
		}
		m, err := util.RequestMessage(ctx.Ses, ctx.Msg.Author.ID, time.Minute*5)
		if err != nil {
			return
		}
		if index, err := strconv.Atoi(m.Content); err == nil {
			n = index
		} else {
			return
		}
		if n < 0 || n >= len(data.Products) {
			return
		}
		ctx.Ses.DG.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	p := data.Products[n]
	embed := dream.NewEmbed().SetTitle(p.Title)

	var d GogResultExpanded
	err = getUnmarshal(fmt.Sprintf("http://api.gog.com/products/%d?expand=downloads,expanded_dlcs,description,screenshots,videos,related_products,changelog", p.ID), &d)
	if err != nil {
		ctx.ReplyError(err)
		return
	}

	var imageurl string

	if strings.HasPrefix(p.Image, "//") {
		imageurl = "http://" + p.Image[2:] + ".jpg"
	}
	embed.SetThumbnail(imageurl)
	embed.SetURL("http://gog.com/" + p.URL)
	embed.AddField("Price", getPrice(&p))

	if p.Price.IsDiscounted {
		embed.AddField("Discount", fmt.Sprintf("%d%% off", p.Price.DiscountPercentage))
	}
	if p.Category != "" {
		embed.AddField("Category", p.Category)
	}
	if d.GameType != "" {
		embed.AddField("Type", d.GameType)
	}
	if d.ReleaseDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05+0200", d.ReleaseDate)
		if err == nil {
			embed.AddField("Release Date", t.Format("Jan 2 2006")+" ("+humanize.Time(t)+") ")
		} else {
			ctx.ReplyError(err)
		}
	}

	var (
		stars, frac float64
		rating      string
	)
	// Prevent division by zero.
	if p.Rating != 0 {
		stars, frac = math.Modf(float64(p.Rating) / 10.0)
		rating = strings.Repeat("⭐", int(stars))
		if frac >= 0.5 {
			rating += "🔸"
		}
	}
	// Display unfilled stars. Stars + frac + 0.5 rounds upwards then subtracts from five
	// Rounding upwards is necessary because half-stars consume a slot.
	rating += strings.Repeat("▪", int(5-(int(stars+frac+0.5))))
	embed.AddField("Rating", rating)

	// Description
	embed.AddField("Description", "```\n"+cutstring(regexp.MustCompile(`\<.*?\>`).ReplaceAllString(d.Description.Full, ""), 500)+"\n```")
	embed.InlineAllFields()

	// If a selection menu has been created, edit it to contain the expanded information,
	// Else send a new message.
	if msg != nil {
		_, err = ctx.Ses.DG.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, embed.MessageEmbed)
		if err != nil {
			ctx.ReplyError(err)
		}
	} else {
		ctx.Ses.DG.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed.MessageEmbed)
	}
}
