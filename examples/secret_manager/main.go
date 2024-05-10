package main

import (
	"flag"
	"fmt"
	"github.com/SamuelTissot/flowconf"
	"github.com/SamuelTissot/flowconf/examples/internal/fake"
	"strings"
)

// Configurations holds our application configurations
type Configurations struct {
	Env               string `json:"environment" toml:"environment"`
	MeaningOfLife     int
	ServiceAccountKey string
	Text              string
	Secret            string
}

var configFiles []string

func init() {
	flag.Func(
		"conf", "List of configuration files", func(flagValue string) error {
			for _, f := range strings.Fields(flagValue) {
				configFiles = append(configFiles, f)
			}
			return nil
		},
	)
}

func main() {
	// the configuration files are passed to the application via
	// configuration flags
	// example:
	// go run --conf file_one.toml --conf file_two.json ....
	flag.Parse()

	// load sources
	sources, err := flowconf.NewSourcesFromFilepaths(configFiles...)
	if err != nil {
		// handle error
		panic(err)
	}

	// CREATE the Builder with static sources
	builder := flowconf.NewBuilder(sources...)
	// add a manager that will fetch values from an external provider like GCP Secret Manager,
	// here we fake it
	builder.SetSecretManagers(fakeManager())

	// BUILD the configurations
	conf := Configurations{}
	err = builder.Build(&conf)
	if err != nil {
		// handle error
		panic(err)
	}

	// USE THE VALUES
	fmt.Printf("%s:\t\t%s\n", "environment", conf.Env)
	fmt.Printf("%s:\t%d\n", "Meaning Of Life", conf.MeaningOfLife)
	fmt.Printf("%s:\t\t\t%s\n", "secret", conf.Secret)
	fmt.Printf("%s:\t\t\t%s\n", "text", strings.Trim(strings.ReplaceAll(conf.Text, "\n", "\n\t\t\t"), "\n\t"))
	fmt.Printf("%s:\t%v\n", "Service Account", strings.ReplaceAll(conf.ServiceAccountKey, "\n", ""))

	// OUTPUT
	// -------------------------------------------------------------
	// environment:            local
	// Meaning Of Life:        42
	// secret:                 !! locally overriden secret value !!
	// text:                   The quick brown
	//                         fox jumps over
	//                         the lazy dog.
	// Service Account:        {"type": "service_account","project_id": "project_id","private_key_id": "private key id","private_key": "-----BEGIN PRIVATE KEY-----Qualisque wisi commodo fabellas homero diam decore consetetur veniam quod duo splendide netus quis animal postulant voluptatibus necessitatibus deterruis-----END PRIVATE KEY-----","client_email": "service_account_name@project_id.iam.gserviceaccount.com","client_id": "client id","auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri": "https://oauth2.googleapis.com/token","auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service_account_name%40project-id.gserviceaccount.com"}
	// --------------------------------------------------------------
}

// fakeManager returns a fake SecretManager implementation with predefined secrets.
func fakeManager() flowconf.SecretManager {
	secrets := map[string]string{
		// simple value secret
		// but it's never fetch because of the locally overriden secret in [[conf-local.json]]
		"project/id/secret": "very secretive content",
		// example with a service gcp service account key (json) but needs to be in []byte format in configs
		"project/id/service_account_key": `{
"type": "service_account",
"project_id": "project_id",
"private_key_id": "private key id",
"private_key": "-----BEGIN PRIVATE KEY-----\nQualisque wisi commodo fabellas homero diam decore consetetur veniam quod duo splendide netus quis animal postulant voluptatibus necessitatibus deterruis\n-----END PRIVATE KEY-----\n",
"client_email": "service_account_name@project_id.iam.gserviceaccount.com",
"client_id": "client id",
"auth_uri": "https://accounts.google.com/o/oauth2/auth",
"token_uri": "https://oauth2.googleapis.com/token",
"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service_account_name%40project-id.gserviceaccount.com"
}`,
	}

	return fake.NewManager(secrets)
}
