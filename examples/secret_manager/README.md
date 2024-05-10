# Secret Manager Example

This example shows how to use a secret manager and override configuration values
with local configuration files.

## Setup

there are two configurations files
- The Main Configuration file that should be committed [conf.toml](./conf.toml)
- A local configuration file that in not committed and is used by each developer
  to override values for their own dev environment [conf-local.json](./conf-local.json)

**NOTE:** The local file does not need to be JSON, it's is just here to show that 
it's possible to mix and match file types.


## Running the application

I this example the configuration are passed via a `conf` flag. From within the 
[example folder](.) you can test the application by running it with different 
configurations

##### Just the conf.toml
```shell
go run . --conf conf.toml

# OUTPUTS
# --------------------------------------------------------------
# environment:            prd
# Meaning Of Life:        42
# secret:                 very secretive content
# text:                   The quick brown
#                         fox jumps over
#                         the lazy dog.
# Service Account:        {"type": "service_account","project_id": "project_id","private_key_id": "private key id","private_key": "-----BEGIN PRIVATE KEY-----Qualisque wisi commodo fabellas homero diam decore consetetur veniam quod duo splendide netus quis animal postulant voluptatibus necessitatibus deterruis-----END PRIVATE KEY-----","client_email": "service_account_name@project_id.iam.gserviceaccount.com","client_id": "client id","auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri": "https://oauth2.googleapis.com/token","auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service_account_name%40project-id.gserviceaccount.com"}
```

##### With the local overrides
```shell
go run . --conf conf.toml --conf conf-local.json

# OUTPUTS:
# --------------------------------------------------------------
# environment:            local
# Meaning Of Life:        42
# secret:                 !! locally overriden secret value !!
# text:                   The quick brown
#                         fox jumps over
#                         the lazy dog.
# Service Account:        {"type": "service_account","project_id": "project_id","private_key_id": "private key id","private_key": "-----BEGIN PRIVATE KEY-----Qualisque wisi commodo fabellas homero diam decore consetetur veniam quod duo splendide netus quis animal postulant voluptatibus necessitatibus deterruis-----END PRIVATE KEY-----","client_email": "service_account_name@project_id.iam.gserviceaccount.com","client_id": "client id","auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri": "https://oauth2.googleapis.com/token","auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service_account_name%40project-id.gserviceaccount.com"}
```