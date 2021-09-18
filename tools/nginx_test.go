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
)

func TestPoistionWithStr(t *testing.T) {
	type args struct {
		str   string
		times int
		arch  []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "Normal Test",
			args: args{
				str: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

				# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;
`,
				times: -1,
				arch:  []string{"#", "checksum"},
			},
			wantResult: []string{
				"# Configuration checksum: 3674439977205751794",
				"# Configuration checksum: 3674439977205751794",
			},
		},
		{
			name: "Find 1 string",
			args: args{
				str: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

				# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;
`,
				times: 1,
				arch:  []string{"#", "checksum"},
			},
			wantResult: []string{
				"# Configuration checksum: 3674439977205751794",
			},
		},
		{
			name: "Find 2 string",
			args: args{
				str: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

				# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;
`,
				times: 2,
				arch:  []string{"#", "checksum"},
			},
			wantResult: []string{
				"# Configuration checksum: 3674439977205751794",
				"# Configuration checksum: 3674439977205751794",
			},
		},
		{
			name: "Find nothing",
			args: args{
				str: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

				# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid /tmp/nginx.pid;
`,
				times: -1,
				arch:  []string{"#", "checksum1"},
			},
			wantResult: nil,
		},
		{
			name: "Not has preFix str",
			args: args{
				str: `# configuration file /etc/nginx/nginx.conf:

# Configuration checksum: 3674439977205751794

				# Configuration checksum: 3674439977205751794

# setup custom paths that do not require root access
pid1 /tmp/nginx.pid;
pid /tmp/nginx.pid;
`,
				times: 2,
				arch:  []string{"", "pid"},
			},
			wantResult: []string{
				"pid1 /tmp/nginx.pid;",
				"pid /tmp/nginx.pid;",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := PoistionWithStr(tt.args.str, tt.args.times, tt.args.arch...); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("PoistionWithStr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

var fullConfig = `# configuration file /etc/nginx/nginx.conf:

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

func Test_findCodeSnippet(t *testing.T) {
	type args struct {
		conf  string
		key   string
		start string
		close string
		all   bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "Find http",
			args: args{
				conf:  fullConfig,
				key:   "http",
				start: " {",
				close: "}",
				all:   false,
			},
			wantResult: []string{
				`http {
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

}`,
			},
		},
		{
			name: "Find server",
			args: args{
				conf:  fullConfig,
				key:   "server",
				start: " {",
				close: "}",
			},
			wantResult: []string{
				`server {
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

}`,
			},
		},
		{
			name: "Find location",
			args: args{
				conf:  fullConfig,
				key:   "location",
				start: " {",
				close: "}",
				all:   false,
			},
			wantResult: []string{
				`location / {

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

}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := findCodeSnippet(tt.args.conf, tt.args.key, tt.args.start, tt.args.close, tt.args.all); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("findCodeSnippet() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestFindHttpSnippet(t *testing.T) {
	type args struct {
		conf string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal Test",
			args: args{conf: fullConfig},
			want: `http {
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

}`,
		},
		{
			name: "Not include http snippet",
			args: args{conf: `# configuration file /etc/nginx/nginx.conf:

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

`},
			want: "",
		},
		{
			name: "Wrong http snippet",
			args: args{conf: `# configuration file /etc/nginx/nginx.conf:

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
	lua_package_path "/etc/nginx/lua/?.lua;;";`},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindHttpSnippet(tt.args.conf); got != tt.want {
				t.Errorf("FindHttpSnippet() = %v, want %v", got, tt.want)
			}
		})
	}
}

var httpConfig = `http {
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
}`

func TestFindServerSnippet(t *testing.T) {
	type args struct {
		conf string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Normal Test",
			args: args{conf: httpConfig},
			want: []string{
				`server {
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

}`,
				`server {
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

}`,
			},
		},
		{
			name: "A correct and a wrong",
			args: args{conf: `http {
lua_package_path "/etc/nginx/lua/?.lua;;";

upstream upstream_balancer {
keepalive_requests 100;

}


## start server _
server {


}
## end server _

## start server www.example.com
server {
}
}
## end server www.example.com
}`},
			want: []string{
				`server {


}`,
				`server {
}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindServerSnippet(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindServerSnippet() = %v, want %v", got, tt.want)
			}
		})
	}
}

var serverConfig = `server {
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


  location /test {

    proxy_pass ;
  }

}`

func TestFindLocationSnippet(t *testing.T) {
	type args struct {
		conf string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Normal Test",
			args: args{conf: serverConfig},
			want: []string{
				`location / {

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

}`,
				`location /test {

proxy_pass ;
}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindLocationSnippet(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLocationSnippet() = %v, want %v", got, tt.want)
			}
		})
	}
}
