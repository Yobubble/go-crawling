package crawler

import (
	"Github.com/Yobubble/go-crawling/pkg/constants"
	"Github.com/Yobubble/go-crawling/pkg/entities"
	"Github.com/Yobubble/go-crawling/pkg/utils"
	"github.com/gocolly/colly/v2"
)

type doraemonGadgetsCrawler struct {
	c *colly.Collector
}

func (d *doraemonGadgetsCrawler) GetGadgetList(letter rune) ([]entities.DoraemonGadget, error) {
	var result []entities.DoraemonGadget

	d.c.OnHTML(constants.TargetHtmlElementLink, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		d.c.Visit(e.Request.AbsoluteURL(link))
	})

	d.c.OnHTML(constants.TargetHtmlElementDetails, func(e *colly.HTMLElement) {
		result = append(result, entities.DoraemonGadget{
			EngName: e.ChildText(constants.TargetHtmlElementTitleEn),
			JpName: e.ChildText(constants.TargetHtmlElementTitleJp),
			Function: e.ChildText(constants.TargetHtmlElementFunction),
			AppearsIn: e.ChildAttrs(constants.TargetHtmlElementAppearsIn, "title"),
			ImageUrl: e.ChildAttr(constants.TargetHtmlElementImg, "src"), 
		})
	})

	d.c.OnError(func(_ *colly.Response, err error) {
		utils.Log.WithField("Error", err).Error()
	})

	d.c.OnRequest(func(r *colly.Request) {
		utils.Log.WithField("Visiting", r.URL.String())
	})

	d.c.Visit(constants.MainUrl + string(letter))

	return result, nil
} 

func NewDoraemonGadgetsCrawler() *doraemonGadgetsCrawler {
	c := colly.NewCollector(
		colly.AllowedDomains(constants.DoraemonPrimaryUrl, constants.DoraemonSecondaryUrl),
	)
	c.UserAgent = constants.UserAgent 
	return &doraemonGadgetsCrawler{
		c: c,
	}
}

