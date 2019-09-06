package contracts

type Job interface {
	GetJobType() string
	Run() (interface{}, error)
	Stop()
}
