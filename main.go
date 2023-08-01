package main

import (
	"fmt"

	"github.com/October-9th/LibraryOfBabel/util"
)

func main() {

	text := "this is the text where you need to find in the library"
	standardText := util.StandardizationText(text)
	libraryCoordinate := util.CalculateLibraryCoordinate()
	address := util.GetAddress(standardText, int64(libraryCoordinate))
	fmt.Println("hexagon address:", address)
	findText := util.GetContent(address)
	fmt.Println("exact match: ", findText)
}
