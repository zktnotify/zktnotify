package typed

type NotifierType uint32

const (
	DINGTALK NotifierType = iota
	SERVERCHAN
)

type Receiver struct {
	All bool
	ID  []string
}

type Notifier interface {
	Notify(url string, msg string, receiver ...Receiver) error
}

func Valid(ntype NotifierType) bool {
	switch ntype {
	case DINGTALK, SERVERCHAN:
		return true
	}
	return false
}
