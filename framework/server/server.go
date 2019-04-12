package server

type ServerType int

const (
	HTTP ServerType = iota
)

type Server interface {
	Start() error
	Stop() error
	IsRunning() bool
}
