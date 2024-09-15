package main

import (
	"fmt"

	"Github.com/Yobubble/go-crawling/pkg/crawler"
	"Github.com/Yobubble/go-crawling/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	utils.LogInit()

	dorae := crawler.NewDoraemonGadgetsCrawler()

	// loop A - Z 
	var i rune
	for i = 'A' ; i < 'Z' ; i++ {
		result, err := dorae.GetGadgetList(i)
		if err != nil {
			utils.Log.WithField("text", err).Error()
		}

		filteredResult := utils.RemoveUncertainData(result)
		for _, data := range filteredResult {
			utils.Log.WithFields(logrus.Fields{
				"en_name" : data.EngName,
				"jp_name" : data.JpName,
				"function" : data.Function,
				"appears_in" : data.AppearsIn,
				"image_url" : data.ImageUrl,
			}).Info()
		}


		utils.Log.WithField(fmt.Sprintf("data deleted at %c", i), len(result) - len(filteredResult)).Debug()
		utils.JsonSerialize(filteredResult, i, "../data")
	}  
}