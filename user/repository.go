package user

import 	(
	"github.com/piotrkowalczuk/netwars-backend/database"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository struct {
	repository database.Repository
}

func NewUserRepository() (repo *UserRepository) {
	repo = &UserRepository{database.Repository{"users"}}

	return
}

func (self *UserRepository) super() *database.Repository {
	return &self.repository
}

func (self *UserRepository) FindOne(id int) (user *User){
	user = new(User)

	err := self.super().FindOne("user_id", id).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.PasswordType,
		&user.PasswordSalt,
		&user.BadLogins,
		&user.Email,
		&user.NTCNick,
		&user.NickHistory,
		&user.Status,
		&user.EmailUsed,
		&user.ReferrerId,
		&user.GaduGadu,
		&user.ExtraInfo,
		&user.Trial,
		&user.ShowEmail,
		&user.NbOfRefs,
		&user.Suspended,
		&user.ChangeIp,
		&user.ChangeUserId,
		&user.ChangeAt,
		&user.LoginAt,
		&user.CreatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Println(user)
	}

	return
}
