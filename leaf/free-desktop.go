package leaf

import (
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/danopia/romaine-head/common"
	"gopkg.in/ini.v1"
)

const appFolder string = "/usr/share/applications"
const dirFolder string = "/usr/share/desktop-directories"

// updateEntries chroot, '/usr/share/desktop-directories', 'Directory', '.directory'
// updateEntries chroot, '/usr/share/applications', 'Application', '.desktop'

var keyTypes = map[string]string{
	"NoDisplay":       "bool",
	"Hidden":          "bool",
	"DBusActivatable": "bool",
	"Terminal":        "bool",
	"StartupNotify":   "bool",

	"OnlyShowIn": "string-set",
	"NotShowIn":  "string-set",
	"Actions":    "string-set",
	"MimeType":   "string-set",
	"Implements": "string-set",
	"Categories": "string-set",
	"Keywords":   "string-set",
}

func (s *Stalk) watchFreeDesktop() {
	log.Println("Watching freedesktop files")
	s.watchFreeDesktopRoot(appFolder, "Application", ".desktop")
	s.watchFreeDesktopRoot(dirFolder, "Directory", ".directory")
}

func (s *Stalk) watchFreeDesktopRoot(root string, Type string, ext string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	// Blank slate
	s.Sink <- common.Packet{
		Cmd:     "collection wipe",
		Context: "fd-apps",
	}

	for _, file := range files {
		if file.IsDir() {
			// TODO: recursive
			// TODO: dir/app.ext has id dir-app.ext
		} else if path.Ext(file.Name()) == ext {
			filePath := path.Join(root, file.Name())

			s.Sink <- common.Packet{
				Cmd: "set field",
				Extras: map[string]interface{}{
					"Collection": "fd-apps",
					"Id":         file.Name(),
					"Field":      "Entry",
					"Value":      readFdFile(filePath, file.Name()),
				},
			}
		}
	}
}

func readFdFile(path string, id string) map[string]interface{} {
	obj := make(map[string]interface{})

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ini.Load(raw)
	if err != nil {
		log.Fatal(err)
	}
	cfg.BlockMode = false

	section, err := cfg.GetSection("Desktop Entry")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range section.Keys() {
		if strings.Contains(entry.Name(), "[") {
			// Ignore localized strings for now

		} else if keyTypes[entry.Name()] == "bool" {
			obj[entry.Name()] = (entry.String() == "true") || (entry.String() == "1")

		} else if keyTypes[entry.Name()] == "string-set" {
			obj[entry.Name()] = strings.Split(strings.TrimRight(entry.String(), ";"), ";")

		} else {
			obj[entry.Name()] = entry.String()
		}
	}

	return obj
}
