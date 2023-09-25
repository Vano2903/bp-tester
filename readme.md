# breapoint tester

## run

```
docker-compose up
```

## todo (in order of priority)

- [ ] rebuild queue at the start by checking for pendings
- [ ] add language support (golang, binary)
- [ ] add login support (username, password) to keep track of attempts
- [ ] add 404 page
- [ ] add leaderboard
- [ ] add page to show all attempts (it should update in real time and be paginated)
- [ ] set attempt output to docker build output when image fails to build
- [ ] add function to create image with attempt status (to think)

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
