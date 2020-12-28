# A Go client to auth against Aruba Central REST API

## Installation
Build 
```
make build
```

## Usage
Update the config/config.production.json
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
./bin/arubacentral_authclient 
./bin/arubacentral_authclient | jq '.access_token'"

output:
{
    "access_token":"xxxxxxxxxxxxxxx",
    "refresh_token":"xxxxxxxxxxxxxx",
    "token_type":"bearer"
}
```

Run client with DEBUG level
```
./bin/arubacentral_authclient  -loglevel DEBUG
```

Get help
```
./bin/arubacentral_authclient  -h
```

## todo
- add error handling on config file not present
- add more tests
- add config file path as cli argument
- add validation of config file format
- add refresh token support option

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
