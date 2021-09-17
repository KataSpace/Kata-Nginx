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

package apis

import (
	"github.com/KataSpace/Kata-Nginx/apis/nginx"
)

type NginxEngine interface {
	// CheckNginx check if it has a running nginx instance.
	// data is full nginx configure
	CheckNginx(data string) (running bool, err error)

	// GetNginxSum get currently nginx configure file sum.
	// data is full nginx configure
	GetNginxSum(data string) (sum string, err error)

	// GenerateTopology generate nginx invoke topology data
	GenerateTopology() (err error)

	// NginxSum  return nginx configure sum
	NginxSum()(sum string)
}

type ParseEngine interface {
	// FindDomain find all domains in specify nginx configure
	// data is configure snippet
	FindDomain(data string) (domains []string, err error)

	// FindLocation find all locations in specify server snippet.
	// data is server snippet content
	FindLocation(data string) (locations []string, err error)

	// FindLocationMetaData find location metadata in each location snippet.
	// data is location snippet content
	FindLocationMetaData(data string) (md nginx.LocationMetaData, err error)
}

type Engine interface {
	// EngineInit Init Runtime Engine
	// config the global runtime configure
	EngineInit(conf Config) (err error)
	NginxEngine
	ParseEngine
}
