module uc-go

go 1.19

replace github.com/cacktopus/theheads => ../heads

replace tinygo.org/x/tinyfs => ../tinyfs

require (
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	github.com/tinylib/msgp v1.1.8
	tinygo.org/x/drivers v0.24.0
	tinygo.org/x/tinyfs v0.2.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
