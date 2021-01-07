# A Go client to execute commands against the Aruba Central REST API
```

Cli tool against the Aruba Central API
Started as a project to learn GO.
Supports authentication, and several basic get command against the API (swarms, devices, aps)

```

## Installation

This will create the binary as ./bin/central

```
make build
```

## Usage
### First update the config/arubacentral.json
config file can be in $HOME/.config/arubacentral.json or current folder ./config/arubacentral.json
```
{
    "clientID": "",
    "customerID": "",
    "clientSecret": "",
    "username": "",  
    "password": ""
}
```

### Run the authentication 
this will output the access and refresh token, and write it to ./config/arubacentraltoken.json
```
./bin/central auth

{
  "access_token": "xxxxxx",
  "refresh_token": "xxxxxxxx",
  "token_type": "bearer"
}
```

### Get swarms 
```
./bin/central get swarms | jq

[...
{
    "firmware_version": "8.6.0.4_74969",
    "group_name": "xxxxx",
    "ip_address": "0.0.0.0",
    "name": "xxxxxx",
    "public_ip_address": "x.x.x.x",
    "status": "Up",
    "swarm_id": "xxxxxxxxx"
  },
...]

select the first one
./bin/central get swarms | jq ".[0] 
```

### Get devices 
```
./bin/central get devices | jq 
[...
{
  "aruba_part_no": "AP-505-US",
  "customer_id": "xxxxxxx",
  "customer_name": "xxxxxxxx",
  "device_type": "iap",
  "macaddr": "xx:xx:xx:xx:xx:xx",
  "model": "XXXXXX",
  "serial": "XXXXXXXXXXX"
}
{
  "aruba_part_no": "AP-505-US",
  "customer_id": "xxxxxxxx",
  "customer_name": "xxxxxxxx",
  "device_type": "iap",
  "macaddr": "xx:xx:xx:xx:xx:xx",
  "model": "XXXXXX",
  "serial": "XXXXXXXXXXX"
}
...]

select only the one with a partial Mac Adress match
./bin/central get devices | jq '.[] | select(.macaddr | contains("B8"))'
```

### Get access points
```
./bin/central get aps | jq
[...
 {
    "ap_deployment_mode": "IAP",
    "ap_group": "",
    "cluster_id": "",
    "model": "505",
    "radios": [
      {
        "band": 1,
        "index": 0,
        "macaddr": "xx:xx:xx:xx:xx:xx",
        "status": "Up"
      },
      {
        "band": 0,
        "index": 1,
        "macaddr": "xx:xx:xx:xx:xx:xx",
        "status": "Up"
      }
    ],
    "serial": "XXXXXXXXXX",
    "firmware_version": "8.7.0.0_75915",
    "ip_address": "x.x.x.x",
    "last_modified": 1610038000,
    "mesh_role": "Unknown",
    "name": "xxxxxxxxxx",
    "status": "Up",
    "macaddr": "xx:xx:xx:xx:xx:xx",
    "notes": "",
    "public_ip_address": "x.x.x.x",
    "subnet_mask": "m.m.m.m",
    "group_name": "yourgroupname",
    "site": "yoursitename",
    "swarm_id": "yourswarm_unique_id",
    "swarm_master": false
  },
  ...]
```


## Get help
```
  ./bin/central -h


Aruba Central management tool

Aruba Central cli to communicate with Aruba Central REST API

Usage:
  central [command]

Available Commands:
  auth        auth using secrets from config file
  get         get [devices, swarms, aps]
  help        Help about any command

Flags:
      --config string     config file (default is ./config/config.json) (default "./config/config.json")
  -h, --help              help for central
      --loglevel string   log level [NONE, INFO, DEBUG] (default "NONE")

Use "central [command] --help" for more information about a command.
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
