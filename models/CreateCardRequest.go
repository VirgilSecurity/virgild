package models

// type CreateCardRequest struct {
// 	CardsResponse `json:"-"`
// 	Identity      string            `json:"identity"`
// 	IdentityType  string            `json:"identity_type"`
// 	PublicKey     []byte            `json:"public_key"` //DER encoded public key
// 	Scope         enums.VirgilEnum  `json:"scope"`
// 	Data          map[string]string `json:"data,omitempty"`
// 	DeviceInfo    DeviceInfo        `json:"info"`
// }

type DeviceInfo struct {
	Device     string `json:"device"`
	DeviceName string `json:"device_name"`
}
