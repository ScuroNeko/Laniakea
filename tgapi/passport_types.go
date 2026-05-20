package tgapi

// PassportData contains information about Telegram Passport data shared with the bot.
// Since: Bot API 4.0
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// PassportFile represents a file uploaded to Telegram Passport.
// Since: Bot API 4.0
type PassportFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileDate     int64  `json:"file_date"`
}

// PassportElementType represents the type of a Telegram Passport element.
type PassportElementType string

const (
	PassportPersonalDetailsType       PassportElementType = "personal_details"
	PassportPassportType              PassportElementType = "passport"
	PassportDriverLicenseType         PassportElementType = "driver_license"
	PassportIdentityCardType          PassportElementType = "identity_card"
	PassportInternalPassportType      PassportElementType = "internal_passport"
	PassportAddressType               PassportElementType = "address"
	PassportUtilityBillType           PassportElementType = "utility_bill"
	PassportBankStatementType         PassportElementType = "bank_statement"
	PassportRentalAgreementType       PassportElementType = "rental_agreement"
	PassportPassportRegistrationType  PassportElementType = "passport_registration"
	PassportTemporaryRegistrationType PassportElementType = "temporary_registration"
	PassportPhoneNumberType           PassportElementType = "phone_number"
	PassportEmailType                 PassportElementType = "email"
)

// EncryptedPassportElement contains information about documents or other Telegram Passport elements.
// Since: Bot API 4.0
type EncryptedPassportElement struct {
	Type        PassportElementType `json:"type"`
	Data        string              `json:"data,omitempty"`
	PhoneNumber string              `json:"phone_number,omitempty"`
	Email       string              `json:"email,omitempty"`
	Files       []PassportFile      `json:"files,omitempty"`
	FrontSide   *PassportFile       `json:"front_side,omitempty"`
	ReverseSide *PassportFile       `json:"reverse_side,omitempty"`
	Selfie      *PassportFile       `json:"selfie,omitempty"`
	Translation *PassportFile       `json:"translation,omitempty"`
	Hash        string              `json:"hash,omitempty"`
}

// EncryptedCredentials contains data required for decrypting and authenticating EncryptedPassportElement.
// Since: Bot API 4.0
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// PassportElementError is a JSON-serializable passport element error object.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#passportelementerror
type PassportElementError struct {
	Source string              `json:"source"`
	Type   PassportElementType `json:"type"`

	FieldName string `json:"field_name,omitempty"`
	DataHash  string `json:"data_hash,omitempty"`

	FileHash   string   `json:"file_hash,omitempty"`
	FileHashes []string `json:"file_hashes,omitempty"`

	ElementHash string `json:"element_hash,omitempty"`

	Message string `json:"message"`
}
