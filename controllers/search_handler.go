package controllers

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
)

// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(c *gin.Context) {
	var err error
	// TODO Don't create a new client for each request
	// TODO Fallback to postgres search if es is down

	page := c.Param("page")

	// db params url
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if pagenum <= 0 {
			NotFoundHandler(c)
			return
		}
	}

	searchParam, torrents, nbTorrents, err := search.SearchByQuery(c, pagenum)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Convert back to strings for now.
	// TODO Deprecate fully SearchParam and only use TorrentParam
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	nav := navigation{int(nbTorrents), int(searchParam.Max), int(searchParam.Offset), "search"}
	searchForm := newSearchForm(c)
	searchForm.TorrentParam, searchForm.Category = searchParam, category

	modelList(c, "site/torrents/listing.jet.html", models.TorrentsToJSON(torrents), nav, searchForm)
}