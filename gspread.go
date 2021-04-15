package main

import (
	"log"
	"strconv"
	"time"

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
		case "Locations" : {
			for i:= 1; i < len(row); i++ {
				settings.Locations = append(settings.Locations, row[i].Value)
			} 
		}
		
		} 
		if err !=  nil {
			panic(err)
		}
	}

	return settings

}


func ( gs *  GspreadHoldes) ReadHoldes ()  []Holde {
	spread, err := gs.service.FetchSpreadsheet(gs.spreadsheetName)
	if err !=  nil {
		panic(err)
	}
	sheet, err := spread.SheetByTitle("Holdes")
	if err !=  nil {
		panic(err)
	}
	var result [] Holde 
	// skip first row as header
	for i:= 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		id, err := strconv.Atoi( row[0].Value)
		if err != nil {
			log.Println("Wrong id in ", i)
			continue
		}
		name := row[1].Value
		result = append(result, Holde{
			Name:      name,
			ID:        id,
			Amount:    0,
			Level:     1,
			Owner:     "",
			LastVisit: time.Time{},
		})

	} 
	
	
	return result
}
