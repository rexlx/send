package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pusher/pusher-http-go"
	"github.com/rexlx/vapi/local/data"
	"github.com/rexlx/vapi/local/utils"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type envelope map[string]interface{}

var upgradeConn = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (app *settings) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var data jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		data.Error = true
		data.Message = "invalid json"
		_ = app.writeJSON(w, http.StatusBadRequest, data)
	}

	user, err := app.models.User.GetByEmail(creds.Email)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	app.infoLog.Println("login attempt for...", user.Email)
	validPwd, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPwd {
		app.errorJSON(w, errors.New("sorry, that didn't work"))
		app.errorLog.Println("failed login attempt for", creds.Email)
		return
	}

	// is user active
	if user.Active == 0 {
		app.errorJSON(w, errors.New("inactive user"))
		return
	}

	token, err := app.models.Token.CreateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	app.infoLog.Println("successful login for", creds.Email)
	err = app.models.Token.InsertToken(*token, *user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data = jsonResponse{
		Error:   false,
		Message: "you're in...",
		Data:    envelope{"token": token, "user": user},
	}

	err = app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *settings) Logout(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token string `json:"token"`
		User  string `json:"user"`
	}

	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, errors.New("bad json"))
		return
	}

	app.infoLog.Printf("logging %v out...", request.User)
	err = app.models.Token.DeleteByToken(request.Token)
	if err != nil {
		app.errorJSON(w, errors.New("bad json"))
		return
	}

	data := jsonResponse{
		Error:   false,
		Message: "logged out",
	}
	_ = app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := jsonResponse{
		Error:   false,
		Message: "boom",
		Data:    envelope{"users": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) EditUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if user.ID == 0 {
		if _, err := app.models.User.CreateUser(user); err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		u, err := app.models.User.GetUser(user.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Active = user.Active

		if err := u.UpdateUser(); err != nil {
			app.errorJSON(w, err)
			return
		}
		if user.Password != "" {
			err := u.ResetPwd(user.Password)
			if err != nil {
				app.errorJSON(w, err)
				return
			}
		}
	}

	data := jsonResponse{
		Error:   false,
		Message: "ok",
	}
	_ = app.writeJSON(w, http.StatusAccepted, data)
}

func (app *settings) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	user, err := app.models.User.GetUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)
}

func (app *settings) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID int `json:"id"`
	}
	err := app.readJSON(w, r, data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.User.DeleteUserByID(data.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "user deleted",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *settings) BootUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	user, err := app.models.User.GetUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Active = 0
	err = user.UpdateUser()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// delete tokens
	err = app.models.Token.DeleteTokenByUID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := jsonResponse{
		Error:   false,
		Message: "user logged out",
	}
	_ = app.writeJSON(w, http.StatusAccepted, data)
}

