package cmdbus

type Validator interface {
	Validate(s any) error
}
