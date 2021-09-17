// Copyright (c) 2021-2021. The Kata-Nginx Authors.
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

type RunParams struct {
}

// Config Web Server Running Configure
type Config struct {
	// Debug whether output debug log
	// default false
	Debug bool
	// Port Web Server Listen Port
	// default 8000
	Port int
	// Cache whether enable cache.
	// If enable cache, only has a difference check sum, KN will re-parse configure.
	// default false.
	Cache bool
}
