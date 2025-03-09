package port

type RequestService interface {
	Run(iterations int) error
	Close() error
}
