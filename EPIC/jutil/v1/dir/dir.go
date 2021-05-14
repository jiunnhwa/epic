package dir

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

type FolderService struct {
	BaseDir,InSubDir,OutSubDir string
	InFile,OutFile string
}

//NewFolderService sets the working dir BaseDir, InDir, OutDir
func NewFolderService(BaseDir, InDir, OutDir string) *FolderService {
	return &FolderService{ BaseDir:BaseDir, InSubDir: filepath.Join(BaseDir,InDir), OutSubDir: filepath.Join(BaseDir,OutDir)}
}

//SetInOutFileName path joins the current working InFile,OutFile, based on fname
func (fs *FolderService )SetInOutFileName (fname string) *FolderService {
	fs.InFile, fs.OutFile  = path.Join(fs.InSubDir,fname),path.Join(fs.OutSubDir,fname)
	return fs
}

//MoveInFileToOutFile archives the infile as outfile
func (fs *FolderService )MoveInFileToOutFile() *FolderService {
	err := os.Rename(fs.InFile, fs.OutFile)
	if err != nil {
		log.Println(err)
	}
	return fs
}