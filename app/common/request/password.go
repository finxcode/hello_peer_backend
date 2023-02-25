package request

type Password struct {
	Password string `json:"password" binding:"required"`
}

func (p *Password) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"password:required": "密码不能为空",
	}
}
