# spoonfed-go

Go adaptation of my spoonfed-gql.

## Roadmap

- [x] Migrate base files
  - Added restructured cursor pagination for recipes
- [x] Add Docker
- [x] Add live reload(Air)
- [x] Replaced MongoDB with Gorm(Postgres)
- [x] Migrate dataloaders
  - [x] Add dataloaders for likes
- [ ] Migrate over the rest of the schema
- [ ] Migrate AWS S3
- [ ] Add recipe folder functionality

## Usage

### To start a development server:

- Create a .env file(See `.env.example` for required fields)
- Run `make dev-dc-build` to create the docker image
- Run `make dev-dc-run` to spin up the docker container
  - PS: On initial load or if you require a rebuild you can just run `make dc-up` to build and run in one command
- Then viola pop into http://localhost:8080/ to access the playground

