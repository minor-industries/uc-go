module host

go 1.21

toolchain go1.21.1

replace uc-go => ../

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/jessevdk/go-flags v1.5.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	uc-go v0.0.0-00010101000000-000000000000
)

require (
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
