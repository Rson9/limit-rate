package limiter

type Limter interface {
	Check() bool
}
