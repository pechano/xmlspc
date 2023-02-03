package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/sqweek/dialog"
)

func main(){
	spcPath :=loadfile() 
	BPF := dialog.Message("Is this a Biocidal Product Family (BPF)?").Title("Single SPC or Family?").YesNo()
	if BPF == true {
		handleBPF(spcPath)
	} else {
	handleSPC(spcPath)
}
}

func handleSPC(spcPath string)(){

spcPath = filepath.Clean(spcPath)
	spcFolder := filepath.Dir(spcPath)
	fmt.Println(spcPath)
	SPC,err := os.Open(spcPath)

	if err != nil {log.Println(err.Error())}
	defer SPC.Close()
	byteValue, _ := ioutil.ReadAll(SPC)
	var baseinfo SPCfile
	err = xml.Unmarshal(byteValue, &baseinfo)
	if err != nil {log.Println(err.Error()+"Is this a BPF SPC?")}

	

	outputFilePath:= filepath.Join(spcFolder,baseinfo.ApplicationName+".txt")




	os.Create(outputFilePath)
	outputFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {log.Println(err.Error())}
	fmt.Fprintln(outputFile,"BP/BPF navn:")
	fmt.Fprintln(outputFile,baseinfo.ApplicationName)

	fmt.Fprintln(outputFile,"Godkendte anvendelser:")
	for _, useslice := range baseinfo.Uses.AuthorisedUse{

	fmt.Fprintln(outputFile,useslice.Name)
		for _, fieldslice := range useslice.Fields{
fmt.Fprintln(outputFile,fieldslice)
		}
	}
}


func handleBPF(spcPath string)(){

spcPath = filepath.Clean(spcPath)
	spcFolder := filepath.Dir(spcPath)
	fmt.Println(spcPath)
	SPC,err := os.Open(spcPath)

	if err != nil {log.Println(err.Error())}
	defer SPC.Close()
	byteValue, _ := ioutil.ReadAll(SPC)
	var baseinfo BPFfile
	err = xml.Unmarshal(byteValue, &baseinfo)
	if err != nil {log.Println(err.Error()+"Is this a single product SPC?")}

	

	outputFilePath:= filepath.Join(spcFolder,baseinfo.BPFname+".txt")




	os.Create(outputFilePath)
	outputFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {log.Println(err.Error())}
	fmt.Fprintln(outputFile,"BP/BPF navn:")
	fmt.Fprintln(outputFile,baseinfo.BPFname)

	fmt.Fprintln(outputFile,"Produkt(er):")
	for _, metaSlice:= range baseinfo.MetaSPCs.MetaSPC{
		for _, productSlice := range metaSlice.Products{
			fmt.Fprintln(outputFile, productSlice.Tradename)}
	}
	fmt.Fprintln(outputFile,"Anvendelser fra SPC:")
	for _, metaSlice := range baseinfo.MetaSPCs.MetaSPC{
		fmt.Fprintln(outputFile,"MetaSPC:"+metaSlice.Name)
		for _, useSlice := range metaSlice.Uses{

	fmt.Fprintln(outputFile,"Anvendelser:")
		fmt.Fprintln(outputFile, useSlice.Name)
		fmt.Fprintln(outputFile, useSlice.Field)
		}
	}

}

func loadfile()(filename string ){
	filename, err := dialog.File().Filter("SPC", "xml").Load()	
	if err != nil {log.Println(err.Error())}
	filename = filepath.Clean(filename)
	return
}

//SPC STRUCTS**************************

type SPCfile struct{
	XMLName xml.Name `xml:"SPC"`
	ApplicationName string `xml:"name,attr"`
	Uses SPCuses `xml:"AuthorisedUses"`
}


type SPCuses struct{

	XMLName xml.Name `xml:"AuthorisedUses"`
	AuthorisedUse []SPCuse `xml:"AuthorisedUse"`
}

type SPCuse struct {

	XMLName xml.Name `xml:"AuthorisedUse"`
	Name string `xml:"Name"`
	Fields []string `xml:"Fields>Field"`

}









//BPF STRUCTS****************************************


type BPFfile struct{
	XMLName xml.Name `xml:"SPFBC"`
	BAS string `xml:"ProductInfo>Composition>ActiveSubstance>BAS"`
	BPFname string `xml:"name,attr"`
	MetaSPCs metaspcs `xml:"MetaSPCs"`
}



type metaspcs struct {
	XMLName xml.Name `xml:"MetaSPCs"`
	MetaSPC []MetaSPC `xml:"MetaSPC"`
}

type MetaSPC struct{

	XMLName xml.Name `xml:"MetaSPC"`
	Function string `xml:"ProductInfo>Composition>ActiveSubstance>Function"`
	Products []product `xml:"Products"`
	Uses []use `xml:"AuthorisedUses"`
	Name string `xml:"name,attr"`
}

type product struct {
	XMLName xml.Name `xml:"Products"`
	Tradename string `xml:"Product>AdminInfo>TradeNames>TradeName>Text"` 
}

type use struct {
	XMLName xml.Name `xml:"AuthorisedUses"`
	Name string `xml:"AuthorisedUse>Name"` 
	Field string `xml:"AuthorisedUse>Fields>Field"` 
}
