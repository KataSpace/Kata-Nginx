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

package v1

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/KataSpace/Kata-Nginx/apis"
	"github.com/KataSpace/Kata-Nginx/services/v1/tmpl"
	"github.com/gin-gonic/gin"
)

type WebService struct {
	engine apis.Engine
}

// GetPing 健康检查接口
func (ws *WebService) GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "running", "nginx check sum": ws.engine.NginxSum()})
}

// PostSync Engine缓存同步接口
func (ws *WebService) PostSync(c *gin.Context) {

	node, _, err := ws.engine.Reflash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, node)
}

func (ws *WebService) GetWeb(c *gin.Context) {
	needOriginal, _ := strconv.ParseBool(c.Query("original"))

	node, original, err := ws.engine.Reflash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	d, _ := template.New("d3").Parse(tmpl.D3)

	type data struct {
		Node     interface{}
		Original string
	}

	if needOriginal {
		d.Execute(c.Writer, data{Node: node, Original: original})
		return
	}

	d.Execute(c.Writer, data{Node: node, Original: ""})
}
func NewWebService(engine apis.Engine) *WebService {
	return &WebService{engine: engine}
}
