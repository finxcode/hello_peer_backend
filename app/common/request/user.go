package request

type Register struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (register Register) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"Name.required":     "user name field cannot be empty",
		"Mobile.required":   "phone no field cannot be empty",
		"Password.required": "password cannot be empty",
		"mobile.mobile":     "format of phone number incorrect",
	}
}

type Login struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (login Login) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"mobile.required":   "phone no field cannot be empty",
		"mobile.mobile":     "format of phone number incorrect",
		"password.required": "password cannot be empty",
	}
}
