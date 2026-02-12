package defense

type AuthPayload struct {
	Signature   string `json:"signature"`
	UserName    string `json:"userName"`
	RandomKey   string `json:"randomKey"`
	PublicKey   string `json:"publicKey"`
	EncryptType string `json:"encryptType"`
	IpAddress   string `json:"ipAddress"`
	ClientType  string `json:"clientType"`
	UserType    string `json:"userType"`
}

type AuthRes struct {
	Token        string `json:"token"`
	Credential   string `json:"credential"`
	SecretKey    string `json:"secretKey"`
	SecretVector string `json:"secretVector"`
	Code         int    `json:"code"`
}

type EncDataPayload struct {
	UserName   string `json:"userName"`
	ClientType string `json:"clientType"`
}

type EncDataRes struct {
	Realm       string `json:"realm"`
	RandomKey   string `json:"randomKey"`
	EncryptType string `json:"encryptType"`
	Publickey   string `json:"publickey"`
}
