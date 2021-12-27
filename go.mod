module github.com/awe76/sagawf

go 1.16

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210920160938-87db9fbc61c7 // indirect
	github.com/awe76/sagaproc v0.0.0-20211227164624-9303312bc6d1
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	go-micro.dev/v4 v4.5.0
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
