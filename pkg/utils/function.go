package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"Github.com/Yobubble/go-crawling/pkg/entities"
)

func JsonSerialize(data interface{}, fileName string, targetPath string) {
	fullFileName := filepath.Join(targetPath, fmt.Sprintf("%s.json", fileName))

	file, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644) 
	if err != nil {
		Log.WithField("text", err).Error()
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") 
	err = encoder.Encode(data)
	if err != nil {
		Log.WithField("text", err).Error()
	}
}

func RemoveUncertainData(result []entities.DoraemonGadget) []entities.DoraemonGadget {
	var filteredResult []entities.DoraemonGadget

	for _, data := range result {
		if data.EngName != "" && data.JpName != "" && data.Description != "" && data.AppearsIn != nil && data.ImageUrl != "" {
			filteredResult = append(filteredResult, data)
		}
	}	
	return filteredResult 

}