### pwKeeper

pwKeeper is a client-server application, that allows to securely save vulnerable information on remote server and access it only with client application.

The following documents are supported:
+ Note - just some text
+ Card - credit card information
+ Credential - login credentials: login and password
+ File - files up to 14 MB

#### Client application 

Client application is a terminal UI and should be run from terminal (command prompt in Windows).
Supported architectures:
+ darwin_amd64
+ darwin_arm64
+ linux_amd64
+ linux_arm64
+ windows_amd64

##### Building

To build client for required architecture use `build-client.sh` script with exactly one argument: the platform to build for. To build for current platform use `current` as argument

##### Using
Client application allows to register a new user as using the already registered one.

Multiple client applications may be run simultaneously from different locations. The data will always be in sync. Upon launching client application caches data from server. All updates are first sent to server and appear in app only after server saved document. Files are always stored on server and fetched upon download request.

To run application just pass the server address in `-a` flag:
```shell
./client -a 127.0.0.1:3200
```
or for Windows
```shell
client_windows_amd64.exe -a 127.0.0.1:3200
```

Application support the following flags to run with:
+ `-v` `--version` print version and exit
+ `-a` `--address` server address host:port (default 127.0.0.1:3200)
+ `-m` `--mouse` enable terminal mouse support (experimental, may be unstable)
+ `--tls-ca-file` path to CA tls certificate, enables secured server connection
+ `--tls-insecure` disables validation of server certificate, use for testing only

Most of the flags have corresponding environment variables, which can be examined using `-h` or `--help` flag.

#### Server Application

Server application uses MongoDB as backend storage and gRPC transport between client and server.

Secured TLS connection is supported between client and server.

##### Building

To build server application for current architecture use provided `build-server.sh` script

##### Running

The most preferred way to run server app is to run it containerized. With included `docker-compose.yml` as an example, you can quickly build and launch application. Command `docker compose up` will build and run server container and mongodb. The application will be listening port `3200`.

Server application options:
+ `-v` `--version` print version and exit
+ `-d` `--db-uri` mongodb connection string, default is `mongodb://mongo:27017`
+ `-a` `--address` server bind address host:port (default 0.0.0.0:3200)
+ `-k` `--token-key` token signing key. To keep user sessions alive after server restart, please provide it. Otherwise, random key will be generated on each server restart.
+ `--tls-cert-file` path to server tls certificate file, if it is signed with intermediate CA, this file should contain full certificate chain. This option enables tsl.
+ `--tls-key-file` path to server certificate file (should be without password protection)
+ `-l` `--loglevel` log level: -1..2, where -1=Debug 0=Info 1=Warning 2=Error, default is `0`
+ `--debug` switches logs output to text mode (default is json) and turns log level to debug

Most of the flags have corresponding environment variable, which can be examined using `-h` or `--help` flag.