package user

import (
	"database/sql"
	"time"
)

type User struct {
	Id sql.NullInt64 `db:"user_id"`
	Name string `db:"user_name"`
	Password string `db:"user_pass"`
	PasswordSalt string `db:"pass_salt"`
	PasswordType uint16 `db:"pass_type"`
	BadLogins sql.NullInt64 `db:"bad_logins"`
	Email sql.NullString `db:"email"`
	NTCNick sql.NullString `db:"ntcnick"`
	NickHistory sql.NullString `db:"nickhistory"`
	Status *uint16 `db:"user_status"`
	EmailUsed sql.NullInt64`db:"email_used"`
	ReferrerId sql.NullInt64 `db:"referrer"`
	GaduGadu sql.NullString `db:"gg"`
	ExtraInfo sql.NullString `db:"extrainfo"`
	Trial *uint16 `db:"trial"`
	ShowEmail sql.NullString `db:"showemail"`
	NbOfRefs sql.NullInt64 `db:"refer_count"`
	Suspended sql.NullString `db:"suspended"`
	ChangeIp sql.NullString `db:"change_ip"`
	ChangeUserId sql.NullInt64 `db:"change_user_id"`
	ChangeAt *time.Time `db:"change_date"`
	LoginAt *time.Time `db:"last_login"`
	CreatedAt *time.Time `db:"created"`
}
