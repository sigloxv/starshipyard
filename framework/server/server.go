package server

type ServerType int

const (
	HTTPServer ServerType = iota
)

type Server interface {
	Start() error
	Stop() error
	IsRunning() bool
}
