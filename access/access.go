package access

type Getter interface {
	Get(fileName string, offset int, length int) ([]byte, error)
}
