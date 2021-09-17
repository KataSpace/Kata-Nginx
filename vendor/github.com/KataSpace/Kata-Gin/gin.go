package KataGin

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine, nameConvert func(string) string, getMethods func(string) (string, string), objects ...interface{}) *gin.Engine {
	nf := defaultConvert
	gc := defaultGetMethods

	for _, o := range objects {
		t := reflect.TypeOf(o)
		v := reflect.ValueOf(o)

		for i := 0; i < v.NumMethod(); i++ {
			method := v.Method(i)

			if t.Method(i).IsExported() {
				f, ok := method.Interface().(func(ctx *gin.Context))
				if ok {

					if nameConvert != nil {
						nf = nameConvert
					}

					if getMethods != nil {
						gc = getMethods
					}

					method, name := gc(nf(t.Method(i).Name))
					switch strings.ToUpper(method) {
					case http.MethodGet:
						r.GET(name, f)
					case http.MethodPost:
						r.POST(name, f)
					case http.MethodPut:
						r.PUT(name, f)
					case http.MethodDelete:
						r.DELETE(name, f)

					}
				}
			}
		}
	}
	return r
}
