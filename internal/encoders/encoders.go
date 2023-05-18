package encoders

type Encoder string

const (
	None   Encoder = ""
	Hex    Encoder = "hex"
	Base64 Encoder = "base64"
)
