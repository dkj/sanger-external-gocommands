# gocommands
iRODS Command-line Tools written in Go


## Download pre-built binary
Please download binary file (bundled with `tar` or `zip`) at (https://github.com/cyverse/gocommands/releases)[https://github.com/cyverse/gocommands/releases].
Be sure to download a binary for your target system architecture.


## Build from source
Use `make` to build `gocommands`. Binaries will be created on `./bin` directory.

```bash
make
```

## How to use

### Using a persistent configuration (compatible to iCommands)
`gocommands` can create a configuration that is compatible to `icommands`.
Run `gocmd init` to configure iRODS account for access in interactive manner.
This will create a configuration directory `.irods` in your home directory and several configuration files will be created.
Now, it is ready to go.

Use any commands, such as `gocmd ls`, to access iRODS.

### Using an external configuration file 
`gocommands` can read configuration from an external file in `YAML` or `JSON` format.
In this example, I'll show you an example `YAML` configuration.

Create a YAML file with iRODS account, say `config.yaml`.
```yaml
irods_host: "data.cyverse.org"
irods_port: 1247
irods_user_name: "your username"
irods_zone_name: "iplant"
irods_user_password: "your password"
```

Then run any commands, such as `gocmd ls`, with `-c` option.
```bash
gocmd ls -c config.yaml
```

You can omit password if you want. In the case, `gocommands` will ask you to type a password in runtime.

### Using environmental variables 
`gocommands` can read configuration from environmental variables.

Set environmental variables
```bash
export IRODS_HOST="data.cyverse.org"
export IRODS_PORT=1247
export IRODS_USER_NAME="your username"
export IRODS_ZONE_NAME="iplant"
export IRODS_USER_PASSWORD="your password"
```

Then run any commands, such as `gocmd ls`, with `-e` option.
```bash
gocmd ls -e
```

You can omit password if you want. In the case, `gocommands` will ask you to type a password in runtime.
