package ex

type Closable interface {
	Close() error
}

func CloseSilent(something Closable) {
	_ = something.Close()
}
