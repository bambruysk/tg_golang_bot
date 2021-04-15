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


func ( gs *  GspreadHoldes) ReadHoldes ()  []Holde {
	
}
