package auth

type AuthClient struct {
	AuthServiceURL string `json:"auth_service_url"`
	JWTSecret string `json:"jwt_secret"`
}


type AuthWebhookEvent struct {
    Table  string `json:"table"`
    Type   string `json:"type"`
    Record struct {
        ID        string `json:"id"`
        Email     string `json:"email"`
        CreatedAt string `json:"created_at"`
    } `json:"record"`
}