package controller

type Context struct{}

type Action func(context *Context) error

type Hooks struct {
	BeforeAction []Action
	AfterAction  []Action
}

type Controller struct {
	Hooks *Hooks
}

type ControllerInterface interface {
	BeforeAction()
}
