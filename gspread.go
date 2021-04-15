package main

import (
	"strconv"

	"gopkg.in/Iwark/spreadsheet.v2"
)

type GspreadHoldes struct {
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
		service : service,
		spreadsheetName : spreadsheetName,
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
		key,val := row[0].Value,row[1].Value
		switch key {
		case "HoldeNums" : settings.HoldeNums, err = strconv.Atoi(val)
		case "WorldSizeX" : settings.WorldSizeX, err = strconv.Atoi(val)
		case "WorldSizeY" : settings.WorldSizeY, err = strconv.Atoi(val)
		case "MoneyPerHour" : settings.MoneyPerHour, err = strconv.ParseFloat(val,64)
		case "TimeDegradation" : settings.TimeDegradation, err = strconv.ParseFloat(val,64)
		
		} 
		if err !=  nil {
			panic(err)
		}
	}

	return settings

}


func ( gs *  GspreadHoldes) ReadHoldes ()  []Holde {
	return []Holde{}
}
