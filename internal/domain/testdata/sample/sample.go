package sample

type Sample struct {
	Default  string
	Named    string   `redact:"named"`
	Allowed  string   `redact:",allow"`
	Redactor string   `redact:",,pan64"`
	All      string   `redact:"all,allow,star"`
	Skipped  string   `redact:"-"`
	JSON     string   `json:"json_name,omitempty"`
	Slice    []string `redact:",,len"`
}