func (app *settings) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	valid := false
	valid, _ = app.models.Token.ValidToken(data.Token)
	payload := jsonResponse{
		Error: false,
		Data:  valid,
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *settings) AllTargets(w http.ResponseWriter, r *http.Request) {
	var targets data.Target
	all, err := targets.GetAllTargets()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := jsonResponse{
		Error:   false,
		Message: "boom",
		Data:    envelope{"targets": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) GetTarget(w http.ResponseWriter, r *http.Request) {
	targetID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	target, err := app.models.Target.GetTarget(targetID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, target)
}

func (app *settings) DeleteTarget(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID int `json:"id"`
	}
	err := app.readJSON(w, r, data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Target.DeleteTarget()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "user deleted",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *settings) EditTarget(w http.ResponseWriter, r *http.Request) {
	var target data.Target
	err := app.readJSON(w, r, &target)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if target.ID == 0 {
		if _, err := app.models.Target.CreateTarget(target); err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		t, err := app.models.Target.GetTarget(target.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		t.Address = target.Address
		t.UserName = target.UserName
		t.Port = target.Port
		t.Protocol = target.Protocol
		t.Password = target.Password
		t.Token = target.Token
		t.UserName = target.UserName
		t.UpdatedAt = target.UpdatedAt
		if err := t.UpdateTarget(); err != nil {
			app.errorJSON(w, err)
			return
		}
		if target.Password != "" {
			err := t.ModifyTargetPwd(target.Password)
			if err != nil {
				app.errorJSON(w, err)
				return
			}
		}
	}

	data := jsonResponse{
		Error:   false,
		Message: "ok",
	}
	_ = app.writeJSON(w, http.StatusAccepted, data)
}

// func (app *settings) DisableTarget(w http.ResponseWriter, r *http.Request) {
// 	targetID, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}
// 	target, err := app.models.Target.GetTarget(targetID)
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	target.Active = 0
// 	err = target.UpdateTarget()
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	data := jsonResponse{
// 		Error:   false,
// 		Message: "target is disabled",
// 	}
// 	_ = app.writeJSON(w, http.StatusAccepted, data)
// }

func (app *settings) AllConfigs(w http.ResponseWriter, r *http.Request) {
	var configs data.Config
	all, err := configs.GetConfigs()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := jsonResponse{
		Error:   false,
		Message: "boom",
		Data:    envelope{"configs": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) GetConfig(w http.ResponseWriter, r *http.Request) {
	configID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	target, err := app.models.Config.GetConfigByID(configID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, target)
}

func (app *settings) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID int `json:"id"`
	}
	err := app.readJSON(w, r, data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Config.DeleteConfig(data.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "configuration deleted",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *settings) EditConfig(w http.ResponseWriter, r *http.Request) {
	var config data.Config
	err := app.readJSON(w, r, &config)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if config.ID == 0 {
		if _, err := app.models.Config.CreateConfig(config); err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		c, err := app.models.Config.GetConfigByID(config.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		c.User = config.User
		c.Key = config.Key
		c.Port = config.Port
		c.Hosts = config.Hosts
		c.Command = config.Command
		c.Timeout = config.Timeout
		c.Fatal = config.Fatal
		c.Ordered = config.Ordered
		c.LogPath = config.LogPath
		c.Name = config.Name
		if err := c.UpdateConfig(); err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	data := jsonResponse{
		Error:   false,
		Message: "config was saved",
	}
	_ = app.writeJSON(w, http.StatusAccepted, data)
}

func (app *settings) SendCommand(w http.ResponseWriter, r *http.Request) {
	response := make(chan utils.GenericResponse)
	var config data.Config
	var reply data.Reply
	type msgContainer struct {
		Message string `json:"message"`
	}

	err := app.readJSON(w, r, &config)
	app.infoLog.Println(config)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "instructions sent",
		Data:    envelope{"config": config},
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
	// now do the work.
	for _, dest := range strings.Split(config.Hosts, " ") {
		if config.Ordered {
			response := utils.SendSSHCmd(dest, config.User, &config, time.Now())
			app.infoLog.Println(response, "ORDERED FEAT NOT IMPL")
		} else {
			go func(dest, user string) {
				response <- utils.SendSSHCmd(dest, config.User, &config, time.Now())
			}(dest, config.User)
		}
	}
	if !config.Ordered {
		timeout, err := strconv.Atoi(config.Timeout)
		if err != nil {
			app.errorJSON(w, err)
		}
		for i := 0; i < len(strings.Split(config.Hosts, " ")); i++ {
			select {
			case results := <-response:
				cfg, _ := json.Marshal(results.Config)
				message := msgContainer{
					Message: results.Message,
					// Message: strings.TrimSpace(results.Message),
				}
				msg, _ := json.Marshal(message)
				reply.Good = results.Good
				reply.Config = string(cfg)
				reply.TimeTX = results.TimeTX
				reply.TimeRX = results.TimeRX
				reply.Host = results.Host
				reply.Reply = string(msg) // dont store as string!
				_, err := app.models.Reply.InsertResponse(reply)
				if err != nil {
					app.errorLog.Println(err)
				}
			case <-time.After(time.Duration(timeout) * time.Second):
				cfg, _ := json.Marshal(config)
				message := msgContainer{
					Message: fmt.Sprintf("configuration timeout reached: %v", timeout),
					// Message: strings.TrimSpace(results.Message),
				}
				msg, _ := json.Marshal(message)
				app.infoLog.Printf("backend  TIMEOUT")
				reply.Good = false
				reply.TimeRX = time.Now()
				reply.Config = string(cfg)
				reply.Reply = string(msg)
				reply.Host = "system"
				_, err := app.models.Reply.InsertResponse(reply)
				if err != nil {
					app.errorLog.Println(err)
				}
				return
			}
		}
	}
}

func (app *settings) Last24responses(w http.ResponseWriter, r *http.Request) {
	var replies data.Reply
	all, err := replies.GetResponses(24)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := jsonResponse{
		Error:   false,
		Message: "bAm!",
		Data:    envelope{"data": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) GetResponses(w http.ResponseWriter, r *http.Request) {
	resPerPage, err := strconv.Atoi(chi.URLParam(r, "num"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var replies data.Reply
	all, err := replies.GetResponses(resPerPage)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := jsonResponse{
		Error:   false,
		Message: "got responses",
		Data:    envelope{"data": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

func (app *settings) GetResponse(w http.ResponseWriter, r *http.Request) {
	responseID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	res, err := app.models.Reply.GetResponse(responseID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, res)
}

func (app *settings) SaveCommand(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		UserID      int    `json:"user_id"`
		Command     string `json:"command"`
		CommandName string `json:"command_name"`
	}
	var p payload
	var c data.SavedCommand

	err := app.readJSON(w, r, &p)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	c.UserID = p.UserID
	c.Command = p.Command
	c.CommandName = p.CommandName
	_, err = app.models.SavedCommand.InsertCommand(c)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := jsonResponse{
		Error:   false,
		Message: "command saved",
	}
	_ = app.writeJSON(w, http.StatusOK, msg)
}

func (app *settings) GetUserSavedCommands(w http.ResponseWriter, r *http.Request) {
	var command data.SavedCommand
	type uid struct {
		UserID int `json:"user_id"`
	}
	var u uid
	err := app.readJSON(w, r, &u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	all, err := command.GetSavedCommands(30, u.UserID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := jsonResponse{
		Error:   false,
		Message: "saved commands retrieved",
		Data:    envelope{"data": all},
	}

	app.writeJSON(w, http.StatusOK, data)
}

// web socket stuff
// type wsRes struct {
// 	Action      string `json:"action"`
// 	Message     string `json:"message"`
// 	MessageType string `json:"message_type"`
// }

func (app *settings) IpeAuth(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	u, err := app.models.User.GetUser(user.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(u.ID),
		UserInfo: map[string]string{
			"name": u.FirstName,
			"id":   strconv.Itoa(u.ID),
		},
	}

	res, err := app.WsClient.AuthenticatePresenceChannel(body, presenceData)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, res)

}
