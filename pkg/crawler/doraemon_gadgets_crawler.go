package crawler

import (
	"strings"
	"time"

	"Github.com/Yobubble/go-crawling/pkg/constants"
	"Github.com/Yobubble/go-crawling/pkg/entities"
	"Github.com/Yobubble/go-crawling/pkg/utils"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

type doraemonGadgetsCrawler struct {
	c *colly.Collector // main collector
	dc *colly.Collector // detail collector
}

func (d *doraemonGadgetsCrawler) ScapeGadgetListFromAtoZ() error {
	var i rune
	var err error
	var total int
	var removedGadgets int
	var result []entities.DoraemonGadget
	var filteredResult []entities.DoraemonGadget
	var mergedResult []entities.DoraemonGadget

	for i = 'A'; i <= 'Z'; i++ {
		result, err = d.GetGadgetList(i)
		if err != nil {
				utils.Log.WithError(err).Error("Error getting gadget list for letter " + string(i))
				return err
		}

		filteredResult = utils.RemoveUncertainData(result)

		// Log results for debugging purposes
		utils.Log.WithFields(logrus.Fields{
				"letter": string(i),
				"original_count": len(result),
				"filtered_count": len(filteredResult),
		}).Info("Scraping results")

		mergedResult = append(mergedResult, filteredResult...)

		total += len(result)
		removedGadgets += len(filteredResult)
		
		utils.Log.WithField("Field", string(i)).Info("Finished scraping letter " + string(i))
	}

	utils.JsonSerialize(mergedResult, "data", constants.MainExportPath)

	utils.Log.WithFields(logrus.Fields{
			"original_result": total,
			"filtered_result": removedGadgets,
	}).Info("Finished scraping from A to Z")

	return nil
}

func (d *doraemonGadgetsCrawler) GetGadgetList(letter rune) ([]entities.DoraemonGadget, error) {
	var result []entities.DoraemonGadget

	// visit at gadget's url
	d.c.OnHTML(constants.TargetHtmlElementLink, func(e *colly.HTMLElement) {
			link := e.Attr("href")
			d.dc.Visit(e.Request.AbsoluteURL(link))
	})

	// scape necessary gadget's information
	d.dc.OnHTML(constants.TargetHtmlElementDetails, func(e *colly.HTMLElement) {
			engName := e.ChildText(constants.TargetHtmlElementTitleEn)
			jpName := e.ChildText(constants.TargetHtmlElementTitleJp)
			function := e.ChildText(constants.TargetHtmlElementFunction)
			appearsIn := e.ChildAttrs(constants.TargetHtmlElementAppearsIn, "title")
			imageUrl := e.ChildAttr(constants.TargetHtmlElementImg, "src")
			var filteredAppearsIn []string

			// episodes filtration
			for _, episode := range appearsIn {
					if !strings.HasPrefix(episode, "Special:") && !strings.HasPrefix(episode, "Category:") {
							filteredAppearsIn = append(filteredAppearsIn, episode)
					}
			}
			
			result = append(result, entities.DoraemonGadget{
					EngName:    engName,
					JpName:    jpName,
					Description: function,
					AppearsIn: filteredAppearsIn,
					ImageUrl:  imageUrl,
			})
	})

	// visit main url
	d.c.Visit(constants.MainUrl + string(letter))

	d.c.Wait()
	d.dc.Wait()

	return result, nil
}

func NewDoraemonGadgetsCrawler() *doraemonGadgetsCrawler {
	c := colly.NewCollector(
		colly.AllowedDomains(constants.DoraemonPrimaryUrl, constants.DoraemonSecondaryUrl),
	)
	c.UserAgent = constants.UserAgent
	dc := c.Clone()

	c.OnRequest(func(r *colly.Request) {
			utils.Log.WithField("Visiting", r.URL.String()).Trace("Main Colllector Visiting URL")
	}) 

	dc.OnRequest(func(r *colly.Request) {
		utils.Log.WithField("Visiting", r.URL.String()).Trace("Detail Collector Visiting URL")
	})

	c.OnError(func(r *colly.Response, err error) {
    utils.Log.WithFields(logrus.Fields{
        "URL": r.Request.URL,
        "Error": err,
    }).Error("Request failed")
	})

	dc.OnError(func(r *colly.Response, err error) {
    utils.Log.WithFields(logrus.Fields{
        "URL": r.Request.URL,
        "Error": err,
    }).Error("Request failed")
	})

	c.Limit(&colly.LimitRule{
    DomainGlob:  "*fandom.com/*",
    Parallelism: 2,
    Delay:      5 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	dc.Limit(&colly.LimitRule{
		DomainGlob: "*fandom.com/*",
		Parallelism: 2,
		Delay: 3 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	return &doraemonGadgetsCrawler{
		c: c,
		dc: dc, 
	}
}
