package form

type FormMethod int
type FormMethods map[FormMethod]bool

const (
	POST FormMethod = iota
	GET
	PUT
	PATCH
	DELETE
)

// TODO: What about initializing the form based on the Method?
// As in, PostMethod.Form()
// TODO: Decide if we should only have the "real" POST/GET then
func (self FormMethod) String() string {
	switch self {
	case GET:
		return `GET`
	case PUT:
		return `PUT`
	case PATCH:
		return `PATCH`
	case DELETE:
		return `DELETE`
	default: // POST
		return `POST`
	}
}

func AllowedFormMethods() FormMethods {
	return FormMethods{
		GET:    true,
		PUT:    false,
		POST:   true,
		PATCH:  false,
		DELETE: false,
	}
}

func MarshalMethod(methodName string) FormMethod {
	for method, allowed := range AllowedFormMethods() {
		if allowed && method.String() == methodName {
			return method
		}
	}
	return POST // Default to POST
}
