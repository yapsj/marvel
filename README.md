# Marvel API

Marvel API is a project that calling Marvel Server to get characters information

## Installation (windows only)

- Download the project from github. 
- Download REDIS for windows [here](https://github.com/downloads/dmajkic/redis/redis-2.4.5-win32-win64.zip)

The project is using Windows version to execute. The redis  is using localhost:6379 by default.

## Description
The application will do get all characters ID and get character by ID that provided from Marvel Server.

## Config.json
```json
{
    "host": "localhost",
    "port": 8080,
    "redis":{
        "host": "localhost",
        "port":6379
    },
    "logFilePath": "", 
    "cacheTTL": 30000
}
```
**host** is the address of the http server.

**port** is the port number.

**redis** is the configuration section for REDIS server, it has **host** and **port** too.

**logFilePath** is the folder path to specify where the logs.txt is going to put.

**cacheTTL** is the cache time (SECONDS) that saved in REDIS server, it expires in the key in REDIS server in X SECONDS.

** Not specifying logFilePath will use current directory of the executable as the log output folder.

## Running the project
- Extract the downloaded REDIS server and nagivate to the 64bit folder, double click **redis-server.exe** to start the server.
- Run the script below to start the application.
```bash
./marvel.exe ./config.json {publickey} {privatekey}
```

## Testing the project
In the project folder /integration_test, execute the command below:-
```bash
go run .
```
It will print out the list of the characterIds and 1 random character information get from the characterIds.


## License
[MIT](https://choosealicense.com/licenses/mit/)
