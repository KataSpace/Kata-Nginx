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

package engine

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/KataSpace/Kata-Nginx/apis"
	"github.com/KataSpace/Kata-Nginx/apis/nginx"
	"github.com/KataSpace/Kata-Nginx/apis/web"
	"github.com/KataSpace/Kata-Nginx/tools"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CommonEngine struct {
	conf apis.Config
	sum  string
	node web.Node
}

func (ce CommonEngine) EngineInit(conf apis.Config) (err error) {

	//fill data in first time
	hasRunning, err := ce.CheckNginx()
	if err != nil {
		return errors.WithMessage(err, "Check Running Nginx")
	}

	if !hasRunning {
		return errors.New("Not Running Master Nginx")
	}

	data, err := ce.GetNginxContent()
	if err != nil {
		return errors.WithMessage(err, "GetNginxContent")
	}

	s, err := ce.GetNginxSum(data)
	if err != nil {
		return errors.WithMessage(err, "GetNginxSum")
	}

	ce.sum = s

	ingress, err := ce.FindDomain(data)
	if err != nil {
		return errors.WithMessage(err, "Generate Nginx Struct")
	}

	node, err := ce.GenerateTopology(ingress)
	if err != nil {
		return errors.WithMessage(err, "Generate Node Data")
	}

	ce.node = node
	return nil
}

func (ce CommonEngine) Reflash() (node web.Node, err error) {
	data, err := ce.GetNginxContent()
	if err != nil {
		return node, errors.WithMessage(err, "GetNginxContent")
	}

	newSum, err := ce.GetNginxSum(data)
	if err != nil {
		return node, errors.WithMessage(err, "GetNginxSum")
	}

	if ce.conf.Cache {
		if ce.sum == newSum {
			return ce.node, nil
		}
	}

	ingress, err := ce.FindDomain(data)
	if err != nil {
		return node, errors.WithMessage(err, "Generate Nginx Struct")
	}

	newNode, err := ce.GenerateTopology(ingress)
	if err != nil {
		return node, errors.WithMessage(err, "Generate Node Data")
	}

	ce.node = newNode
	ce.sum = newSum
	return ce.node, nil
}
func (ce CommonEngine) GetNginxContent() (string, error) {

	nginxPath, err := tools.CheckNginxProcess()
	if err != nil {
		return "", err
	}

	path := parseNginxCommand(nginxPath)

	log.Debugf("Get Nginx Binary: [%s]", path)

	command := strings.Split(path, " ")
	command = append(command, "-T")

	log.Debugf("GetNginxContent use: [%s]", strings.Join(command, " "))
	cmd := exec.Command("/bin/sh", "-c", strings.Join(command, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), errors.New(string(out))
	}

	return string(out), nil
}

func (ce CommonEngine) CheckNginx() (running bool, err error) {
	path, err := tools.CheckNginxProcess()
	if err != nil {
		return false, err
	}
	if path != "" {
		return true, nil
	}

	return false, nil
}

func (ce CommonEngine) GetNginxSum(data string) (sum string, err error) {

	result := tools.PoistionWithStr(data, 1, "#", "checksum")
	if len(result) == 0 {
		return "", errors.New("Not find check sum in this nginx configure")
	}

	_str := strings.Split(result[0], ":")
	if len(_str) != 2 {
		return "", errors.New("Wrong check sum format in this nginx configure")
	}

	return strings.TrimSpace(_str[1]), nil
}

func (ce CommonEngine) GenerateTopology(ingress nginx.Ingress) (node web.Node, err error) {
	return ingress2Node(ingress), nil
}

func (ce CommonEngine) NginxSum() (sum string) {
	return ce.sum
}

// FindDomain find http snippet from nginx configure
// data is nginx configure
func (ce CommonEngine) FindDomain(data string) (domains nginx.Ingress, err error) {
	httpSnippet := tools.FindHttpSnippet(data)
	if httpSnippet == "" {
		return domains, errors.New("Not find http in this config")
	}

	_domains := tools.PoistionWithStr(httpSnippet, -1, "##", "start server")
	if len(_domains) == 0 {
		return domains, errors.New("Not find domains in this config")
	}

	var serverName []string
	for _, d := range _domains {
		names := strings.Split(d, "start server")
		if len(names) == 2 {
			serverName = append(serverName, strings.TrimSpace(names[1]))
		}
	}

	_domainSnippet := tools.FindServerSnippet(data)

	if len(serverName) > len(_domainSnippet) {
		log.Errorf("There has [%d] server and has [%d] snippet", len(serverName), len(_domainSnippet))
		log.Errorf("server: [%v]", serverName)
		return domains, errors.New("Server name and server content num not match")
	}

	var do []nginx.Domain
	for idx, ds := range _domainSnippet {
		if idx < len(serverName) {
			s, _ := ce.FindServer(ds)
			s.Server = serverName[idx]
			do = append(do, s)
			continue
		}

		break
	}

	domains.Name = "Ingress"
	domains.Domains = do
	return domains, nil
}

