package dtos

type GlobalHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DTOCreateUserReq struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Cep          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
}
