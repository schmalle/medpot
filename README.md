# medpot
Medpot is a honeypot that tries to emulate [HL7](https://en.wikipedia.org/wiki/Health_Level_7) / [FHIR](https://en.wikipedia.org/wiki/Fast_Healthcare_Interoperability_Resources) services



## Installation
Requires go 1.17 or newer

1. Installation of dependencies can be handled by running `bash scripts/dependencies.sh`
2. Now you can either do<br/>
    a) `bash scripts/run_medpot.sh` or `go run go/*.go` - To run the files<br/>
    
    b) `bash scripts/compile_medpot.sh` or  `go build -o medpot go/*.go` - To compile the files into a binary<br/>
    
    c) `make` and `make install` to create a copy on disk and also create all necessary files<br/>
   
    d) `bash scripts/compile_docker.sh` to create a docker container
3. You're done now! My suggestion is to check the arguments that you can send in to the program

By default the honeypot will try to bind and listen on port `2575`

## Arguments
Medpot utilizes an arugment parser to be able to less static in some areas, the supported arguments at this point of time are.<br>
```
#### Definied Arguments ####
--help, -h | Displays all defined arguments
--set_logo, -sl <value>  | Allows you to pick a logo that is shown on boot | options are [1, 2]
--set_port, -sp <value>  | Allows for a different port to be used, default = 2575
--set_log_location, -sll <value>  | Changes the directory where the logs will be placed, default = '/var/log/medpot/'
```

All arguments can easily be checked by passing the `-h` flag.

## Templates, configuration & Logs
Templates and configurations are located at `/etc/medpot/`

The default location for log files are located at `/var/log/medpot.log` but this can be changed by sending the `-sll` flag followed by the new location

<b>Example</b>
```
{
    "level":"info",
    "message":"Connection found",
    "time":"2018.09.09 17:20:49",
    "port":"57905",
    "ip":"127.0.0.1",
    "data":"TVMpEUk98S0FURV5TTUlUSF5FEK"
}
```
