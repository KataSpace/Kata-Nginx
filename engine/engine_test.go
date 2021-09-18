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
	"reflect"
	"testing"

	"github.com/KataSpace/Kata-Nginx/apis"
	"github.com/KataSpace/Kata-Nginx/apis/nginx"
	"github.com/KataSpace/Kata-Nginx/apis/web"
)

func TestCommonEngine_GetNginxSum(t *testing.T) {
	type fields struct {
		conf apis.Config
		sum  string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSum string
		wantErr bool
	}{
		{
			name: "Normal Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;`,
			},
			wantSum: "3674439977205751794",
			wantErr: false,
		},
		{
			name: "Wrong Format Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `# configuration file /etc/nginx/nginx.conf:

# Configuration: checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;`,
			},
			wantSum: "",
			wantErr: true,
		},
		{
			name: "Wrong Content Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `# configuration file /etc/nginx/nginx.conf:
			
# setup custom paths that do not require root access
pid /tmp/nginx.pid;`,
			},
			wantSum: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CommonEngine{
				conf: tt.fields.conf,
				sum:  tt.fields.sum,
			}
			gotSum, err := ce.GetNginxSum(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNginxSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("GetNginxSum() gotSum = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestCommonEngine_FindDomain(t *testing.T) {

	var httpSnippet = `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;

daemon off;

worker_processes 4;

worker_rlimit_nofile 261120;

worker_shutdown_timeout 240s ;

events {
	multi_accept        on;
	worker_connections  16384;
	use                 epoll;
}

http {
	lua_package_path "/etc/nginx/lua/?.lua;;";
	
	upstream upstream_balancer {
		keepalive_requests 100;
		
	}
	
	
	## start server _
	server {
		server_name _ ;
		
		listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location / {
			
			set $namespace      "";
			set $ingress_name   "";
			set $service_name   "";
			set $service_port   "";
			set $location_path  "/";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}
		
	}
	## end server _
	

	## start server www.example.com
	server {
		server_name www.example.com ;
		
		listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location / {
			
			set $namespace      "";
			set $ingress_name   "";
			set $service_name   "";
			set $service_port   "";
			set $location_path  "/";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}
		
	}
	## end server www.example.com
}

stream {
	
	# TCP services
	
	# UDP services
	
}


# configuration file /etc/nginx/mime.types:

types {
    
    video/x-ms-wmv                                   wmv;
    video/x-msvideo                                  avi;
}

`

	type fields struct {
		conf apis.Config
		sum  string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDomains nginx.Ingress
		wantErr     bool
	}{
		{
			name: "Find all domains",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: httpSnippet},
			wantDomains: nginx.Ingress{
				Name: "Ingress",
				Domains: []nginx.Domain{
					nginx.Domain{
						Server: "_",
						Locations: []nginx.Location{
							nginx.Location{
								Path: "/",
								MetaData: []nginx.LocationMetaData{
									nginx.LocationMetaData{
										Namespace: "",
										Ingress:   "",
										Service:   "",
										Port:      "",
									},
								},
							},
						},
					},
					nginx.Domain{
						Server: "www.example.com",
						Locations: []nginx.Location{
							nginx.Location{
								Path: "/",
								MetaData: []nginx.LocationMetaData{
									nginx.LocationMetaData{
										Namespace: "",
										Ingress:   "",
										Service:   "",
										Port:      "",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CommonEngine{
				conf: tt.fields.conf,
				sum:  tt.fields.sum,
			}
			gotDomains, err := ce.FindDomain(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDomains, tt.wantDomains) {
				t.Errorf("FindDomain() gotDomains = %v, want %v", gotDomains, tt.wantDomains)
			}
		})
	}
}

func Test_stripContentFromLocation(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Test",
			args: args{str: "set $namespace      \"default\";"},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripContentFromLocation(tt.args.str); got != tt.want {
				t.Errorf("stripContentFromLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocationPath(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Test",
			args: args{str: "location /abc/eft {"},
			want: "/abc/eft",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocationPath(tt.args.str); got != tt.want {
				t.Errorf("getLocationPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocation(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Test",
			args: args{str: `listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location /abc/def {
			
			set $namespace      "";
			set $ingress_name   "";`},
			want: "/abc/def",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocation(tt.args.str); got != tt.want {
				t.Errorf("getLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommonEngine_FindLocationMetaData(t *testing.T) {
	type fields struct {
		conf apis.Config
		sum  string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMd  nginx.LocationMetaData
		wantErr bool
	}{
		{
			name: "Normal Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `location / {
			
			set $namespace      "default";
			set $ingress_name   "name";
			set $service_name   "service";
			set $service_port   "8000";
			set $location_path  "/abc";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}`},
			wantMd: nginx.LocationMetaData{
				Namespace: "default",
				Ingress:   "name",
				Service:   "service",
				Port:      "8000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CommonEngine{
				conf: tt.fields.conf,
				sum:  tt.fields.sum,
			}
			gotMd, err := ce.FindLocationMetaData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLocationMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMd, tt.wantMd) {
				t.Errorf("FindLocationMetaData() gotMd = %v, want %v", gotMd, tt.wantMd)
			}
		})
	}
}

func TestCommonEngine_FindLocation(t *testing.T) {
	type fields struct {
		conf apis.Config
		sum  string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantLocations []nginx.Location
		wantErr       bool
	}{
		{
			name: "Normal Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `server {
		server_name www.example.com ;
		
		listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location /abc/ef {
			
			set $namespace      "default";
			set $ingress_name   "name";
			set $service_name   "service";
			set $service_port   "8000";
			set $location_path  "/abc";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}
		
	}`},
			wantLocations: []nginx.Location{nginx.Location{
				Path: "/abc/ef",
				MetaData: []nginx.LocationMetaData{nginx.LocationMetaData{
					Namespace: "default",
					Ingress:   "name",
					Service:   "service",
					Port:      "8000",
				}},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CommonEngine{
				conf: tt.fields.conf,
				sum:  tt.fields.sum,
			}
			gotLocations, err := ce.FindLocation(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLocations, tt.wantLocations) {
				t.Errorf("FindLocation() gotLocations = %v, want %v", gotLocations, tt.wantLocations)
			}
		})
	}
}

func TestCommonEngine_FindDomain1(t *testing.T) {
	type fields struct {
		conf apis.Config
		sum  string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDomains nginx.Ingress
		wantErr     bool
	}{
		{
			name: "Normal Test",
			fields: fields{
				conf: apis.Config{},
				sum:  "",
			},
			args: args{data: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;

daemon off;

worker_processes 4;

worker_rlimit_nofile 261120;

worker_shutdown_timeout 240s ;

events {
	multi_accept        on;
	worker_connections  16384;
	use                 epoll;
}

http {
	lua_package_path "/etc/nginx/lua/?.lua;;";
	
	upstream upstream_balancer {
		keepalive_requests 100;
		
	}
	
	
	## start server _
	server {
		server_name _ ;
		
		listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location / {
			
			set $namespace      "root_default";
			set $ingress_name   "root_ingress";
			set $service_name   "root_service";
			set $service_port   "80";
			set $location_path  "/";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}
		
	}
	## end server _
	

	## start server www.example.com
	server {
		server_name www.example.com ;
		
		listen [::]:443 default_server reuseport backlog=32768 ssl http2 ;
		
		
		location /abc/ef {
			
			set $namespace      "default";
			set $ingress_name   "name";
			set $service_name   "service";
			set $service_port   "8000";
			set $location_path  "/abc";
			set $token_decode "";
			
			rewrite_by_lua_block {
				lua_ingress.rewrite({
					force_ssl_redirect = false,
					ssl_redirect = false,
					force_no_ssl_redirect = false,
					use_port_in_redirects = false,
				})
				balancer.rewrite()
				plugins.run()
			}
			
			log_by_lua_block {
				balancer.log()
				
				monitor.call()
				
				plugins.run()
			}
			
			proxy_redirect                          off;
			
		}
		
	}
	## end server www.example.com
}

stream {
	
	# TCP services
	
	# UDP services
	
}


# configuration file /etc/nginx/mime.types:

types {
    
    video/x-ms-wmv                                   wmv;
    video/x-msvideo                                  avi;
}

`},
			wantDomains: nginx.Ingress{
				Name: "Ingress",
				Domains: []nginx.Domain{
					nginx.Domain{
						Server: "_",
						Locations: []nginx.Location{
							nginx.Location{
								Path: "/",
								MetaData: []nginx.LocationMetaData{nginx.LocationMetaData{
									Namespace: "root_default",
									Ingress:   "root_ingress",
									Service:   "root_service",
									Port:      "80",
								}},
							},
						},
					},
					nginx.Domain{
						Server: "www.example.com",
						Locations: []nginx.Location{
							nginx.Location{
								Path: "/abc/ef",
								MetaData: []nginx.LocationMetaData{nginx.LocationMetaData{
									Namespace: "default",
									Ingress:   "name",
									Service:   "service",
									Port:      "8000",
								}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CommonEngine{
				conf: tt.fields.conf,
				sum:  tt.fields.sum,
			}
			gotDomains, err := ce.FindDomain(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDomains, tt.wantDomains) {
				t.Errorf("FindDomain() gotDomains = %v, want %v", gotDomains, tt.wantDomains)
			}
		})
	}
}

func Test_location2Node(t *testing.T) {
	type args struct {
		l nginx.Location
	}
	tests := []struct {
		name     string
		args     args
		wantNode web.Node
	}{
		{
			name: "Normal Test",
			args: args{l: nginx.Location{
				Path: "/abc",
				MetaData: []nginx.LocationMetaData{
					nginx.LocationMetaData{
						Namespace: "default",
						Ingress:   "ingress",
						Service:   "service",
						Port:      "8000",
					},
				},
			}},
			wantNode: web.Node{
				Name: "/abc",
				Children: []web.Node{
					web.Node{
						Name:     "service-8000",
						Children: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNode := location2Node(tt.args.l); !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("location2Node() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}

func Test_server2Node(t *testing.T) {
	type args struct {
		s nginx.Domain
	}
	tests := []struct {
		name     string
		args     args
		wantNode web.Node
	}{
		{
			name: "Normal Test",
			args: args{s: nginx.Domain{
				Server: "www.example.com",
				Locations: []nginx.Location{
					nginx.Location{
						Path: "/abc",
						MetaData: []nginx.LocationMetaData{
							nginx.LocationMetaData{
								Namespace: "default",
								Ingress:   "ingress",
								Service:   "service",
								Port:      "8000",
							},
						},
					},
				},
			}},
			wantNode: web.Node{
				Name: "www.example.com",
				Children: []web.Node{
					{
						Name: "/abc",
						Children: []web.Node{
							web.Node{
								Name:     "service-8000",
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
			if gotNode := server2Node(tt.args.s); !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("server2Node() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}

func Test_ingress2Node(t *testing.T) {
	type args struct {
		i nginx.Ingress
	}
	tests := []struct {
		name     string
		args     args
		wantNode web.Node
	}{
		{
			name: "Normal Test",
			args: args{i: nginx.Ingress{
				Name: "ingress",
				Domains: []nginx.Domain{
					{
						Server: "www.example.com",
						Locations: []nginx.Location{
							nginx.Location{
								Path: "/abc",
								MetaData: []nginx.LocationMetaData{
									nginx.LocationMetaData{
										Namespace: "default",
										Ingress:   "ingress",
										Service:   "service",
										Port:      "8000",
									},
								},
							},
						},
					},
				},
			}},
			wantNode: web.Node{
				Name: "ingress",
				Children: []web.Node{
					web.Node{
						Name: "www.example.com",
						Children: []web.Node{
							{
								Name: "/abc",
								Children: []web.Node{
									web.Node{
										Name:     "service-8000",
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
			if gotNode := ingress2Node(tt.args.i); !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("ingress2Node() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}

func Test_parseNginxCommand(t *testing.T) {
	type args struct {
		ps string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Use specify configure",
			args: args{ps: "nginx -c /etc/nginx/default.conf -g daemon off;"},
			want: "nginx -c /etc/nginx/default.conf ",
		},
		{
			name: "Not use specify configure",
			args: args{ps: "nginx -g daemon off;"},
			want: "nginx",
		},
		{
			name: "Normal Test",
			args: args{ps: ` nginx -g daemon off;
	`},
			want: "nginx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNginxCommand(tt.args.ps); got != tt.want {
				t.Errorf("parseNginxCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
