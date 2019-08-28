# Rest API example in Golang. Dockerized.

Specification [details](https://github.com/kodersky/golang-api-example/blob/master/specifcation.md).


## Usage

1. Create `config.yaml` and `.env` files. Use `config.yaml.example` and `env.example`
as a templates.

2. Make start.sh file executable `$ sudo chmod +x start.sh`.

3. Run `$ ./start.sh`.

## Testing

Command `$ ./start.sh` will build **small** docker container optimized for production, 
therefore `go` command is not available there.

Please use Golang on your host machine.

## Database

You can connect to DB (MySQL) from your host machine on port `33306`

## Troubleshoot:

If any problems with Docker try to use Edge version. Setup was tested only 
on MacOS with Docker 2.1.1.0 Edge version.

Make sure you have `config.yaml` and `.env` in project root directory.
Make sure you have **valid API KEY for Google Maps** placed in `config.yaml`.

Use `config.yaml.example` and and `env.example` as templates.