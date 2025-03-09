package port

type ResponseService interface {
	Run() error
	Close() error
}
