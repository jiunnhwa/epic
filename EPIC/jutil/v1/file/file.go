package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//ReadFile the given filename
func ReadFile(filename string) ([]byte ){
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("ReadFile error for",filename, err)
	}
	//fmt.Print(string(data))
	return data
}

//supercedes filewrite
//WriteFile writes the data as filename, , with create or truncate
func WriteFile(filename string, data string) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	err := ioutil.WriteFile(filename, []byte(data), 0666)
	if err != nil {
		log.Println("WriteFile error for", filename, err)
		//return false, err
	}
}

func RenameFile(oldfilename, newfilename string){
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	err := os.Rename(oldfilename, newfilename)
	if err != nil {
		log.Println("RenameFile error on", oldfilename, newfilename, err)
	}
}

func GetFileInfos(dirname, extension string)   ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println("GetFileInfos on dir", dirname, "with error", err)
		return nil,err
	}
	var filesext []fs.FileInfo
	for _, f := range files {
		if !f.IsDir()  {
			s := strings.Split(f.Name(),".")
			if extension== s[len(s)-1]{
				filesext = append(filesext, f)
			}
		}
	}
	return filesext,nil
}

//GetFileNames returns a list of files with the full path, matching the given file extension
func GetFileNames(dirname, fileext string) []string {
	var filenames []string
	files, _ := GetFileInfos(dirname, fileext)
	for _, f := range files {
		filenames = append(filenames,filepath.Join(dirname, f.Name()))
	}
	return filenames
}

//
////GetFiles returns a list of FileInfo structs, excluding direntries
////Returns filtered list by explicit extension type, not wildcards, if given
//func GetFiles(dirname, extension string) ([]fs.FileInfo, error) {
//	files, err := ioutil.ReadDir(dirname)
//	if err != nil {
//		log.Println(err)
//		return nil,err
//	}
//	var filesext []fs.FileInfo
//	for _, f := range files {
//		if !f.IsDir() {
//			if len(extension) > 0 {
//				s := strings.Split(f.Name(), ".")
//				if extension == s[len(s)-1] {
//					filesext = append(filesext, f)
//				}
//			} else{
//				filesext = append(filesext, f)
//			}
//		}
//	}
//	return filesext,nil
//}


//FileWrite writes the data as filename, , with create or truncate
//returns err
func FileWrite(filename string, data string) error {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	err := ioutil.WriteFile(filename, []byte(data), 0666)
	if err != nil {
		log.Println("FileWrite: cannot write " + filename)
		return err
	}
	return nil
}

//FileAppend safely appends text to the file.
//returns err
func FileAppend(fname, text string) error {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetMaxNextID uses a closure to generate ID sequences
func GetMaxNextID(prefix, dirname, extension string) string {
	files, _ := GetFileInfos(dirname, extension)
	if len(files) ==0 {
		return fmt.Sprintf("%s-%0004d",prefix,1) //Start as ??-0001.txt
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() > files[j].Name() })
	i, _:= strconv.Atoi(strings.Replace(strings.Split(files[0].Name(),".")[0],prefix+"-","",-1))
	i++
	return fmt.Sprintf("%s-%0004d",prefix,i)
}
