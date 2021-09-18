// Copyright (c) 2021. The Kata-Nginx Authors.
//
// Licensed under the GPL License, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// If you has any question, plz contact me. ztao8607@gmail.com

package tools

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/golang-collections/collections/Stack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	httpCode     = "http"
	serverCode   = "server"
	locationCode = "location"
	commonStart  = " {"
	commonClose  = "}"
)

func CheckNginxProcess() (path string, err error) {
	c := "/bin/ps -ef | grep nginx | grep master| grep -v grep"
	cmd := exec.Command("/bin/bash", "-c", c)
	data, err := cmd.Output()
	if err != nil {
		return path, errors.WithMessage(err, "Get Nginx Pid")
	}

	log.Debugf("CheckNginxProcess: %s", string(data))

	//if len(data) > 0 {
	//	return path, nil
	//}

	p := strings.Split(string(data), "master process")

	log.Debugf("Parse Nginx Path: %v", p)
	if len(p) != 2 {
		return path, errors.New("Can not get nginx execute path")
	}

	return p[1], nil
}

// PoistionWithStr find all match string by arch
// If times == -1, then return all match strings.
// else only return specify times strings
// arch[0] is the prefix of string
func PoistionWithStr(str string, times int, arch ...string) (result []string) {
	if len(arch) == 0 {
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(str))
	idx := 0
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, arch[0]) {
			digo := true
			for _, a := range arch {
				if !strings.Contains(s, a) {
					digo = false
					break
				}
			}
			if digo {
				idx++
				result = append(result, s)
			}
		}

		if times <= 0 {
			continue
		}

		if times == idx {
			return result
		}

	}

	return result
}

// FindHttpSnippet find http snippet from nginx configure
// conf is full nginx configure
func FindHttpSnippet(conf string) string {
	result := findCodeSnippet(conf, httpCode, commonStart, commonClose, false)
	if len(result) == 0 {
		return ""
	}

	return result[0]
}

// FindServerSnippet find all server codeSnippet in http snippet.
// conf is full http configure
func FindServerSnippet(conf string) []string {
	return findCodeSnippet(conf, serverCode, commonStart, commonClose, true)
}

// FindLocationSnippet find all location codeSnippet in http snippet.
// conf is full server configure
func FindLocationSnippet(conf string) []string {
	return findCodeSnippet(conf, locationCode, commonStart, commonClose, true)
}

// findCodeSnippet find specify code snippet from fully configure
// key is the code id. e.g. http、server、location
// start is '{'
// and close is '}'
// more detail info please reference unit tests.
// ps:
// In some function, the param is a json object like :
// lua_ingress.rewrite({
//	force_ssl_redirect = false,
//	ssl_redirect = false
// })
// in this case, the `close` substr should be " {". If not will match a wrong string.
func findCodeSnippet(conf, key, start, close string, all bool) (result []string) {
	scanner := bufio.NewScanner(strings.NewReader(conf))

	save := false
	str := ""

	keyStack := stack.New()

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())

		if !save && strings.HasPrefix(s, key) && strings.HasSuffix(s, start) {
			save = true
			str = s
			keyStack.Push(s)
			continue
		}

		if save {
			str += "\n" + s
			switch {
			case strings.HasSuffix(s, close):
				keyStack.Pop()
			case strings.HasSuffix(s, start):
				keyStack.Push(s)
			}

			if keyStack.Len() == 0 {
				save = false
				result = append(result, str)
				str = ""
				if !all {
					break
				}
			}
		}
	}

	if keyStack.Len() != 0 {
		return nil
	}

	return result
}
