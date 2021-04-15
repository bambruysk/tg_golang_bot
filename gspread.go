package main

import (
	"gopkg.in/Iwark/spreadsheet.v2"
)

type GspreadHoldes {
	service * spreadsheet.Service
	spreadsheetName string
}

func NewGspreadHoldes() GspreadHoldes{

	spreadsheetName := "1tODyRkeelf-YXnpiGeOH9v71LCEV9epw9IzSLi3mv2s"

	service, err := spreadsheet.NewService()

	if err != nil {
		panic(err)
	}

	return GspreadHoldes{
		service : service
		spreadsheetName : spreadsheetName
	} 
}

func (gs * GspreadHoldes) ReadSettings () HoldeGameSettings {

	spread, err := gs.service.FetchSpreadsheet(gs.spreadsheetName)
	if err !=  nil {
		panic(err)
	}
	sheet, err := spread.SheetByTitle("Settings")
	if err !=  nil {
		panic(err)
	}
	settings :=  HoldeGameSettings{}
	for _, row :=  range sheet.Rows{
		if row[0].Value == "" {
			continue
		}
		switch val := row[0].Value {
		case "HoldeNums" : settings.HoldeNums = strconv.Atoi(val)
		case "WorldSizeX" : settings.WorldSizeX = strconv.Atoi(val)
		case "WorldSizeY" : settings.WorldSizeY = strconv.Atoi(val)
		case "MoneyPerHour" : settings.MoneyPerHour = strconv.Atoi(val)
		} 
	}


}


func ( gs *  GspreadHoldes) ReadHoldes ()  []Holde {

}
