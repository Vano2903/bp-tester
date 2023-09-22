# breapoint tester

## run

```
docker-compose up
```

## todo

- [ ] rebuild queue at the start by checking for pendings
- [ ] set attempt output to docker build output when image fails to build
- [ ] add language support (golang, binary)
- [ ] add login support (username, password) to keep track of attempts

## stuff

entità:

utente:

- id
- telegramName
- telegramID
- tentativi []tentativo

tentativo:

- id
- status => pending, invalid, building, running, buildError, runError, fail, success
- output => output del programma && status code
- tempi []tempo
- best
- avg
- fileContent

leaderboard:
