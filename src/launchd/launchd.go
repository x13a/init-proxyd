package launchd

import (
	"encoding/xml"
	"errors"
	"log"
	"os"
	"runtime"

	"./launch"
)

const DefaultConfig = "/Library/LaunchDaemons/me.lucky.init-proxyd.plist"

func Is() bool {
	return runtime.GOOS == "darwin"
}

func Sockets(config string) ([]int, error) {
	var names []string
	var err error
	if config != "-" {
		names, err = getSocketNames(config)
		if err != nil {
			return nil, err
		}
		if len(names) == 0 {
			return nil, errors.New("Socket names not found in " + config)
		}
	} else {
		names = []string{"Socket"}
	}
	res := make([]int, 0, len(names))
	for _, name := range names {
		fds, err := launch.ActivateSocket(name)
		if err != nil {
			return nil, err
		}
		res = append(res, fds...)
	}
	return res, nil
}

func getSocketNames(config string) ([]string, error) {
	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dec := xml.NewDecoder(file)
	const (
		Key   = "key"
		Array = "array"
		Dict  = "dict"
	)
	depth := -1
	found := false
	res := []string{}
Loop:
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch element := token.(type) {
		case xml.StartElement:
			switch element.Name.Local {
			case Key:
				if (!found && depth != 0) || (found && depth != 1) {
					continue
				}
				var key string
				if err = dec.DecodeElement(&key, &element); err != nil {
					log.Println(err)
					continue
				}
				if !found {
					if key == "Sockets" {
						found = true
					}
				} else {
					res = append(res, key)
				}
			case Array, Dict:
				depth++
			}
		case xml.EndElement:
			switch element.Name.Local {
			case Array:
				depth--
			case Dict:
				depth--
				if found && depth == 0 {
					break Loop
				}
			}
		}
	}
	return res, nil
}
