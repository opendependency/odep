module github.com/opendependency/odep

go 1.17

require (
	github.com/gofrs/flock v0.8.1
	github.com/opendependency/go-spec v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.27.1
)

// Testing
require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.10.1
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20201224043029-2b0845dc783e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/opendependency/go-spec => ../go-spec
