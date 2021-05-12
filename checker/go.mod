module checker

go 1.16

require (
	github.com/julienschmidt/httprouter v1.3.1-0.20200921135023-fe77dd05ab5a
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf
	google.golang.org/grpc v1.37.0
	gosship v0.0.0
)

replace gosship v0.0.0 => ../service
