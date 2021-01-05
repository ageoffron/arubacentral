# A Go client to execute commands against the Aruba Central REST API
```

Support the auth command for now with plans to support get a series of targets

```

## Installation
Build 
```
make build
```

## Usage
Update the config/config.json
```
{
    "clientID": "",
    "customerID": "",
    "clientSecret": "",
    "username": "",  
    "password": ""
}
```

Run client
```
./bin/arubacentral auth
./bin/arubacentral auth | jq

output:
{
  "access_token": "xxxxxx",
  "refresh_token": "xxxxxxxx",
  "token_type": "bearer"
}
```

Get help

  ./bin/central -h

```
Aruba Central management tool

Aruba Central cli to communicate with Aruba Central REST API

Usage:
  central [command]

Available Commands:
  auth        auth against Aruba Central API
  help        Help about any command

Flags:
      --config string     config file (default is ./config/config.json) (default "./config/config.json")
  -h, --help              help for central
      --loglevel string   log level [NONE, INFO, DEBUG] (default "NONE")

Use "central [command] --help" for more information about a command.
```

## todo
- add error handling on config file not present
- add more tests
- add config file path as cli argument
- add validation of config file format
- add refresh token support option

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
