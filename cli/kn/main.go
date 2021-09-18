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

package main

import (
	"fmt"
	"os"

	kg "github.com/KataSpace/Kata-Gin"
	"github.com/KataSpace/Kata-Nginx/config"
	"github.com/KataSpace/Kata-Nginx/engine"
	v1 "github.com/KataSpace/Kata-Nginx/services/v1"
	"github.com/gin-gonic/gin"
)

func main() {

	conf, eg, err := engine.InitEngine(os.Getenv(config.KataConfigPath))
	if err != nil {
		panic(err)
	}
	ws := v1.NewWebService(eg)

	r := gin.Default()
	r = kg.RegisterRouter(r, nil, nil, ws)

	panic(r.Run(fmt.Sprintf(":%d", conf.Port)))
}
