package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/delonce/socialnetwork/internal/service"
	"github.com/delonce/socialnetwork/internal/service/user"

	"github.com/julienschmidt/httprouter"
)

func (handler *NetworkHandler) GetIndexPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.HandlerLogger.Infof("Connection from %s to %s, method=%s", r.RemoteAddr, r.RequestURI, r.Method)

	err := INDEX_TEMPLATE.Execute(w, nil)

	if err != nil {
		handler.HandlerLogger.Errorf("Failed connection with error: %v", err)
		if _, err = w.Write([]byte("Something wrong...")); err != nil {
			panic(err)
		}
	}
}

func (handler *NetworkHandler) GetRegisterPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.HandlerLogger.Infof("Connection from %s to %s, method=%s", r.RemoteAddr, r.RequestURI, r.Method)

	err := REGISTER_TEMPLATE.Execute(w, nil)

	if err != nil {
		handler.HandlerLogger.Errorf("Failed connection with error: %v", err)
		if _, err = w.Write([]byte("Something wrong...")); err != nil {
			panic(err)
		}
	}
}

func (handler *NetworkHandler) GoRegisterUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.HandlerLogger.Infof("Connection from %s to %s, method=%s", r.RemoteAddr, r.RequestURI, r.Method)

	username := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")

	regService, err := user.NewRegisterService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		w.Write([]byte("Something wrong..."))
	}

	_, err = regService.RegisterNewUser(username, password, email)

	if err != nil {

		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, LOGIN_URL, http.StatusSeeOther)
}

func (handler *NetworkHandler) SignIn(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.HandlerLogger.Infof("Connection from %s to %s, method=%s", r.RemoteAddr, r.RequestURI, r.Method)

	username := r.FormValue("login")
	password := r.FormValue("password")

	authService, err := user.NewAuthService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		w.Write([]byte("Something wrong..."))
	}

	session, err := authService.CreateNewSession(username, password)

	if err != nil {
		// TODO Error msg on html page
		w.Write([]byte(err.Error()))
		return
	}

	setTokenCookie(w, session)
	http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
}

func (handler *NetworkHandler) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accessToken, _ := r.Cookie(accessTokenCookie)
	refreshToken, _ := r.Cookie(refreshTokenCookie)

	authService, err := user.NewAuthService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Errorf("Failed to create auth service when logouting")
	}

	authService.Logout(accessToken.Value, refreshToken.Value)
	clearCookie(w)

	http.Redirect(w, r, LOGIN_URL, http.StatusSeeOther)
}

func (handler *NetworkHandler) GetLoginPage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	handler.HandlerLogger.Infof("Connection from %s to %s, method=%s", r.RemoteAddr, r.RequestURI, r.Method)

	ok, session := handler.checkUserSession(w, r)

	if ok {
		setTokenCookie(w, session)

		userID := session.GetTokenPair().Refresh.UserID

		ctx := context.WithValue(context.Background(), userIDKey, userID)
		r = r.WithContext(ctx)

		http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
		return
	}

	err := LOGIN_TEMPLATE.Execute(w, nil)

	if err != nil {
		handler.HandlerLogger.Errorf("Failed connection with error: %v", err)
		if _, err = w.Write([]byte("Something wrong...")); err != nil {
			panic(err)
		}
	}
}

func (handler *NetworkHandler) CheckAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ok, session := handler.checkUserSession(w, r)

		if !ok {
			http.Redirect(w, r, LOGIN_URL, http.StatusSeeOther)
			return
		}

		setTokenCookie(w, session)

		userID := session.GetTokenPair().Refresh.UserID

		ctx := context.WithValue(context.Background(), userIDKey, userID)
		r = r.WithContext(ctx)

		next(w, r, params)
	}
}

func (handler *NetworkHandler) checkUserSession(w http.ResponseWriter, r *http.Request) (bool, user.Session) {
	accessToken, err := r.Cookie(accessTokenCookie)

	if err != nil {
		return false, nil
	}

	refreshToken, err := r.Cookie(refreshTokenCookie)

	if err != nil {
		return false, nil
	}

	authServ, err := user.NewAuthService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		handler.HandlerLogger.Error("Error when creating authService")
		w.Write([]byte("Something wrong"))
		return false, nil
	}

	return authServ.CheckSession(accessToken.Value, refreshToken.Value)
}

func setTokenCookie(w http.ResponseWriter, session user.Session) {
	accessToken := session.GetTokenPair().Access
	refreshToken := session.GetTokenPair().Refresh

	w.Header().Set("Access-Control-Allow-Credentials", "true")

	accessCookie := http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	refreshCookie := http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken.UUID,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func clearCookie(w http.ResponseWriter) {
	accessCookie := http.Cookie{
		Name:   accessTokenCookie,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	refreshCookie := http.Cookie{
		Name:   refreshTokenCookie,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func (handler *NetworkHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) *service.User {
	userID := fmt.Sprintf("%v", r.Context().Value(userIDKey))

	authService, err := user.NewAuthService(handler.HandlerLogger, handler.HandlerConfig)

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Something wrong"))
		return nil
	}

	user, err := authService.GetUserByID(userID)

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, LOGIN_URL, http.StatusSeeOther)
		return nil
	}

	return user
}
