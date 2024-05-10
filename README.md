# FlowConf
The extendable, configurable application configuration and secret manager package for the Go language.

At its core, **FlowConf** permits the declaration of configurations and secrets in an easily understandable format. **FlowConf** does not hide the implementation details or where each secret is stored. In order to resolve a secret, a user needs to have permission to view that secret.


## Usage

```go
import "github.com/SamuelTissot/flowconf"

// load source from file
sources, err := flowconf.NewSourcesFromFilepaths("fileOne.json", "fileTwo.toml")
if err != nil {
	// handle error
}

// create the builder
builder := flowconf.NewBuilder(sources...)

var conf Configuration // your declared Configuration struct
err = builder.Build(&conf)
if err != nil {
// handle error
}

// USE Configuration
```


## Examples

#### With Only Static (files) Source
see [static source example](./doc_test.go)

#### Full Example with a Secret Manager
see [Secret Manager Example](./examples/secret_manager)