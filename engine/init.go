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
	"github.com/KataSpace/Kata-Nginx/apis"
	"github.com/KataSpace/Kata-Nginx/config"
)

func InitEngine(path string) (*apis.Config, apis.Engine, error) {
	conf, err := config.NewWebConfig(path)
	if err != nil {
		return nil, nil, err
	}

	engine, err := initEngine(conf)
	return conf, engine, err
}

func initEngine(conf *apis.Config) (apis.Engine, error) {
	return CommonEngine{conf: *conf}, nil
}
