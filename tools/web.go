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
	"fmt"

	"github.com/KataSpace/Kata-Nginx/apis/nginx"
	"github.com/KataSpace/Kata-Nginx/apis/web"
)

// ConvertNginxToNode convert nginx.Ingress data to web.Node
func ConvertNginxToNode(ingress nginx.Ingress) web.Node {
	var node []web.Node

	for _, s := range ingress.Domains {
		node = append(node, fillDomainToNode(s))
	}

	return web.Node{
		Name:     ingress.Name,
		Children: node,
	}
}

// fillDomainToNode
// e.g.
// l := nginx.Domain{
//		Server: "test.example.com",
//		Locations: []nginx.Location{
//			{
//				Path: "/",
//				MetaData: []nginx.LocationMetaData{
//					{
//						Namespace: "name",
//						Ingress:   "ingress",
//						Service:   "service",
//						Port:      "8000",
//					},
//				},
//			},
//		},
//	}
//  convert to
//  web.Node{
//		Name: "test.example.com",
//		Children: {
//			web.Node{
//				Name: "/",
//				Children: []web.Node{
//					web.Node{
//						Name:     "service:8000",
//						Children: nil,
//					},
//				},
//			},
//		},
//	}
func fillDomainToNode(domain nginx.Domain) web.Node {
	var node []web.Node

	for _, d := range domain.Locations {
		node = append(node, fillLocationToNode(d))
	}

	return web.Node{
		Name:     domain.Server,
		Children: node,
	}
}

// fillLocationToNode
// e.g.
// l := nginx.Location{
//		Path:     "/",
//		MetaData: []nginx.LocationMetaData{
//			{
//				Namespace: "name",
//				Ingress:   "ingress",
//				Service:   "service",
//				Port:      "8000",
//			},
//		},
//	}
//  convert to
//  web.Node{
// 		Name: "/",
// 		Children:{
//  		web.Node{
// 				Name: "service:8000",
// 				Children:nil
//			}
//		}
//	}
func fillLocationToNode(location nginx.Location) web.Node {
	var node []web.Node

	for _, l := range location.MetaData {
		node = append(node, fillLocationMetaDataToNode(l))
	}

	return web.Node{Name: location.Path, Children: node}
}

// fillLocationMetaDataToNode
// e.g.
// l := nginx.LocationMetaData{
//  Namespace: "name",
//  Ingress:   "ingress",
//  Service:   "service",
//  Port:      "8000",
// }
// convert to
//  web.Node{
// 		Name: "service:8000",
// 		Children:nil
//	}
func fillLocationMetaDataToNode(md nginx.LocationMetaData) web.Node {

	return web.Node{Name: fmt.Sprintf("%s:%s", md.Service, md.Port)}
}
