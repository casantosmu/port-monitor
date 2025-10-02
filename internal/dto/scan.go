package dto

type ScanRequest struct {
	IP   string `json:"ip" validate:"required,ip"`
	Port int    `json:"port" validate:"required,min=1,max=65535"`
}

type ScanResponse struct {
	Open bool `json:"open"`
}
