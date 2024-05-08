# FlowConf
The extendable, configurable application configuration package for go.

## Example Usage
```go

	type MyConfiguration struct {
		MeaningOfLife   int
		Cats            []string
		Pi              float64
		Perfection      []int
		BackToTheFuture time.Time
		Secret          string
		Tag             string `json:"TagValue" toml:"TagValue"`
	}

	tomlSource := `
MeaningOfLife = 42
Cats = [ "James", "Bond" ]
Pi = 3.14
Perfection = [ 6, 28, 496, 8128 ]
BackToTheFuture = 1985-10-21T01:22:00Z

# get the secret from gcp secret manager ( we are not using a secret manager in this example so this will stay the same)
Secret = "@managerprefix::projects/id/secrets/name-of-secret"

# change the name with tags
TagValue = 'tag value'
`

	overrideWith := `
{
	"Cats" : [ "Bob", "Morane" ]
}
`

	// Build the source from files or io.reader
	// to build sources from multiple files you can use
	// sources, err := flowconf.NewSourcesFromFilepaths("config.toml", "config-local.json", "another-config.json")
	//
	// HERE we will only build manually one source from an io.Reader
	source := flowconf.NewSource("example.toml", flowconf.Toml, strings.NewReader(tomlSource))
	overrideSource := flowconf.NewSource("override.json", flowconf.Json, strings.NewReader(overrideWith))

	// create the builder
	builder := flowconf.NewBuilder(source, overrideSource)
	// add secret manager if you need to
	// builder.SetSecretManagers(....)

	// populate the configuration with the builder
	conf := new(MyConfiguration)
	err := builder.Build(conf)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println()
	fmt.Println(conf.MeaningOfLife)
	fmt.Println(strings.Join(conf.Cats, ", "))
	fmt.Println(conf.Pi)
	fmt.Println(conf.Tag)
	// ...

	// Output:
	// 42
	// Bob, Morane
	// 3.14
	// tag value
```