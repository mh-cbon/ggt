package lib

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// File is descriptor
type File struct {
	PhysicalFile io.ReadCloser
	PhysicalPath string
	Name         string
	Index        int
	Filename     string
	PathExt      string
	MimeType     string
	Written      int64
}

// Fd ...
func (f *File) Fd() io.Reader {
	return f.PhysicalFile
}

type multiError struct {
	error
	errors []error
}

func (m multiError) set(e ...error) error {
	if len(e) == 0 {
		return nil
	}
	m.errors = e
	msg := ""
	for _, err := range e {
		msg += err.Error() + ":"
	}
	m.error = errors.New("multiple errors:" + msg[:len(msg)-1])
	return m
}

// Close ...
func (f *File) Close() error {
	return multiError{}.set(
		f.PhysicalFile.Close(),
		os.Remove(f.PhysicalPath),
	)
}

// AttachmentName ...
func (f *File) AttachmentName() string {
	return f.Filename
}

// ContentType ...
func (f *File) ContentType() string {
	return f.MimeType
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
		MaxSize:    (1 << 10) * 24,
		UploadPath: os.TempDir(),
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

	var tmp *os.File
	if tmp, err = ioutil.TempFile(f.UploadPath, "up-"); nil != err {
		return ret, err
	}
	var written int64
	if written, err = io.Copy(tmp, infile); nil != err {
		return ret, err
	}
	if err = infile.Close(); nil != err {
		return ret, err
	}
	if err = infile.Close(); nil != err {
		return ret, err
	}

	ret = File{
		PhysicalFile: tmp,
		PhysicalPath: filepath.Join(f.UploadPath, filepath.Base(tmp.Name())),
		Written:      written,
		Name:         name,
		Index:        index,
		Filename:     hdr.Filename,
		PathExt:      filepath.Ext(hdr.Filename),
		MimeType:     "",
	}

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
