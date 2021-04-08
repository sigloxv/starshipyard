package helpers

const (
	EmptyString = ""
	Blank       = ""
)

func IsBlank(s string) bool {
	return (s == Blank)
}

func IsEmpty(s string) bool {
	return (s == EmptyString)
}
