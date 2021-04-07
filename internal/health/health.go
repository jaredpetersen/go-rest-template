package health

type Status struct {
	Up bool
}

func Check() *Status {
	return &Status{
		Up: true,
	}
}
