# vechaintool
A developer-friendly, vechain smart contract development debug tools

#### complier

##### requires
- Go >= 1.11

##### clone project

```shell
git clone https://github.com/wupeaking/vechaintool && cd vechaintool
```

##### download module

```shell
go mod download
# using proxy
GOPROXY=https://goproxy.io go mod download
``` 

##### generate bin file

```shell
go build -o tools main.go
```

```shell
# if you need cross complie ...
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -v -o tools.exe main.go
```

#### Screenshot

1. setting page

![setting](./img/setting.png)

2. transfer page

![tx](./img/transfer.png)

3. contract page

![contract](./img/contract.png)

4. encode page

![contract](./img/encode.png)


#### demo

<img width=500 src="./img/demogif.gif">
