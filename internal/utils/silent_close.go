package utils

type Closable interface {
	Close() error
}

func SilentClose(obj Closable) {
	_ = obj.Close()
}