// FindServer find all servers info from server snippet
// data is server snippet
func (ce CommonEngine) FindServer(data string) (servers nginx.Domain, err error) {
	locations, err := ce.FindLocation(data)
	servers.Locations = locations
	return
}

// FindLocation find all locations from server snippet
// data is server snippet.
func (ce CommonEngine) FindLocation(data string) (locations []nginx.Location, err error) {
	_locations := tools.FindLocationSnippet(data)
	if len(_locations) == 0 {
		return nil, errors.New("Not find locations in this server configure")
	}

	for _, l := range _locations {
		path := getLocation(l)
		if path == "" {
			return nil, errors.New("Not find location path in this location configure")
		}
		lo := nginx.Location{Path: path}

		md, _ := ce.FindLocationMetaData(l)
		lo.MetaData = []nginx.LocationMetaData{md}

		locations = append(locations, lo)
	}
	return locations, nil
}

// FindLocationMetaData find all location metadata from location snippet
// data is location snippet
func (ce CommonEngine) FindLocationMetaData(data string) (md nginx.LocationMetaData, err error) {

	_md := nginx.LocationMetaData{}
	_namespace := tools.PoistionWithStr(data, 1, "set", "namespace")
	if len(_namespace) != 0 {
		_md.Namespace = stripContentFromLocation(_namespace[0])
	}

	_ingress_name := tools.PoistionWithStr(data, 1, "set", "ingress_name")
	if len(_ingress_name) != 0 {
		_md.Ingress = stripContentFromLocation(_ingress_name[0])
	}

	_service_name := tools.PoistionWithStr(data, 1, "set", "service_name")
	if len(_ingress_name) != 0 {
		_md.Service = stripContentFromLocation(_service_name[0])
	}

	_service_port := tools.PoistionWithStr(data, 1, "set", "service_port")
	if len(_ingress_name) != 0 {
		_md.Port = stripContentFromLocation(_service_port[0])
	}

	return _md, nil
}

func getLocation(str string) string {
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(s, "location") {
			return getLocationPath(s)
		}
	}

	return ""
}

func getLocationPath(str string) string {
	return strings.TrimSpace(strings.Replace(strings.Replace(str, "location", "", 1), "{", "", 1))
}

// stripContentFromLocation get data content from location.
// e.g.
// get namespace default from `set $namespace      "default";`
func stripContentFromLocation(str string) string {
	data := strings.Split(str, "\"")
	if len(data) != 3 {
		return ""
	}
	return data[1]
}

// ingress2Node change ingress struct to node struct
func ingress2Node(i nginx.Ingress) (node web.Node) {
	var n []web.Node
	for _, s := range i.Domains {
		n = append(n, server2Node(s))
	}

	node.Name = i.Name
	node.Children = n
	return
}

// server2Node change server struct to node struct
func server2Node(s nginx.Domain) (node web.Node) {
	var n []web.Node
	for _, l := range s.Locations {
		n = append(n, location2Node(l))
	}
	node.Name = s.Server
	node.Children = n
	return
}

// location2Node change location struct to node struct
func location2Node(l nginx.Location) (node web.Node) {

	var n []web.Node
	for _, lo := range l.MetaData {
		n = append(n, web.Node{
			Name: fmt.Sprintf("%s-%s", lo.Service, lo.Port),
		})
	}

	node.Name = l.Path
	node.Children = n
	return
}

// parseNginxCommand get nginx execute binary from ps result.
// If nginx running via specify configure file, then return with this configure file
// If nginx running without specify configure file, then only return nginx.
// e.g.
// get `nginx` from `nginx -g daemon off;`
// get `nginx -c /etc/nginx/nginx.conf` from `nginx -c /etc/nginx/nginx.conf`
func parseNginxCommand(ps string) string {
	log.Debugf("Parse nginx command: [%s]", ps)
	result := ""
	command := strings.Split(strings.TrimSpace(ps), " ")

	if strings.Contains(ps, "-c") {
		saveOnce := false
		for _, c := range command {
			result += c + " "

			switch c {
			case "-c":
				saveOnce = true
			default:
				if saveOnce {
					return result
				}
			}
		}

		return result
	}

	return command[0]
}
