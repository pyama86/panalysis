package panalysis

const MaxBuffer = 4096

type Parser interface {
	Parse() (interface{}, error)
}
