// Code generated by go-bindata.
// sources:
// index.html
// DO NOT EDIT!

package tpl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x58\xcb\x8e\xdb\xb6\x1a\x5e\x7b\x9e\x82\x50\x90\xcd\x41\x64\xcb\x9e\xf1\x5c\x14\x8d\x81\x83\x93\xb3\x0b\x0e\x0e\x9a\xa2\x5d\xd3\xe2\x2f\x9b\x18\x8a\x14\x28\xfa\x16\x63\x80\x2e\xdb\x02\x5d\x15\x45\xda\x22\x40\x90\x45\x16\x05\xba\xe8\x4b\xd5\xcf\x51\x90\x94\x28\xca\x92\x27\x8e\x8d\x11\x68\xf2\xbf\x5f\x3e\xfe\x9a\x64\xa9\x72\x86\x18\xe6\x8b\xfb\x60\xa5\xb2\xf0\x36\x98\x5d\xa0\x64\x09\x98\xcc\xd0\x05\x42\x89\xa2\x8a\xc1\xec\xf0\xc3\xcf\x7f\xff\xf4\xf1\xf0\xfe\xfb\xc3\x8f\x9f\x92\x91\xdd\x33\xc7\xa5\xda\x31\x40\x6a\x57\xc0\x7d\xa0\x60\xab\x46\x69\x59\x06\xb3\x8b\xc1\x60\x2e\xc8\x0e\xed\x2f\x06\x83\x41\x26\xb8\x0a\x33\x9c\x53\xb6\x8b\x51\xf0\x0d\x48\x82\x39\x0e\x5e\xa0\x12\xf3\x32\x2c\x41\xd2\xec\xa5\x23\x2b\xe9\x5b\x88\xd1\x38\x8a\x9e\x9b\xbd\x39\x4e\x1f\x16\x52\xac\x38\x89\xd1\xb3\x2c\xd2\x5f\xb3\x9f\x0a\x26\x64\x8c\x9e\x45\xe6\x63\xb7\x56\xb2\xd4\x7b\x04\x32\xbc\x62\x4a\xef\x3d\x5e\x54\x76\xbc\xb8\x18\x0c\x14\x9e\x33\x30\x0b\xa2\x9f\xcb\xb1\x79\x4e\xcc\x53\x5a\x43\x73\x2c\x17\x94\xc7\x28\x2a\xb6\x46\x64\x81\x09\xa1\x7c\xe1\x36\x8c\xbc\x7f\xc7\x8c\xf2\x87\x17\x66\xb5\xa6\x25\x55\x40\xec\x0f\x9c\x2a\xba\x06\x2b\xa9\xb6\xef\xf2\xf2\xfa\xfa\xee\xce\xe3\x5d\x8a\x35\xc8\x36\x4d\x14\x69\x2a\x47\xb3\x1c\xdb\x63\xa7\x7c\x7c\x5d\x6c\xb5\x05\xe8\xb6\xd8\xa2\xf1\xb4\xb2\x6d\x2e\x24\x81\x46\x85\x21\x28\x05\xa3\xc4\x3b\x0d\x37\x94\xa8\xa5\x31\xdf\xfc\x55\x82\x3a\xd1\xbe\xac\xa2\x6d\xf6\x36\x40\x17\x4b\x15\x23\x2e\x64\x8e\x59\x37\x0b\x8d\xb9\xce\x83\xcc\x7c\x1a\x0f\x26\xed\x70\x4e\x6a\xfd\x37\xc5\x16\x4d\xbb\xfa\x6f\x2b\xf5\xce\xe3\xcb\x86\xce\x08\x34\xa9\xb3\x32\x2b\xc7\x52\xc1\x18\x2e\x4a\x88\x51\xbd\x3a\x22\x56\xc4\xd2\x57\x11\x70\x05\xb5\x11\x92\x84\x73\x09\xf8\x21\x46\x0f\x00\x45\x88\x99\xf5\x71\xb3\xa4\x0a\xc2\xb2\xc0\x29\x68\xdf\x37\x12\x17\x66\x5f\xe7\x2b\x63\x62\x13\xa3\x25\x25\x04\xb8\xd9\xd4\x85\x1e\x36\x27\xc0\x18\x2d\x4a\x5a\x9e\xf4\xc2\xf3\x76\x72\x55\x34\xb5\x49\xe8\x7a\x98\x02\x57\x70\x54\x80\x63\x1d\xad\xb1\x5f\x75\x43\x02\x8c\xe6\xd4\x51\x12\x5a\x16\x0c\xef\x62\x34\x67\x22\x7d\x30\x3a\x96\x55\xe2\x5c\x89\xf4\x99\x6e\x85\x71\xa1\xe0\xa4\x1c\xbf\xea\xad\x25\x51\x7f\xde\x6e\x7a\x9b\xd4\xd5\x82\x57\xe0\xd7\xd7\xaf\x5e\x35\xda\x33\x21\x9c\x1f\x4d\x95\x4f\xbd\xe2\xac\xb2\xac\x44\x11\xa3\xab\xba\xb0\xd1\xb3\xd4\x7c\x1a\x41\xe5\x2a\xcf\xb1\xdc\x79\xe5\x71\x9c\xee\x3e\x6b\xad\xec\xb9\x50\x4a\xe4\x31\x1a\x37\xe2\xe1\x06\x6e\x52\x38\x21\x7e\xa8\x11\x11\xa4\x2b\xac\x56\xaf\xcc\x05\x23\xdd\x50\xa4\x90\x66\x77\x69\xab\x5b\x3d\x6d\x9d\xd3\x3e\x93\x8e\x3a\xeb\xc8\xa4\x54\x70\x05\x5c\x39\x9b\x5a\xda\x1b\x6f\x9c\xf6\xa8\x4f\x70\x0f\x56\x8c\xfd\xc2\x6b\xab\xd4\xb9\x0b\xdb\xb9\x0b\x65\x55\x77\x75\xf6\x4c\x77\x60\x46\x17\x3c\x46\xe6\xec\x49\x61\xac\x2d\x8c\x41\x76\x4a\x96\x3e\x3a\x21\x4a\x91\x61\xe5\x06\x6b\xe1\x44\x25\xed\xcc\x1c\x37\x52\x64\x4b\x4a\xed\xe0\xb9\x62\xa4\x0f\xf3\xed\xee\x98\xe8\x6f\x1b\xab\x86\x04\x32\x89\x17\x94\x67\xc2\x72\x98\xdd\x90\xe1\x9d\x58\xa9\x18\x65\x74\x0b\xe4\x34\x87\x03\x3a\x0d\x6c\x1a\xb5\x62\x64\xe0\x2d\xd4\x1b\x26\x84\x61\x5e\x86\x27\x4f\xad\xf9\x73\xc8\x84\x04\x46\x39\xf4\xd9\x9c\x35\x77\xac\x03\xfc\x7a\xcb\xf2\xe3\x4c\x81\x3c\xc1\x1e\x45\xb7\xc7\xec\xcd\x96\x65\xb7\xf0\x74\x04\x48\x5c\x70\x2f\xc0\xa6\x0e\x4c\x1e\xaa\x18\x9d\xae\x31\x73\xe4\xa3\xaa\x4f\x6b\xf7\x1d\xf1\xbf\xea\x3c\x6f\x35\x4a\x18\x20\x72\xcd\xb8\x6d\x8d\x10\x7d\xa3\xcc\x1b\xb1\x92\x29\xa0\x37\x7a\x84\x41\xff\x97\x22\x68\x79\x39\x8d\xa6\xd1\x74\x7a\x04\xa9\x2f\x5b\x73\x86\xf9\xa5\xe3\x16\x3a\xf0\x1e\x5e\x79\x93\x80\x37\xa2\x5c\x9a\xe7\x95\x79\x4e\x9f\x82\x20\x6f\xac\xb8\xbd\xba\x79\xd9\xb5\xfa\x2b\xcc\x60\x83\x77\x41\x67\xe2\xf0\xc0\x72\x02\x79\xe7\x3a\xf7\x27\x86\xe1\x95\x25\xa8\x5c\xb1\x48\x3d\x1e\x4e\x3d\x36\xec\x05\x9f\x40\x2a\x24\x56\x54\xf0\x26\xaf\xb5\xa1\x94\x2f\x41\xd2\x26\x7f\xd8\x6f\x9e\x0e\xf7\x8a\x13\x5b\x69\x8e\xbe\x40\xb8\x3d\x4f\x8d\xb3\x5b\xec\xf5\xa6\xb9\x63\x05\x57\x98\x72\x90\x68\xc5\x3e\x43\x9e\x0a\x02\x1d\x7f\xa3\xe1\x5d\xe5\x6f\x6f\x01\xfc\x47\xf3\xd4\xf9\xaf\x43\xaa\x05\xd9\xfc\x35\xcb\xcb\x66\x79\xd5\x2c\xa7\x6e\x59\xb8\xd5\x8a\xa1\xc6\x12\xbf\x9f\x5c\xf8\x5c\xa3\x50\x6e\x2a\xa8\xef\x02\x6f\x37\x2d\x4c\xb2\x4b\xdc\xc4\x4d\xd6\x77\xa6\x3f\xf9\x14\x12\x42\x37\xfb\xf8\xc3\x92\x85\x0d\x37\x2d\x9d\x44\x9b\xde\xa8\x39\x8b\xcc\x28\x38\x71\x33\xac\xf3\x2b\x74\x46\xb6\x01\xf2\x8b\x9a\xf9\x8b\x41\x42\xdf\x0e\x5d\xca\xf6\x35\x63\xb1\x29\x84\xad\x92\x38\x2c\x73\xcc\xd8\x93\x38\x95\x63\xca\xc3\x6a\x54\xb0\xf7\x9a\x28\xa9\x2d\x5d\x09\x0c\xeb\x37\x04\xe3\xfa\xdb\x90\x72\x02\xdb\x18\x79\xbe\xae\x41\x2a\x9a\x62\x16\xe6\x94\x90\x7a\xa2\x71\x9b\x95\x75\xf6\xcc\x32\x0d\x92\x91\x79\xfd\xd2\x6f\x62\xc9\xa8\x7e\x65\x4b\x34\x5c\xd9\x97\xb3\xb9\x44\x23\xbb\x22\x74\x8d\x52\x86\xcb\xf2\x3e\xb0\x51\xd3\x6f\x79\x08\xa1\xfd\x1e\x49\xcc\x17\x80\x86\x6f\x20\xd5\x76\xa2\xc7\x47\xcd\x80\x12\x7b\x97\x55\x3c\xad\xfb\xad\xb9\x7d\x02\x44\xc9\x7d\x20\xa1\x10\x52\xed\xf7\xc3\xff\xe1\x1c\x1e\x1f\x03\x64\x46\x89\xfb\x40\x4f\x61\x81\x51\x8f\x50\x92\x0a\xa6\x53\x5d\x58\xb5\xe6\x77\x4d\x37\x89\x9e\x07\x95\x9d\x67\x9f\x8c\x4f\xf3\x4c\xcf\x3a\x49\x46\x6d\x8b\x12\x65\xc2\x56\x71\x29\x59\x7b\x6e\x73\xd9\xeb\x93\xa6\x23\xfa\x05\xa4\x2c\x30\xbf\x0f\xa6\x01\x32\x39\x72\x11\xae\x45\xb8\xc9\xa4\x1e\x27\x82\x19\x4a\x74\x9f\x20\xdd\x27\xf7\xc1\x6d\x30\x73\xc1\x4b\x46\xfa\x60\x86\x92\x91\x22\xce\x87\x91\x92\x6e\xdd\x58\x56\x4d\x7f\x4f\x98\xe6\x47\xb1\x6d\xda\xec\xf0\xee\xd3\xe1\xfd\x07\x5f\xcb\x67\x39\x7e\xfd\xf0\x04\xc7\xb8\x8f\xe3\xcf\x8f\x87\x5f\xfe\x3a\xa9\x63\xda\xc3\xf1\xfe\xbb\xc3\xbb\x3f\xec\xff\x1a\xbe\x94\xef\xb7\xdf\xbb\x7c\x36\x72\x4d\x8d\xbf\xd6\x23\x4a\x55\xe0\x26\x96\x2d\x05\x47\x42\xf7\xfb\xe1\x2b\xac\x4c\x4e\x8e\x2c\xe9\x12\x7e\x0b\xf0\x40\xf0\xee\x2c\xda\xaf\x69\x0e\xe5\x59\x94\xff\xc5\x92\x51\x28\xd5\x59\xc4\xaf\xb1\xea\x90\xd6\xfe\x03\x27\xb5\xd7\xc9\xa8\x29\xf4\x64\x64\x1a\xda\x52\xfb\x10\x61\xdf\xd0\x2a\x88\x68\xf1\x27\x23\x42\xd7\x33\x64\x30\xc5\x2c\x35\xf2\x58\x81\xc9\x68\xa9\x72\x36\xbb\xf8\x27\x00\x00\xff\xff\x92\xdd\x8d\xca\x51\x12\x00\x00")

func indexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_indexHtml,
		"index.html",
	)
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"index.html": indexHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

