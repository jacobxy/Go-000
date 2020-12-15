module newclient

go 1.14

require (
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.2.2
	myService v0.0.0-00010101000000-000000000000
)

replace myService => ../..
