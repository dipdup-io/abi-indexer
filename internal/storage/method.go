package storage

// IMethod -
type IMethod interface {
	Table[*Method]
}

// Method -
type Method struct {
	// nolint
	tableName struct{} `pg:"methods"`

	ID uint64

	Name        string
	Type        int
	Mutability  string
	IsConst     bool `pg:"default:false"`
	IsPayable   bool `pg:"default:false"`
	Signature   string
	SignatureID []byte
	MetadataID  uint64

	Metadata *Metadata `pg:",rel:has-one"`
}

// TableName -
func (Method) TableName() string {
	return "methods"
}
