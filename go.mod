module github.com/minor-industries/uc-go

go 1.21.1

replace tinygo.org/x/drivers => github.com/minor-industries/drivers v0.0.0-20240403222057-a994655999c5

require (
	github.com/minor-industries/rfm69 v0.0.2
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
