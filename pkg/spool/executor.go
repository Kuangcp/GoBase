package spool

var ()

type (
	SimExecutor struct {
		max int
	}
)

func (s *SimExecutor) Submit(labels map[string]string, runnable Runnable) Future {
	//TODO implement me
	panic("implement me")
}

func (s *SimExecutor) SubmitR(labels map[string]string, runnable Runnable, result any) Future {
	//TODO implement me
	panic("implement me")
}

func (s *SimExecutor) SubmitC(labels map[string]string, runnable Callable) Future {
	//TODO implement me
	panic("implement me")
}
