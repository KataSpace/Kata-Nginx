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
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/KataSpace/Kata-Nginx/apis"
)

func setupConfigure() {
	var c1 = `debug=true
port=80
cache=true	
	`

	err := ioutil.WriteFile("/tmp/c1.toml", []byte(c1), 0777)
	if err != nil {
		panic(err)
	}

	var c2 = `
port=80
cache=true	
	`

	err = ioutil.WriteFile("/tmp/c2.toml", []byte(c2), 0777)
	if err != nil {
		panic(err)
	}


	var c3 = `debug=true
	cache=true	
	`

	err = ioutil.WriteFile("/tmp/c3.toml", []byte(c3), 0777)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setupConfigure()

	os.Exit(m.Run())
}

func TestNewWebConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantCf  *apis.Config
		wantErr bool
	}{
		{
			name: "Default Configure",
			args: args{path: ""},
			wantCf: &apis.Config{
				Debug: false,
				Port:  10000,
				Cache: false,
			},
			wantErr: false,
		},
		{
			name: "Parse all config file",
			args: args{path: "/tmp/c1.toml"},
			wantCf: &apis.Config{
				Debug: true,
				Port:  80,
				Cache: true,
			},
			wantErr: false,
		},
		{
			name: "Parse config file without debug",
			args: args{path: "/tmp/c2.toml"},
			wantCf: &apis.Config{
				Debug: false,
				Port:  80,
				Cache: true,
			},
			wantErr: false,
		},
		{
			name: "Parse config file without port",
			args: args{path: "/tmp/c3.toml"},
			wantCf: &apis.Config{
				Debug: true,
				Port:  10000,
				Cache: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCf, err := NewWebConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCf, tt.wantCf) {
				t.Errorf("NewWebConfig() gotCf = %v, want %v", gotCf, tt.wantCf)
			}
		})
	}
}
