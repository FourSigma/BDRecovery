package main

import (
	//"bytes"
	"fmt"
	"os"
	"time"
	//"path/filepath"
	"strconv"
	"strings"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)

	}
}

type FCSFile struct {
	version  string
	txtStart int
	txtEnd   int
	txtDict  map[string]string
	f        *os.File
}

func (self *FCSFile) InitFCS(path string) {

	//Open the binary FCS file for parsing by
	//using byte offsets.
	f, err := os.Open(path)
	self.f = f
	self.readTextSegment(f) //Populates txtDict with paramters from TEXT segment.
	check(err)
	defer f.Close()

}

//Reads the TEXT segment of the FCS binary and creates
//a dictionary map of the key-value pairs in that
//segment
func (self *FCSFile) readTextSegment(f *os.File) {

	//Offsets based on FCS specs
	self.version = self.readBytes(f, 6, 0)
	tmp := self.readBytes(f, 8, 10)
	self.txtStart, _ = strconv.Atoi(tmp)
	tmp = self.readBytes(f, 8, 18)
	self.txtEnd, _ = strconv.Atoi(tmp)

	//Size of the TEXT segment in the FCS file
	txtSize := self.txtEnd - self.txtStart

	//Stores the content of the TEXT Segment after reading
	txtContent := self.readBytes(f, int64(txtSize), int64(self.txtStart))

	//Data from TEXT segment contained in continous array
	pairs := strings.Split(txtContent, string(12))

	self.txtDict = map[string]string{}

	//Construct a dictionary of parameters and their values
	for i := 1; i < len(pairs); i = i + 2 {

		x, y := pairs[i-1], pairs[i]
		self.cleanString(&x, true)  //Take away any $ or spaces from keys
		self.cleanString(&y, false) //Trims spaces from values
		self.txtDict[x] = y

	}

	//z, _ := filepath.Glob("./*.fcs")
	//for k, v := range self.txtDict {

	//	fmt.Println("Key: " + k)
	//	fmt.Println(v)
	//}

	const shortForm = "2006-Jan-02"
	t, _ := time.Parse(shortForm, "2013-FEB-03")
	fmt.Println(t)

}

//Removes $ (replaced with "") and spaces from string (replaced with "_") for
//only keys (key == true). All strings are trimed
func (self *FCSFile) cleanString(s *string, key bool) {

	if key == true {
		*s = strings.Replace(*s, "$", "", -1)
		*s = strings.Replace(*s, " ", "_", -1)
	}

	*s = strings.TrimSpace(*s) //Trims whitespace

}

//Reads a particular size of bytes (byteSize) starting at a certain part of the file (f)
// (offset).  Returns a cleaned string value.
func (self *FCSFile) readBytes(f *os.File, byteSize int64, offset int64) string {

	readBytes := make([]byte, byteSize)
	f.ReadAt(readBytes, offset)
	byteValue := strings.TrimSpace(string(readBytes)) //Bytes into string conversion

	return byteValue

}

/*****************************************************************************
**   This is the END of the FCSFile defintion and methods.					**
******************************************************************************/

type FCSInfo struct {
	oldFN   string //Numeric file names ex. 10203030202302.fcs
	newFN   string //New Filename ex. EXP_Name_
	srcPath string //Source Path - This is where the BDData file is located
	desPath string //Destination Path - Where the recovered files will be placed
	expName string //Name is experiment as read from TEXT segment of FCS
	expDate string //Date of experiment as read from TEXT segment of FCS
	expSrc  string //Specimen name as read from TEXT segment of FCS

}

func (self *FCSInfo) InitFCSInfo(fcs *FCSFile) {
	self.expName = fcs.txtDict["EXPERIMENT_NAME"]
	self.expDate = fcs.txtDict["DATE"]
	self.expSrc = fcs.txtDict["SRC"]
	self.newFN = fcs.txtDict["FIL"]
	self.oldFN = fcs.f.Name()

}
func (self *FCSInfo) SetPath(src string, des string) {
	self.srcPath = src
	self.desPath = des

}

func main() {

	newFile := &FCSFile{}
	newFile.InitFCS("test.fcs")
	fileInfo := FCSInfo{}
	fileInfo.InitFCSInfo(newFile)

}
