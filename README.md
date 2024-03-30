# RACOOK

Racook, the cook book of your dreams where the yumiest recipes
can be found and recorded!

## Functionality

- [ ] Able see list of recipes
  - [ ] list has total count
  - [ ] pagination
  - [ ] list accept filter
- [ ] filter
  - [ ] by category
  - [ ] by ingidience
  - [ ] by type?
  - [ ] cooking time
- [ ] find recipe and id
- [ ] record recipes
- [ ] able to register/login as a user
- [ ] mark recipe as favorite
- [ ] like/dislike a recipe

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

migrations up

```bash
make migrations-up
```

migrations down

```bash
make migrations-down
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
