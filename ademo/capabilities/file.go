package capable

import (
	"fmt"
	"io"

	ggt "github.com/mh-cbon/ggt/lib"
)

// File ...
type File struct{}

// ReadOneFile ...
func (c File) ReadOneFile(fileName io.Reader) {
	fmt.Printf(`fileName %v
    `, fileName)
}

// ReadOneTmpFile ...
func (c File) ReadOneTmpFile(fileName ggt.File) {
	fmt.Printf(`fileName %v
    `, fileName)
}

// ReadMany ...
func (c File) ReadMany(fileName ggt.File, fileName2 ggt.File) {
	fmt.Printf(`fileName %v
    `, fileName)
	fmt.Printf(`fileName2 %v
    `, fileName2)
}

// ReadSlice ...
func (c File) ReadSlice(fileName []io.Reader, fileName2 []ggt.File) {
	fmt.Printf(`fileName %v
    `, fileName)
	fmt.Printf(`fileName2 %v
    `, fileName2)
}

// ReadAll ...
func (c File) ReadAll(fileValues []io.Reader) {
	fmt.Printf(`fileName %v
    `, fileValues)
}

// ReadAll2 ...
func (c File) ReadAll2(fileValues []ggt.File) {
	fmt.Printf(`fileName %v
    `, fileValues)
}

// WriteFile ...
func (c File) WriteFile() (fileResBody io.Reader, fileResName string, fileResContentType string) {
	fmt.Printf(`fileResBody %v
    `, fileResBody)
	fmt.Printf(`fileResName %v
    `, fileResName)
	fmt.Printf(`fileResContentType %v
    `, fileResContentType)
	return
}

// WriteFile2 ...
func (c File) WriteFile2() (fileResBody *ggt.File) {
	fmt.Printf(`fileResBody %v
    `, fileResBody)
	return
}

// WriteFile3 ...
func (c File) WriteFile3() (fileResBody ggt.File) {
	fmt.Printf(`fileResBody %v
    `, fileResBody)
	return
}
