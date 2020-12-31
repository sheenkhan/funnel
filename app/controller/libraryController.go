package controller

import (
	"encoding/json"
	"funnel/app/errors"
	"funnel/app/helps"
	"funnel/app/model"
	"funnel/app/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var system = service.LibrarySystem{}

func LibraryLogin(context *gin.Context) {
	err := libraryLoginHandle(context)
	if err == nil {
		helps.ContextDataResponseJson(context, helps.SuccessResponseJson(nil))
	}
}

func LibraryBorrowHistory(context *gin.Context) {
	session := sessions.Default(context)
	libraryJson := session.Get("library").([]byte)
	if string(libraryJson) == "{}" {
		helps.ContextDataResponseJson(context, helps.FailResponseJson(errors.NotLogin, nil))
		return
	}

	user := model.LibraryUser{}
	_ = json.Unmarshal(libraryJson, &user)
	books := system.GetBorrowHistory(&user)

	helps.ContextDataResponseJson(context, helps.SuccessResponseJson(books))
}

func LibraryCurrentBorrow(context *gin.Context) {
	session := sessions.Default(context)
	libraryJson := session.Get("library").([]byte)
	if string(libraryJson) == "{}" {
		helps.ContextDataResponseJson(context, helps.FailResponseJson(errors.NotLogin, nil))
		return
	}

	user := model.LibraryUser{}
	_ = json.Unmarshal(libraryJson, &user)
	books := system.GetCurrentBorrow(&user)

	helps.ContextDataResponseJson(context, helps.SuccessResponseJson(books))
}

func libraryLoginHandle(context *gin.Context) error {
	isValid := helps.CheckPostFormEmpty(
		context,
		[]string{"username", "password"},
	)

	if !isValid {
		helps.ContextDataResponseJson(context, helps.FailResponseJson(errors.RequestFailed, nil))
		return errors.ERR_INVALID_ARGS
	}

	user := model.LibraryUser{Username: context.PostForm("username"), Password: context.PostForm("password")}
	err := system.Login(&user)

	if err == errors.ERR_WRONG_PASSWORD {
		helps.ContextDataResponseJson(context, helps.FailResponseJson(errors.WrongPassword, nil))
		return err
	}
	if err != nil {
		helps.ContextDataResponseJson(context, helps.FailResponseJson(errors.UnKnown, nil))
		return err
	}

	session := sessions.Default(context)
	libraryJson, _ := json.Marshal(user)
	session.Set("library", libraryJson)
	_ = session.Save()

	return nil
}