package utils

type Scanneable interface {
	Scan(dest ...any) error
}
