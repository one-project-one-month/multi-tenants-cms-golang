package types

type EmailRequest struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateName string                 `json:"templateName"`
	Data         map[string]interface{} `json:"data"`
}

type TenantCreationResponse struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
