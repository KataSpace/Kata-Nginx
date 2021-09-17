# Kata-Gin

Package for Gin. Automatic generation of route maps

## How to use this package?

+ Declare a struct and implement some gin functions like this:

```go
package main

import (
	"net/http"

	kg "github.com/KataSpace/Kata-Gin"
	"github.com/gin-gonic/gin"
)

type example struct{}

func (e *example) GetAllName(c *gin.Context) {
	c.JSON(http.StatusOK, "Response by GetAllName")
}

func main() {
	r := gin.Default()
	r = kg.RegisterRouter(r, nil, nil, new(example))

	r.Run()
}

```

`Kata-Gin` will register `GetAllName` in `Gin` like the follow:
```go
GET /AllName  GetAllName
```

If there has a request `curl http://127.0.0.1:8080/AllName`, then get response: `Response by GetAllName`

## Details about `KataGin`

+ name convert function

1. slashConvert

If there are multiple consecutive uppercase letters in the function name, it will be ignored by default. If user wants to keep these uppercase letters, can use this function.

This function will add slash before every upper letter. It uses regex(`[A-Z][a-z]+|([A-Z]|[0-9]){%d}`,`%d` is n value.) for split string.

As I test, n = 3 gets the better effects.

If user set n = -1, then ignore this feature.
e.g.

```go
n == -1
AAAGetAllName ==> Get/All/Name

n == 2
GetAPIAllName ==> Get/AP/IA/Name

n == 3
GetAPIAllName ==> Get/API/All/Name

n == 10
GetAPIAllName ==> Get/All/Name
```