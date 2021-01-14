# Writed by yijian on 2021/01/14

.PHONY: build
build:
	go build -o client client.go # 编译出依赖 libc 程序; \
	go build -o server -ldflags '-linkmode "external" -extldflags "-static"' server.go # 编译出不依赖 libc 程序

.PHONY: clean
clean:
	rm -f server rt
