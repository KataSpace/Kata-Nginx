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
	"reflect"
	"testing"

	"github.com/KataSpace/Kata-Nginx/apis/nginx"
	"github.com/KataSpace/Kata-Nginx/apis/web"
)

func Test_fillLocationMetaDataToNode(t *testing.T) {
	type args struct {
		md nginx.LocationMetaData
	}
	tests := []struct {
		name string
		args args
		want web.Node
	}{
		{
			name: "Normal Test",
			args: args{md: nginx.LocationMetaData{
				Namespace: "nameSpace",
				Ingress:   "inGress",
				Service:   "serVice",
				Port:      "8000",
			}},
			want: web.Node{
				Name:     "serVice:8000",
				Children: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillLocationMetaDataToNode(tt.args.md); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillLocationMetaDataToNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fillLocationToNode(t *testing.T) {
	type args struct {
		location nginx.Location
	}
	tests := []struct {
		name string
		args args
		want web.Node
	}{
		{
			name: "Normal Test",
			args: args{nginx.Location{
				Path: "/",
				MetaData: []nginx.LocationMetaData{
					{
						Namespace: "nameSpace",
						Ingress:   "inGress",
						Service:   "serVice",
						Port:      "8000",
					},
				},
			}},
			want: web.Node{
				Name: "/",
				Children: []web.Node{
					{
						Name:     "serVice:8000",
						Children: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillLocationToNode(tt.args.location); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillLocationToNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fillDomainToNode(t *testing.T) {
	type args struct {
		domain nginx.Domain
	}
	tests := []struct {
		name string
		args args
		want web.Node
	}{
		{
			name: "Normal Test",
			args: args{domain: nginx.Domain{
				Server: "test.example.com",
				Locations: []nginx.Location{
					{
						Path: "/",
						MetaData: []nginx.LocationMetaData{
							{
								Namespace: "name",
								Ingress:   "ingress",
								Service:   "serVice",
								Port:      "8000",
							},
						},
					},
				},
			}},
			want: web.Node{
				Name: "test.example.com",
				Children: []web.Node{
					{
						Name: "/",
						Children: []web.Node{
							{
								Name:     "serVice:8000",
								Children: nil,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillDomainToNode(tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillDomainToNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertNginxToNode(t *testing.T) {
	type args struct {
		ingress nginx.Ingress
	}
	tests := []struct {
		name string
		args args
		want web.Node
	}{
		{
			name: "Normal Test",
			args: args{nginx.Ingress{
				Name: "ingress",
				Domains: []nginx.Domain{
					{
						Server: "test.example.com",
						Locations: []nginx.Location{
							{
								Path: "/",
								MetaData: []nginx.LocationMetaData{
									{
										Namespace: "name",
										Ingress:   "ingress",
										Service:   "serVice",
										Port:      "8000",
									},
								},
							},
						},
					},
				},
			}},
			want: web.Node{
				Name: "ingress",
				Children: []web.Node{
					{
						Name: "test.example.com",
						Children: []web.Node{
							{
								Name: "/",
								Children: []web.Node{
									{
										Name:     "serVice:8000",
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertNginxToNode(tt.args.ingress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertNginxToNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
