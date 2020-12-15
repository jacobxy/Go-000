module newservice

go 1.14

replace myService => ../../

require (
	google.golang.org/grpc v1.34.0
	myService v0.0.0-00010101000000-000000000000
)
