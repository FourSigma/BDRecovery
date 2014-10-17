/* The MIT License (MIT)

Copyright (c) 2014  Siva Manivannan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WITHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/****USAGE****
Usage:  go build recovery.go
Command Line:  recovery -src <BDData Dir> -des <Backup Dir>
Example in MacOS:  recovery -src /Users/JDoe/BDdata -des /Users/JDoe/RecoveredFCS
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const SEP = string(os.PathSeparator)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)

	}
}

//Converts date into a more convient format
func convertDate(date *string) {

	const shortForm = "02-Jan-2006"
	t, _ := time.Parse(shortForm, *date)
	timeString := t.String()
	*date = strings.Split(timeString, " ")[0]

}

//Copies from a source file to a new files (des)
func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
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
	oldFN    string //Numeric file names ex. 10203030202302.fcs
	newFN    string //New Filename ex. EXP_Name_
	expName  string //Name is experiment as read from TEXT segment of FCS
	expDate  string //Date of experiment as read from TEXT segment of FCS
	expSrc   string //Specimen name as read from TEXT segment of FCS
	expTube  string //Experimental Tube
	expUser  string //Export username (person who conducted the experiment)
	filePath string //Where the file should be located
}

func (self *FCSInfo) InitFCSInfo(fcs *FCSFile) {

	self.expName = fcs.txtDict["EXPERIMENT_NAME"]
	self.expTube = fcs.txtDict["TUBE_NAME"]
	self.oldFN = fcs.f.Name()
	self.expSrc = fcs.txtDict["SRC"]
	self.expUser = fcs.txtDict["EXPORT_USER_NAME"]

	self.expDate = fcs.txtDict["DATE"]
	convertDate(&self.expDate) //Coverts date to a better string format

	self.newFN = self.expName + "_" + self.expSrc + "_" + self.expTube + ".fcs"
	self.cleanName(&self.newFN, true)

	self.filePath = SEP + self.expUser + SEP + self.expDate + " " + self.expName + SEP + self.expSrc
	self.expName = self.expDate + " " + self.expName

}

//Cleans file names of "/" and "\" characters that might
//interfer with output.
func (self *FCSInfo) cleanName(s *string, isFile bool) {

	if isFile == true {
		*s = strings.Replace(*s, "/", "-", -1)
		*s = strings.Replace(*s, "\\", "-", -1)
	}

}

/*****************************************************************************
**   This is the END of the FCSInfo defintion and methods.					**
******************************************************************************/
type Path struct {
	srcPath string //Source Path - This is where the BDData file is located
	desPath string //Destination Path - Where the recovered files will be placed
}

//Set the path of the BDData directory and the destiantion  of the recovered files.
func (self *Path) SetPath(src string, des string) {
	self.srcPath = src
	self.desPath = des
}

//Reads the the names of all *.fcs files and puts them in
//a slice and returns the slice.
func (self *Path) GlobIt() []string {
	os.Chdir(self.srcPath)
	f, err := filepath.Glob("*.fcs")

	check(err)

	return f

}

//Copies files and moves them to the desination directory.
func (self *Path) RenameMove(fcsInfo *FCSInfo) {
	os.MkdirAll(self.desPath+fcsInfo.filePath, 0777)
	cwd, _ := os.Getwd()
	err := cp(filepath.Join(cwd, fcsInfo.oldFN), filepath.Join(self.desPath, fcsInfo.filePath, fcsInfo.newFN))
	if err == nil {
		fmt.Println(fcsInfo.oldFN + "   ------>" + fcsInfo.newFN)
	}
}

/*****************************************************************************
**   This is the END of the Path
 defintion and methods.					**
******************************************************************************/

func main() {

	var src = flag.String("src", "", "Location of BDData Directory")
	var des = flag.String("des", "", "Location where recoverd files will be stored")
	flag.Parse()

	paths := &Path{}
	paths.SetPath(*src, *des)
	files := paths.GlobIt()

	newFile := &FCSFile{}
	fileInfo := &FCSInfo{}

	for _, fileName := range files {

		newFile.InitFCS(fileName)
		fileInfo.InitFCSInfo(newFile)
		paths.RenameMove(fileInfo)
	}

}
