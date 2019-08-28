# Rest API example in Golang. Dockerized.

Specification [details](https://github.com/kodersky/golang-api-example/specification.md).


## Usage

1. Create `config.yaml` and `.env` files. Use `config.yaml.example` and `env.example`
as a templates.

2. Make start.sh file executable `$ sudo chmod +x start.sh`.

3. Run `$ ./start.sh`.

## Testing

Command `$ ./start.sh` will build small docker container with Nginx as a reverse
proxy. `go` command is not available there.

If you want to run `go` commands inside container like for example for testing 
`docker exec -it golang-api-example go test  ./...` please use 
`docker-compose-dev.yml` file.

Run: 
1. `$ docker-compose -f docker-compose-dev.yml build`
2. `$ docker-compose -f docker-compose-dev.yml up`

Don't forget to run:

 `$ docker exec -i golang-example-api-db mysql -uroot -p"${password}" orders < db.sql` 

if you haven't run it before.


## Troubleshoot:

If any problems with Docker try to use Edge version. Setup was tested only 
on MacOS with Docker 2.1.1.0 Edge version.

Make sure you have `config.yaml` and `.env` in project root directory.
Make sure you have **valid API KEY for Google Maps** placed in `config.yaml`.

Use `config.yaml.example` and and `env.example` as templates.