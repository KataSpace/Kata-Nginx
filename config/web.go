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

package config

import (
	"github.com/BurntSushi/toml"
	"github.com/KataSpace/Kata-Nginx/apis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// NewWebConfig Create a new web server configure
func NewWebConfig(path string) (cf *apis.Config, err error) {
	cf = defaultWebConfig()

	if path != "" {
		var fo []funcOption

		var _cf apis.Config
		_, err := toml.DecodeFile(path, &_cf)
		if err != nil {
			return cf, errors.WithMessage(err, "Parse Config Path")
		}

		if _cf.Debug {
			fo = append(fo, withDebug(_cf.Debug))
		}

		if _cf.Port != 0 {
			fo = append(fo, withPort(_cf.Port))
		}

		if _cf.Cache {
			fo = append(fo, withCache(_cf.Cache))
		}

		for _, opt := range fo {
			opt.apply(cf)
		}
	}

	output(cf)

	return cf, nil
}

type funcOption struct {
	f func(config *apis.Config)
}

func (o funcOption) apply(cf *apis.Config) {
	o.f(cf)
}

// defaultWebConfig create a default web config
func defaultWebConfig() (c *apis.Config) {
	return &apis.Config{
		Debug: true,
		Port:  10000,
		Cache: false,
	}
}

func withDebug(d bool) funcOption {
	return funcOption{f: func(config *apis.Config) {
		config.Debug = d
	}}
}

func withPort(port int) funcOption {
	return funcOption{f: func(config *apis.Config) {
		config.Port = port
	}}
}

func withCache(d bool) funcOption {
	return funcOption{f: func(config *apis.Config) {
		config.Cache = d
	}}
}

func output(cf *apis.Config) {
	log.Printf("Debug: %t", cf.Debug)
	log.Printf("Port: %d", cf.Port)
	log.Printf("Cache: %t", cf.Cache)
}
