package controller

import (
	"bytes"
	"context"
	"crypto/sha512"

	"github.com/vano2903/bp-tester/model"
)

func (c *Controller) hashPassword(password string) string {
	hash := sha512.Sum512([]byte(password))
	return bytes.NewBuffer(hash[:]).String()
}

func (c *Controller) newUserModel(username, password string) *model.User {
	return &model.User{
		Username: username,
		Password: c.hashPassword(password),
	}
}

func (c *Controller) CreateUser(ctx context.Context, username, password string) (*model.User, error) {
	user := c.newUserModel(username, password)
	return user, c.userRepo.InsertOne(ctx, user)
}

func (c *Controller) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := c.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if !c.AreUserCredentialsCorrect(user, password) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (c *Controller) AreUserCredentialsCorrect(user *model.User, password string) bool {
	hashedPassword := c.hashPassword(password)
	return user.Password == hashedPassword
}

/*
pagina di login|register
redirect alla homepage se il cookie access token è settato e valido
se non lo è mostra pagina di login

login-> username password
prendi user a username
controlla password
se corretto genera tokens
se non corretto messaggio error

register -> username password
controlla che username non sia già preso
genera utente
genera token dell'id dell'utente
ritorna tokens e info utente
*/
