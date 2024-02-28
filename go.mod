module uc-go

go 1.21.1

replace github.com/minor-industries/theheads => ../heads

replace github.com/minor-industries/rfm69 => ../minor-industries/rfm69

replace github.com/minor-industries/max31856 => ../minor-industries/max31856

require (
	github.com/minor-industries/rfm69 v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	github.com/tinylib/msgp v1.1.9
	tinygo.org/x/drivers v0.27.0
	tinygo.org/x/tinyfs v0.2.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
