package main

import (
	//"bytes"
	"fmt"
	"os"
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
}

func (self *FCSFile) initFCS(path string) {

	//Open the binary FCS file for parsing by
	//using byte offsets.
	f, err := os.Open(path)
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

	pairs := strings.Split(txtContent, string(12))

	self.txtDict = map[string]string{}

	//Construct a dictionary of parameters and their values
	for i := 1; i < len(pairs); i = i + 2 {

		x, y := pairs[i-1], pairs[i]
		self.removeChar(&x) //Take away any $ or spaces
		self.txtDict[x] = y

	}

	fmt.Println(self.txtDict)

}

//Removes $ (replaced with "") and spaces from string (replaced with "_")
func (self *FCSFile) removeChar(s *string) {

	*s = strings.Replace(*s, "$", "", -1)
	*s = strings.Replace(*s, " ", "_", -1)

}

//Reads a partiuclar size of bytes (byteSize) starting at a certain part of the file (f)
// (offset).  Returns a cleaned string value.
func (self *FCSFile) readBytes(f *os.File, byteSize int64, offset int64) string {

	readBytes := make([]byte, byteSize)
	f.ReadAt(readBytes, offset)
	byteValue := strings.TrimSpace(string(readBytes))

	return byteValue

}

func main() {

	newFile := FCSFile{}
	newFile.initFCS("./test.fcs")

}
