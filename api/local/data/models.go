package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pusher/pusher-http-go"
	"golang.org/x/crypto/bcrypt"
)

const timeout = time.Second * 5

var db *sql.DB

func New(pool *sql.DB) Models {
	db = pool

	return Models{
		User:         User{},
		Token:        Token{},
		Target:       Target{},
		Config:       Config{},
		Reply:        Reply{},
		SavedCommand: SavedCommand{},
	}
}

type Models struct {
	User         User
	Token        Token
	Target       Target
	Config       Config
	Reply        Reply
	SavedCommand SavedCommand
}

// --user
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     Token     `json:"token"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `select id, email, first_name, last_name, password, user_active, created_at, updated_at,
	case
		when (select count(id) from tokens t where user_id = users.id and t.death > NOW()) > 0 then 1
		else 0
	end as has_token
	from users order by last_name`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Token.ID,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`
	var user User

	row := db.QueryRowContext(ctx, q, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		// if a user doesn't exist, the returned error on the frontend is an sql error from the back.
		// create a custom error for right now instead
		myErr := errors.New("incorrect email or password")
		return nil, myErr
	}
	return &user, nil
}

func (u *User) GetUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`
	var user User

	row := db.QueryRowContext(ctx, q, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) PasswordMatches(input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) UpdateUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `update users set
	  email = $1,
	  first_name = $2,
	  last_name = $3,
	  updated_at = $4,
	  user_active = $5
	  where id = $6
	  `
	_, err := db.ExecContext(ctx, q, u.Email, u.FirstName, u.LastName, time.Now(), u.Active, u.ID)

	if err != nil {
		return err
	}
	return nil

}

func (u *User) DeleteUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, q, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) DeleteUserByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CreateUser(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int

	q := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7) returning id
	`

	err = db.QueryRowContext(ctx, q,
		user.Email,
		user.FirstName,
		user.LastName,
		pwd,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (u *User) ResetPwd(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	q := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, q, pwd, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) TestPwd(input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

//__user
//--token

type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Death     time.Time `json:"death"`
	TokenHash []byte    `json:"-"`
}

func (t *Token) GetByToken(input string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `select user_id, email, token, token_hash, created_at, updated_at, death
	from tokens where token = $1`

	var token Token

	row := db.QueryRowContext(ctx, q, input)
	err := row.Scan(
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Death,
	)

	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *Token) GetUserByToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var user User
	q := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`
	row := db.QueryRowContext(ctx, q, token.UserID)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (t *Token) CreateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Death:  time.Now().Add(ttl),
	}

	hotSauce := make([]byte, 16)
	_, err := rand.Read(hotSauce)
	if err != nil {
		return nil, err
	}
	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hotSauce)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]
	return token, nil
}

func (t *Token) TestToken(r *http.Request) (*User, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, errors.New("no auth headers")
	}

	values := strings.Split(header, " ")
	if len(values) != 2 || values[0] != "Bearer" {
		return nil, errors.New("bad auth hheaders")
	}

	data := values[1]
	if len(data) != 26 {
		return nil, errors.New("that token is too small (lol)")
	}

	token, err := t.GetByToken(data)
	if err != nil {
		return nil, errors.New("that token does not exist")
	}

	if token.Death.Before(time.Now()) {
		return nil, errors.New("that token is no longer with us")
	}

	user, err := t.GetUserByToken(*token)
	if err != nil {
		return nil, errors.New("no such user")
	}

	if user.Active == 0 {
		return nil, errors.New("inactive user")
	}

	return user, nil
}

func (t *Token) InsertToken(token Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// purge tokens
	purge := `delete from tokens where user_id = $1`
	_, err := db.ExecContext(ctx, purge, token.UserID)
	if err != nil {
		return err
	}

	token.Email = u.Email
	q := `insert into tokens (user_id, email, token, token_hash, created_at, updated_at, death)
		values ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.ExecContext(ctx, q,
		token.UserID,
		token.Email,
		token.Token,
		token.TokenHash,
		time.Now(),
		time.Now(),
		token.Death,
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) ValidToken(input string) (bool, error) {
	token, err := t.GetByToken(input)
	if err != nil {
		return false, err
	}
	_, err = t.GetUserByToken(*token)
	if err != nil {
		return false, err
	}

	if token.Death.Before(time.Now()) {
		return false, errors.New("dead tokens dont talk")
	}

	return true, nil
}

func (t *Token) DeleteByToken(input string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from tokens where token = $1`

	_, err := db.ExecContext(ctx, q, input)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteTokenByUID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from tokens where user_id = $1`
	_, err := db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

//__token
//--target

type Target struct {
	ID        int       `json:"id"`
	Port      int       `json:"port"`
	Protocol  int       `json:"protocol"`
	Address   string    `json:"address"`
	UserName  string    `json:"username"`
	Key       string    `json:"key"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Target) GetAllTargets() ([]*Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `select id,
	      port,
	      protocol,
	      address,
	      user,
	      key,
	      password,
	      token,
	      created_at,
		  updated_at from targets
		  `
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var targets []*Target

	for rows.Next() {
		var target Target
		err := rows.Scan(
			&target.ID,
			&target.Port,
			&target.Protocol,
			&target.Address,
			&target.UserName,
			&target.Key,
			&target.Password,
			&target.Token,
			&target.CreatedAt,
			&target.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		targets = append(targets, &target)
	}
	return targets, nil
}

func (t *Target) CreateTarget(target Target) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pwd, err := bcrypt.GenerateFromPassword([]byte(target.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int

	q := `insert into targets (port, protocol, address, user, key, password, token, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`
	err = db.QueryRowContext(ctx, q,
		target.Port,
		target.Protocol,
		target.Address,
		target.UserName,
		target.Key,
		pwd,
		target.Token,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (t *Target) ModifyTargetPwd(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	q := `update targets set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, q, pwd, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Target) GetTarget(id int) (*Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	q := `select id,
	      port,
	      protocol,
	      address,
	      user,
	      key,
	      password,
	      token,
	      created_at,
		  updated_at from targets
		  `
	var target Target

	row := db.QueryRowContext(ctx, q, id)
	err := row.Scan(
		&target.ID,
		&target.Port,
		&target.Protocol,
		&target.Address,
		&target.UserName,
		&target.Key,
		&target.Password,
		&target.Token,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (t *Target) UpdateTarget() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	q := `update targets set
		port = $1,
		protocol = $2,
		address = $3,
		user = $4,
		key = $5,
		password = $6,
		token = $7,
		updated_at = $8 where id = $9`
	_, err := db.ExecContext(ctx, q, t.Port, t.Protocol, t.Address, t.UserName, t.Key, t.Password, t.Token, time.Now(), t.ID)
	if err != nil {
		return err
	}
	return nil

}

func (t *Target) DeleteTarget() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from targets where id = $1`
	_, err := db.ExecContext(ctx, q, t.ID)
	if err != nil {
		return err
	}
	return nil
}

//__target
//--config

// type Config struct {
// 	ID        int       `json:"id"`
// 	Name      string    `json:"name"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// 	Plan      string    `json:"plan"`
// }

type Config struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	User    string `json:"user"`
	Key     string `json:"key"`
	LogPath string `json:"logpath"`
	Hosts   string `json:"hosts"`
	Command string `json:"command"`
	Timeout string `json:"timeout"`
	Port    string `json:"port"`
	Fatal   bool   `json:"fatal"`
	Ordered bool   `json:"ordered"`
}

type Plan struct {
	Name  string        `json:"name"`
	Steps []Instruction `json:"steps"`
}

type Instruction struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Hosts   []Target  `json:"hosts"`
	Command string    `json:"command"`
	Fatal   bool      `json:"fatal"`
	Timeout int       `json:"timeout"`
}

func (c *Config) GetConfigs() ([]*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `
	  select
	    id,
		username,
		key,
		logpath,
		hosts,
		command,
		timeout,
		port,
		fatal,
		ordered,
		name
	  from configurations`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var configs []*Config

	for rows.Next() {
		var cfg Config

		err := rows.Scan(
			&cfg.ID,
			&cfg.User,
			&cfg.Key,
			&cfg.LogPath,
			&cfg.Hosts,
			&cfg.Command,
			&cfg.Timeout,
			&cfg.Port,
			&cfg.Fatal,
			&cfg.Ordered,
			&cfg.Name,
			// pq.Array(&cfg.Plan),
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &cfg)
	}
	return configs, nil
}

func (c *Config) GetConfigByID(id int) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `
	  select
	    id,
		username,
		key,
		logpath,
		hosts,
		command,
		timeout,
		port,
		fatal,
		ordered,
		name
	  from configurations
	  where id = $1`

	var cfg Config
	var tmpHosts string

	row := db.QueryRowContext(ctx, q, id)
	err := row.Scan(
		&cfg.ID,
		&cfg.User,
		&cfg.Key,
		&cfg.LogPath,
		&tmpHosts,
		&cfg.Command,
		&cfg.Timeout,
		&cfg.Port,
		&cfg.Fatal,
		&cfg.Ordered,
		&cfg.Name,
	)
	if err != nil {
		return nil, err
	}
	cfg.Hosts = strings.ReplaceAll(strings.Trim(tmpHosts, "{}"), ",", " ")

	return &cfg, nil
}

func (c *Config) UpdateConfig() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	hosts := strings.ReplaceAll(fmt.Sprintf("{%v}", c.Hosts), " ", ",")

	q := `update configurations set
		username = $1,
		key = $2,
		logpath = $3,
		hosts = $4,
		command = $5,
		timeout = $6,
		port = $7,
		fatal = $8,
		ordered = $9,
		name = $10
		where id = $11`
	_, err := db.ExecContext(
		ctx,
		q,
		c.User,
		c.Key,
		c.LogPath,
		hosts,
		c.Command,
		c.Timeout,
		c.Port,
		c.Fatal,
		c.Ordered,
		c.Name,
		c.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeleteConfig(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `delete from configurations where id = $1`

	_, err := db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) CreateConfig(cfg Config) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	hosts := strings.ReplaceAll(fmt.Sprintf("{%v}", cfg.Hosts), " ", ",")

	var newID int
	q := `
	  insert into configurations
		(username, key, logpath, hosts, command, timeout, port, fatal, ordered, name)
	  values
	  	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	  returning id`
	err := db.QueryRowContext(
		ctx,
		q,
		cfg.User,
		cfg.Key,
		cfg.LogPath,
		hosts,
		cfg.Command,
		cfg.Timeout,
		cfg.Port,
		cfg.Fatal,
		cfg.Ordered,
		cfg.Name).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// you changed config to string from Config
type Reply struct {
	ID      int       `json:"id"`
	Reply   string    `json:"reply"`
	TimeTX  time.Time `json:"time_tx"`
	TimeRX  time.Time `json:"time_rx"`
	Config  string    `json:"config"`
	Host    string    `json:"host"`
	Good    bool      `json:"good"`
	ReplyTo string    `json:"reply_to"`
}

func (res *Reply) GetResponses(limit int) ([]*Reply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `
	  select id,
	  reply,
	  command_sent,
	  reply_received,
	  config,
	  host,
	  good,
	  reply_to from replies
	  order by id desc
	  limit %v`
	rows, err := db.QueryContext(ctx, fmt.Sprintf(q, limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var replies []*Reply

	for rows.Next() {
		var reply Reply
		err := rows.Scan(
			&reply.ID,
			&reply.Reply,
			&reply.TimeTX,
			&reply.TimeRX,
			&reply.Config,
			&reply.Host,
			&reply.Good,
			&reply.ReplyTo,
		)

		if err != nil {
			return nil, err
		}
		replies = append(replies, &reply)
	}
	return replies, nil
}

func (res *Reply) InsertResponse(response Reply) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var newID int

	if len(response.ReplyTo) > 59 {
		response.ReplyTo = response.ReplyTo[:59]
	}

	q := `insert into replies (command_sent, reply_received, reply, config, good, host, reply_to)
	values ($1, $2, $3, $4, $5, $6, $7) returning id
	`
	fmt.Println(response)
	err := db.QueryRowContext(ctx, q,
		response.TimeTX,
		response.TimeRX,
		response.Reply,
		response.Config,
		response.Good,
		response.Host,
		response.ReplyTo,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (res *Reply) GetResponse(id int) (*Reply, error) {
	var reply Reply
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	q := `
	  select id,
	  reply,
	  command_sent,
	  reply_received,
	  config,
	  host,
	  good,
	  reply_to from replies
	  where id = $1`
	row := db.QueryRowContext(ctx, q, id)
	err := row.Scan(
		&reply.ID,
		&reply.Reply,
		&reply.TimeTX,
		&reply.TimeRX,
		&reply.Config,
		&reply.Host,
		&reply.Good,
		&reply.ReplyTo,
	)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

type SavedCommand struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
	Command     string    `json:"command"`
	CommandName string    `json:"command_name"`
}

func (c *SavedCommand) InsertCommand(cmd SavedCommand) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var newID int
	q := `insert into saved_commands (created_at, user_id, command, name)
	values ($1, $2, $3, $4) returning id
	`

	err := db.QueryRowContext(ctx, q,
		time.Now(),
		cmd.UserID,
		cmd.Command,
		cmd.CommandName,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (c *SavedCommand) GetSavedCommands(limit, id int) ([]*SavedCommand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	q := `
	  select id, created_at, user_id, command, name
	  from saved_commands where user_id = %v
	  order by id desc
	  limit %v`
	rows, err := db.QueryContext(ctx, fmt.Sprintf(q, id, limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var commands []*SavedCommand

	for rows.Next() {
		var command SavedCommand
		err := rows.Scan(
			&command.ID,
			&command.CreatedAt,
			&command.UserID,
			&command.Command,
			&command.CommandName,
		)

		if err != nil {
			return nil, err
		}
		commands = append(commands, &command)
	}
	return commands, nil
}

// websocket stuff

type WSClient interface {
	Trigger(channel string, eventName string, data interface{}) error
	TriggerMulti(channels []string, eventName string, data interface{}) error
	TriggerExclusive(channel string, eventName string, data interface{}, socketID string) error
	TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) error
	TriggerBatch(batch []pusher.Event) error
	Channels(additionalQueries map[string]string) (*pusher.ChannelsList, error)
	Channel(name string, additionalQueries map[string]string) (*pusher.Channel, error)
	GetChannelUsers(name string) (*pusher.Users, error)
	AuthenticatePrivateChannel(params []byte) (response []byte, err error)
	AuthenticatePresenceChannel(params []byte, member pusher.MemberData) (response []byte, err error)
	Webhook(header http.Header, body []byte) (*pusher.Webhook, error)
}
