# breapoint tester

## run

```
docker-compose up
```

## part for vano

todo:

- controlla che login e register funzionino
- aggiungi pagine di login e register
- pagina di login|register deve reindirizzare alla home se l'utente è già loggato
- implementa logout (rimuovi access e refresh token)
- pagina dell'utente (sarà la nuova homepage per l'utente)
  - mostra tutti i tentativi fatti (paginata e aggiornata in real time)
  - bottone con cui fare un nuovo tentativo
- l'homepage per gli utenti non loggati sarà la leaderboard (/leaderboard)
- gli utenti non loggati possono fare tenativi ma non saranno considerati nella leaderboard

## todo (in order of priority)

- [ ] routine to clean db from expired tokens
- [ ] check if the unique constraint on user db works correctly
- [x] rebuild queue at the start by checking for pendings or attempts in status building
- [ ] add login support (username, password) to keep track of attempts
- [ ] add language support (golang, binary)
- [ ] add 404 page
- [ ] add leaderboard
- [ ] add page to show all attempts (it should update in real time and be paginated)
- [ ] set attempt output to docker build output when image fails to build
- [ ] add function to create image for social (like the image that comes with a gh link) with attempt status (to think)

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
