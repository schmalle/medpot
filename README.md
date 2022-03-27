# medpot
Is a honeypot that tries to emulate HL7 / FHIR honeypot

## Installation
Requires go 1.17 or newer

1. Installation of dependencies can be handled by running `bash scripts/dependencies.sh`
2. Now you can either do
    a) `bash scripts/run_medpot.sh` or `go run go/*.go` - To run the files<br>
    b) `bash scripts/compile_medpot.sh` or  `go build -o medpot go/*.go` - To compile the files into a binary<br>
    c) `bash scripts/build_docker.sh` - To create a working docker image<br>
3. You're done now! My suggestion is to check the arguments that you can send in to the program

By default the honeypot will try to bind and listen on port `2575`

## Arguments


## Log
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
