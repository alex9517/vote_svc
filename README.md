# Vote Service (microservice, Go + Go kit)

This is a web-service (microservice, backend, ...) that can support (provide, perform) miscellaneous kinds
of voting, polling, election, questioning, etc.

To use it,
- you need a [database](#database) (Apache Cassandra, noSQL), and you have to create [ manually ] all the required database stuff (keyspace, tables, records); this is **not so difficult** as it sounds. The project has `./scripts` dir containing **database-related scripts and files**;
- you have to create a proper client to communicate with the service; the project includes a simple demo client (HTML, CSS, JS, jQuery) demonstrating the interaction between client and service. You can use it as is, with minor modifications, but remember: it's just a demo.

There is a short description of [essential scripts](#scripts) at the end of this doc.

Note that `$PROJ_ROOT` in the context of this document designates the name of the **root dir of your project** (the dir where `go.mod` resides).


### Warning!

The current version of the service does not provide authentication of clients (voters, electors). So, it cannot be used for the sensitive procedures like .. I fear to say it .. "presidential elections" or something similar and potentially scandalous.

Also, an "evil" person can vote multiple times unless your client app cannot provide some kind of unique identification of the voter. Currently it's just a UUID created by the browser and stored in a local cookie. Each time when somebody tries to vote using this specific user account, JS checks the cookie and blocks request if the appropriate cookie exists. However, an experienced user can easily bypass this trivial protection arrangement.


## Technical details

### Programming Tools and Dependencies

The language is **Go** (aka Golang), the framework is **Go kit** see [Go kit main site](https://gokit.io).

Note that **Go 1.22** (or higher) is required to compile the code.


### Endpoints

In case of success, all endpoints return HTTP Status 200.

| Method | Endpoint | Description |
| ------------ | ---------------------- | ------------------------------------- |
| GET | `/health` | Returns `HealthStatus` struct with 2 fields: `health_status`, `health_message`. If service is running and ready to accept requests, status is 200 and message is "UP"; otherwise, it can be HTTP 500 or 503 ("Unsufficient memory", "Database connection failure", etc) |
| GET | `/votes/{id}` | Returns `VoteData` struct with a nested array of `Contender` structs |
| GET | `/votes/{id}/results` | .. (same as previous) |
| PUT | `/votes` | Updates (increments) the `Count` of the selected contender and saves UUID of the voter in the database. Browser client app is supposed to send `vote_id int, co_id int16, user_id string`. The response can be `nil` in case of success, or error (HTTP 500, 503, 400, 403) |


### Ports, potocols and certificates

Service is supposed to be accessed using **HTTPS** as a typical RESTful web-service, **gRPC is not supported** in this version. **HTTP** can be used, but it is not recommended.

All stuff related to certificate creation is in `./scripts2/` dir. By default service tries to load self-signed certificate from `/etc/x509/https` dir, which should contain `server.pem` (certificate) and `server.key` (private key). This config is supposed to be used when service is running as a **Docker container** (see details in [Docker image](#docker_img) below). If you start this service as a regilar app using `./runsvc.sh`, you should specify the cert-related files location in this script. Currently, the location is `./cert/`, and certificate was created for the `localhost`. Probably, you'll need to create your own cert. [More about certificates..](#certs)

| Protocol | Default TCP port |
| --------------- | -----------------|
| HTTP | 8081 |
| HTTPS | 8443 |
| HTTP Debug | 8080 |


### Environment variables used by this service

| Variable | Description |
| -------- | ----------- |
| ALLOW_ORIGINS | CORS-related; a comma-separated list of URLs allowed to access this service; there must be NO SPACES between items; ReactJS client during development is usually specified as http://localhost:3000 |
| GOMEMLIMIT | This is Go specific env var (since Go 1.19) that affects the RAM usage by Go runtime; |


### <a name="database"></a>Database

This version of the service stores the data in the Apache Cassandra database (noSQL). Unless you know,

- Cassandra (?) Query Language (CQL) has many similarities with SQL, but it's not SQL;
- Referential integrity: `none`;
- Data duplication: `ok`;

In the process of development and testing I was using the simplest version: a single node database in the Docker container. See an [official image](https://hub.docker.com/_/cassandra) on the Docker Hub. It can be installed with
```
docker pull cassandra:latest
```
Currently all database objects are supposed to be created manually using `cqlsh` and the prepared scripts (`*.cql`). If you have some
experience with Apache Cassandra, you should know that the `cqlsh` is a standard cmdline client app for interacting with Cassandra database (like `sqlplus` for Oracle database or `psql` for PostgreSQL). It comes with the above mentioned docker image.

You may ask: why this web servce does not have a functionality to create tables, insert records, etc? - Probably because a microservice, as I
understand it, should not be overloaded with the miscellaneous functionality. Of course it's a controversial point.


#### How to work with Apache Cassandra running in a Docker container

The commands described here are supposed to be run in terminal mode (e.g. Gnome Terminal). The `cqlsh` is also a cmdline tool. The scripts mentioned below were created to save time. Some of them may need editing (hostnames, IPs, options), and, of course, you can ignore those scripts and enter all commands manually.

The service (Vote Service, the main subj of this proj) would not work without a database created and filled with some records beforehand.

**IMPORTANT!** Remember that `vote_id` differentiates one project (topic, poll, survey, data set) from another. Technically, you can use the same table for important data and for some transient test data. Of course, this is not a good practice, but it's possible.

To start the database container use `runCQLsvc.sh`. Technically, Cassandra's Docker image does not include a database per se, that container runs a database server instance. The database must be created on a Docker volume, but I skip this part, see the doc related to Cassandra Docker image.

To create tables, load data, perform admin tasks, you need `cqlsh`. It can be started
with `runCQLcli.sh` (to be exact, it must be: `$HOME/bin/runCQLcli.sh`). This script is simple and limited, probably you can make something better.

**IMPORTANT!** If you're going to type CQL statements manually, it doesn't matter how or from where you start `cqlsh`. But if you want to exec CQL scripts, it becomes a little bit more complicated: before starting `cqlsh`, you must move (`cd ...`) to `$PROJ_ROOT/scripts` dir (assuming that it contains your CQL scripts). When `cqlsh` starts and successfully connects to the database (no auth required unless you configured it),

- you can manually enter CQL statements and special commands like `exit`;

- you can load/exec CQL script named `data.cql` stored in `$PROJ_ROOT/scripts/` dir. That name is hard-coded: you can have many `*.cql` files, but when you need to exec one of them, copy its contents to `data.cql` (erase/overwrite ...). Alternatively, you can edit or totally rewrite `runCQLcli.sh`;

Assuming that `cqlsh` starts and connects to the database successfully, you should do following steps:

- create a keyspace (you can use `keyspace.cql`, it creates a keyspace named `polls`);
- create `votes` table (see `table1.cql`);
- create `voters` table (see `table2.cql`);
- insert necessary records into `votes` table;

The `records1.cql` script is just a demo. You can load (insert) these records to see how service works. You must load these or similar records if want to run service tests (see `./pkg/service/service_test.go`). The `vote_id = 1` used by this demo proj should be considered occupied, reserved, etc., because it's hard-coded in some tests. Probably, it would be reasonable to avoid 1..10 range.

Anyway, real working records must be totally yours.

Do not forget to rename/copy each script in turn to `data.cql`.

So, keyspace:
```
cqlsh> source './scripts/data.cql';
```
If a keyspace was successfully created, you can create tables and insert records:
```
cqlsh> use polls;
cqlsh:polls> source './scripts/data.cql';
```
It looks like the same cmd is repeated many times, but each time you exec totally different `data.cql`.


## Rate limiting

This app uses **token-based rate limiter**. It may be necessary to perform some tests to ajust the values. By default the limit is 300 requests/sec for GET endpoints and 60 req/sec for PUT endpoint. To change defaults, you have to modify the source code and to rebuild (see `main.go`). To change values quickly, use cmdline parameters `rate-limit` and `rate-limit-put` (and restart, of course).

Probably the database is the main limiting factor. However, load tests (1000 users, ramp-up 20 sec) perform good on my moderate desktop.

Note that `/health` endpoint has the essentially reduced rate limit 2 req/sec because ... security reasons, ... and why would you need to ping it more often?


## <a name="howto"></a>Howto ...

### About testing

Each `pkg` dir includes a test:

- `./pkg/endpoint/endpoint_test.go`
- `./pkg/http/handler_test.go`
- `./pkg/service/service_test.go`

So, to run a unit test on `endpoint.go`, chdir to `./pkg/endpoint` and do
```
go test -v
```
Similarly you can test the transport layer (`./pkg/http`).

However, there is a little problem with the `service` test: it's not exactly a unit test, it requires a database with test/demo records, because interaction with database is the most complicated part of the service; besides this, there is not much to test.

The `runtests.sh` script runs all tests.


#### Load testing

The `testdata` dir holds files and scripts related to load testing with `JMeter` (see [Apache JMeter](https://jmeter.apache.org)). To use these `*.jmx`, you should probably modify them, at least set your hostname. Also, you must edit `runJMeter.sh` script. Or, ignore it and run these tests as you like. Or, just skip this part.


### How to build and run

(see also how to build and run a [Docker image](#docker_img)).

First of all, go to project's root dir (that's where `go.mod` resides). Download (or update) dependencies:
```
go get -u ./...
```
Or, you can do
```
make dep
```
Also, you should probably exec:
```
go mod tidy -v
```
To build the app, run
```
make
```
There are no config files, you only need `vote-svc` (executable), `runsvc.sh` (simple startup script which is probably good for Linux OS only), and `cert` directory with certificate and private key files. Without certificate the service will auto switch to HTTP.

Before running the startup script **check/edit parameters**, then:
```
./runsvc.sh
```
All non-default config parameters must be specified on the command line, or as environment variables. And this is the reason why it's better to use a startup script. To see available options, you can try
```
./vote-svc --help
```


### <a name="docker_img"></a>How to create the Docker Image

The `Dockerfile.alpine` (multistage) is used to build an image based on `golang:1.22-alpine`.

Assuming that the Docker is already installed and running, exec `docker images` to find what you have, and pull what you need. Then,
```
make image_alpine
```
It first builds a **static binary**, then creates an image (~24MB);

If it goes without errors, just to be sure, exec
```
docker images
```
and look for `ms-vote-svc-alpine`.


### How to run the Docker Image

To run the created image use the script
```
runServiceAlpine.sh
```
Before trying, you should probably inspect the script carefully and modify it according to your requirements.

Pay attention to env variables and TCP ports. If you don't know what is `GOMEMLIMIT`, see official Go doc (this env var does not work with old Go versions). `ALLOW_ORIGINS` is important because it allows (or not) your browser client to pass through the CORS protection. This is the URL of the web-server you use to run your client app (HTML/CSS/JS or ReactJS or ... whatever you prefer).

Notice `-d` in that script. It is used to run the service in the background. However, if there are some problems (e.g. HTTPS does not work), and you want to see app's log, remove `-d` (temporarily).

Also, you have to create a network required to run those images (unless it already exists). See `scripts/dock_net.sh`.

To access the containerized services, your client apps should be able to resolve services' hostnames. The simplest way is to add the appropriate record to `/etc/hosts` of your system, e.g.:
```
...
172.25.70.26    s7026
...
```
The name assigned to the started container can be used to stop it later:
```
docker stop vote-service
```


## <a name="certs"></a>Self-signed certificates

The **self-signed server certificate** is usually created for a specific host/node/server that cannot be accessed from the Internet at all, or can be made accessible through a gateway. That means, limit your self-signed certificates usage to local area networks (LANs).

Since cert is usually bound to IP or hostname, this project has 2 (two) different self-signed certificates: one for regular use, one for the Docker container. Each cert consists of two files:

- `service.pem` (certificate),
- `service.key` (private key)

The `./cert/` dir contains certificate files that are supposed to be used when you run the service as a regular executable (not as a Docker container). The `./runsvc.sh` can be used to save you from typing cmdline parameters like `service-cert`, `service-key`, etc.

The Docker container uses slightly different files loaded from `./scripts2/` dir. These files (`service.pem` and `service.key`) are external to the image, so you can create and use your own certs without rebuilding Docker image.

Of course, you can keep the certificate files wherever you want, but in that case you have to edit `run*.sh` scripts or to rebuild app. Or, you can handle all this stuff as you like.

**IMPORTANT!** If service fails to locate and load certificate, it will start anyway, but the protocol will be HTTP, not HTTPS!


### About certificate creation

**You can ignore this chapter** if you prefer to handle these things in your own way.

The `create_cert.sh` script in `./scripts2` can be used to create **self-signed** certificates.

But (!) before creating new cert/key, carefully inspect the content of the `certificate.conf` file, and make sure that `CN`, `DNS`, `IP` are good for you. Or change them as you need.

**Certificate Authority** (CA), which in this case is represented by `ca.pem` and `ca.key` files, is required to sign server's certificate. Usually it is good "as is" and does not require modification (unless it's close to expiration). If you need to create/re-create CA files, run `create_ca.sh` script.

Remember that certificate files are bound to hostname and/or IP addr: changing hostname/IP can render your certificate useless.

So, if you have CA files (`ca.pem` and `ca.key`), and `certificate.conf` satisfies your requirements, run `create_cert.sh`.

It's supposed to create two files: `service.pem` (certificate) and `service.key` (private key). The names can be different, extensions (?) ... probably too (not sure).


## <a name="scripts"></a>Essential scripts

Usually I store the majority of the scripts in the `./scripts` dir, but in this case it was randomly occupied by the database-related stuff. So, there is also `./scripts2`. What are all those scripts supposed to do?

Scripts in proj root dir (before trying, look inside and modify if necessary):
- `runsvc.sh` runs app/service which is executable file `vote-svc`;
- `runtests.sh` runs all unit tests;
- `runServiceAlpine.sh` runs app/service as a Docker container;

Scripts in `./scripts`:
- `runCQLsvc.sh` runs Docker container with the Apache Cassandra database;
- `runCQLcli.sh` starts Apache Cassandra CQLSH (cmdline client, it comes with Apache Cassandra database image);

Scripts in `./scripts2`:
- `create_ca.sh` is used to create Certificate Authority (CA) certificate `ca.cert` and key `ca.key`; these files are required to sign the server's certificate; you're supposed to run this script once in a year or two, three, five years - depends on the expiration time you choose;
- `create_cert.sh` is used to create self-signed server certificate; it also creates some transient files, but you only need `server.key` and `server.pem`; check `certificate.conf` before running this script, and make sure `ca.key` and `ca.cert` files are present in the same dir;
- `dock_net.sh` (not related to certificates) is used to create a Docker net (one time action, unless you randomly delete it); before creating this [ bridge ] network, you may want to select the appropriate IP address, mask, and modify this script. I prefer private network Class B.


## Just for your information (you can ignore this)

When service starts, it tries to register with [Consul Service Discovery](https://www.consul.io), if it is available. See more about Consul in [Consul Intro](https://developer.hashicorp.com/consul/docs/intro).

Logging, metrics, and tracing are handled by the Go kit packages. Metrics data is supposed to be pulled by
[Prometheus service](https://prometheus.io), tracing is provided by Zipkin, OpenTracing, ...


## Some resources related to this project

[Go-kit / Main Site](https://gokit.io)

[Go-kit / Main Site, Source Code on GitHub](https://github.com/go-kit/kit)

[Go-kit / Main Site, Example on GitHub](https://github.com/go-kit/examples)

[Go-kit / Kit (code generation tool)](https://github.com/kujtimiihoxha/kit)

[How to Parse a JSON Request Body in Go](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)

[Best practices for containerizing Go applications with Docker](https://snyk.io/blog/containerizing-go-applications-with-docker)

[Dockerizing your Go application](https://blog.logrocket.com/dockerizing-go-application)

[Containerizing your Go Applications with Docker - Tutorial](https://tutorialedge.net/golang/go-docker-tutorial)

[Securing gRPC connection with SSL/TLS Certificate using Go](https://medium.com/@mertkimyonsen/securing-grpc-connection-with-ssl-tls-certificate-using-go-db3852fe89dd)

[How to secure gRPC connection with SSL/TLS in Go](https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph)

[Docker Compose Tutorial: advanced Docker made simple](https://www.educative.io/blog/docker-compose-tutorial)
