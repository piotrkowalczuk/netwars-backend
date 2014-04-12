package main

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(dbPool *sql.DB) (repository *UserRepository) {
	repository = &UserRepository{dbPool}

	return
}

func (ur *UserRepository) Insert(user *User) (sql.Result, error) {
	query := `
		INSERT INTO users (user_name, user_pass, pass_type, pass_salt, last_login, bad_logins, email, ntcnick, nickhistory, user_status, change_date, change_user_id, change_ip, email_used, referrer, gg, extrainfo, created, trial, showemail, refer_count, suspended)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
		RETURNING user_id
	`
	result, err := ur.db.Exec(
		query,
		&user.Name,
		&user.Password,
		&user.PasswordType,
		&user.PasswordSalt,
		&user.LoginAt,
		&user.BadLogins,
		&user.Email,
		&user.NTCNick,
		&user.NickHistory,
		&user.Status,
		&user.ChangeAt,
		&user.ChangeUserId,
		&user.ChangeIp,
		&user.EmailUsed,
		&user.ReferrerId,
		&user.GaduGadu,
		&user.ExtraInfo,
		&user.CreatedAt,
		&user.Trial,
		&user.ShowEmail,
		&user.NbOfRefs,
		&user.Suspended,
	)

	return result, err
}

func (ur *UserRepository) Update(user *User) (sql.Result, error) {
	query := `
		UPDATE users
		SET user_name = $1,
			user_pass = $2,
			pass_type = $3,
			pass_salt = $4,
			last_login = $5,
			bad_logins = $6,
			email = $7,
			ntcnick = $8,
			nickhistory = $9,
			user_status = $10,
			change_date = $11,
			change_user_id = $12,
			change_ip = $13,
			email_used = $14,
			referrer = $15,
			gg = $16,
			extrainfo = $17,
			created = $18,
			trial = $19,
			showemail = $20,
			refer_count = $21,
			suspended = $22
		WHERE user_id = $23
	`
	result, err := ur.db.Exec(
		query,
		&user.Name,
		&user.Password,
		&user.PasswordType,
		&user.PasswordSalt,
		&user.LoginAt,
		&user.BadLogins,
		&user.Email,
		&user.NTCNick,
		&user.NickHistory,
		&user.Status,
		&user.ChangeAt,
		&user.ChangeUserId,
		&user.ChangeIp,
		&user.EmailUsed,
		&user.ReferrerId,
		&user.GaduGadu,
		&user.ExtraInfo,
		&user.CreatedAt,
		&user.Trial,
		&user.ShowEmail,
		&user.NbOfRefs,
		&user.Suspended,
		&user.Id,
	)

	return result, err
}

func (ur *UserRepository) FindOne(id int64) (*User, error) {
	user := new(User)

	err := ur.db.QueryRow(
		"SELECT * FROM users as u WHERE u.user_id = $1",
		id,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.PasswordType,
		&user.PasswordSalt,
		&user.LoginAt,
		&user.BadLogins,
		&user.Email,
		&user.NTCNick,
		&user.NickHistory,
		&user.Status,
		&user.ChangeAt,
		&user.ChangeUserId,
		&user.ChangeIp,
		&user.EmailUsed,
		&user.ReferrerId,
		&user.GaduGadu,
		&user.ExtraInfo,
		&user.CreatedAt,
		&user.Trial,
		&user.ShowEmail,
		&user.NbOfRefs,
		&user.Suspended,
	)

	return user, err
}

func (ur *UserRepository) FindOneByEmailAndPassword(email string, password string) (error, *User) {
	user := new(User)

	err := ur.db.QueryRow(
		"SELECT * FROM users as u WHERE u.email = $1 AND u.user_pass = $2",
		email,
		password,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.PasswordType,
		&user.PasswordSalt,
		&user.LoginAt,
		&user.BadLogins,
		&user.Email,
		&user.NTCNick,
		&user.NickHistory,
		&user.Status,
		&user.ChangeAt,
		&user.ChangeUserId,
		&user.ChangeIp,
		&user.EmailUsed,
		&user.ReferrerId,
		&user.GaduGadu,
		&user.ExtraInfo,
		&user.CreatedAt,
		&user.Trial,
		&user.ShowEmail,
		&user.NbOfRefs,
		&user.Suspended,
	)

	return err, user
}
