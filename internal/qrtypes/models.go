package qrtypes

type QRTypeData struct {
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type WiFiData struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Security string `json:"security"` // WPA/WEP
}

type VCardData struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Website  string `json:"website"`
}
