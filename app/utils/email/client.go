package email

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/email/my_email"
)

func NewEmail() *my_email.Email {
	return &my_email.Email{
		Host:     variable.ConfigYml.GetString("Email.Host"),
		Port:     variable.ConfigYml.GetInt("Email.Port"),
		IsSSL:    variable.ConfigYml.GetBool("Email.IsSSL"),
		UserName: variable.ConfigYml.GetString("Email.UserName"),
		Password: variable.ConfigYml.GetString("Email.Password"),
		From:     variable.ConfigYml.GetString("Email.From"),
	}
}
