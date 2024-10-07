package main

import (
	"Github.com/Yobubble/go-crawling/pkg/constants"
	"Github.com/Yobubble/go-crawling/pkg/crawler"
	"Github.com/Yobubble/go-crawling/pkg/utils"
)

func main() {
	// init logger
	utils.LogInit()

	// init doraemon gadget crawler 
	dora := crawler.NewDoraemonGadgetsCrawler()

	// call function to scrape doraemon gadgets and output in json format
	result, err := dora.ScapeGadgetListFromAtoZ()
	if err != nil {
		utils.Log.WithError(err)
	}

	// generate json file from the result
	utils.JsonSerialize(result, "documents", constants.MainExportPath)
}