package lib

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// File is descriptor
type File struct {
	PhysicalFile *os.File
	Name         string
	Index        int
	FilePath     string
	PathExt      string
	MimeType     string
	Written      int64
}

// Uploader ...
type Uploader interface {
	Check(r *http.Request, w http.ResponseWriter) error
	Get(r *http.Request, w http.ResponseWriter, name string) (File, error)
	GetSlice(r *http.Request, w http.ResponseWriter, name string) ([]File, error)
	GetAll(r *http.Request, w http.ResponseWriter) ([]File, error)
}

// FileProvider ...
type FileProvider struct {
	MaxSize    int64
	UploadPath string
}

// NewFileProvider ...
func NewFileProvider() *FileProvider {
	return &FileProvider{
		MaxSize: (1 << 10) * 24,
	}
}

// Check ...
func (f *FileProvider) Check(r *http.Request, w http.ResponseWriter) error {
	return r.ParseMultipartForm(f.MaxSize)
}

// makeFile ...
func (f *FileProvider) makeFile(hdr *multipart.FileHeader, name string, index int) (File, error) {
	var ret File
	var err error
	var infile multipart.File
	if infile, err = hdr.Open(); nil != err {
		return ret, err
	}
	defer infile.Close()

	var tmp *os.File
	if tmp, err = ioutil.TempFile(f.UploadPath, "up-"); nil != err {
		return ret, err
	}
	defer tmp.Close()

	ret = File{
		PhysicalFile: tmp,
		Name:         name,
		Index:        index,
		FilePath:     hdr.Filename,
		PathExt:      filepath.Ext(hdr.Filename),
		MimeType:     "",
	}

	var written int64
	if written, err = io.Copy(ret.PhysicalFile, infile); nil != err {
		return ret, err
	}
	ret.Written = written
	return ret, nil
}

// Get ...
func (f *FileProvider) Get(r *http.Request, w http.ResponseWriter, name string) (File, error) {
	var ret File
	var err error
	if fheaders, ok := r.MultipartForm.File[name]; ok {
		for index, hdr := range fheaders {
			ret, err = f.makeFile(hdr, name, index)
			break
		}
	}
	return ret, err
}

// GetSlice ...
func (f *FileProvider) GetSlice(r *http.Request, w http.ResponseWriter, name string) ([]File, error) {
	var ret []File
	if fheaders, ok := r.MultipartForm.File[name]; ok {
		for index, hdr := range fheaders {
			file, err := f.makeFile(hdr, name, index)
			if err != nil {
				return ret, err
			}
			ret = append(ret, file)
		}
	}
	return ret, nil
}

// GetAll ...
func (f *FileProvider) GetAll(r *http.Request, w http.ResponseWriter) ([]File, error) {
	var ret []File
	for name := range r.MultipartForm.File {
		files, err := f.GetSlice(r, w, name)
		if err != nil {
			return ret, err
		}
		ret = append(ret, files...)
	}
	return ret, nil
}

// MaxSizeExceededError ..
type MaxSizeExceededError struct {
	error
}
