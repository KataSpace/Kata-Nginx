
BUILD := $(shell date +"%Y%m%d-%H%M%S")

clean:
	rm -f bin/kn

kn-alpine: $(shell find . -type f|grep -v vendor|grep ".go$$")
	docker run -it --rm -v ${PWD}:/go/src/github.com/KataSpace/Kata-Nginx -w /go/src/github.com/KataSpace/Kata-Nginx --entrypoint go golang:alpine3.14 build -o bin/kn cli/kn/main.go
	docker build -t registry.cn-beijing.aliyuncs.com/vikings/kn:alpine-$(BUILD) .
	docker push registry.cn-beijing.aliyuncs.com/vikings/kn:alpine-$(BUILD)

kn: $(shell find . -type f|grep -v vendor|grep ".go$$")
	docker run -it --rm -v ${PWD}:/go/src/github.com/KataSpace/Kata-Nginx -w /go/src/github.com/KataSpace/Kata-Nginx --entrypoint go golang:1.17.1 build -o bin/kn cli/kn/main.go
	docker build -t registry.cn-beijing.aliyuncs.com/vikings/kn:$(BUILD) .
	docker push registry.cn-beijing.aliyuncs.com/vikings/kn:$(BUILD)