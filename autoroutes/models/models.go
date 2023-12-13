package models

type GenerateFunction func(args ...interface{}) interface{}

type GenerateFields map[string]GenerateFunction

type ValidationFunction func(value interface{}) bool

type ValidationFields map[string]ValidationFunction

type Route struct {
	Methods    []string
	DataEntity interface{}
	Config     Config
}

type Config struct {
	SelfResource         bool     // GET PUT POST DELETE
	SelfResourceField    string   // GET PUT POST DELETE
	SelectFields         []string // GET PUT POST
	OmitFields           []string // GET PUT POST
	GetLimit             int
	GetSearchFields      []string
	GetMessage           string
	PostMessage          string
	PostDuplicateMessage string
	PostSkipOmit         bool // genetated fields cannot be omitted
	PostNewFields        []string
	PostGenerateValues   GenerateFields
	PostValidationFields ValidationFields
	PutMessage           string
	PutUpdatableFields   []string
	DeleteMessage        string
}

type Routes []Route

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "Delete"
)
