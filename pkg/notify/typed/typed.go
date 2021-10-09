package typed

type NotifierType uint32

const (
	DINGTALK NotifierType = iota
	SERVERCHAN
	WXPUSHER
)

// TemplateID message template type id
type TemplateID uint64

const (
	Invalid TemplateID = iota
	ToWork
	Worked
	Midway
	Lated
	Remind    // notify to take a card
	DelayWork // delay in work
	MonthDaily
)

type WorkType uint64

const (
	Working WorkType = iota
	OffWork
)

func (wt WorkType) String() string {
	if wt == Working {
		return "上班"
	}
	return "下班"
}

func (s TemplateID) String() string {
	status := "未知"
	switch s {
	case Invalid:
		status = "未知"
	case ToWork:
		status = "正常"
	case Worked:
		status = "正常"
	case Midway:
		status = "正常"
	case Lated:
		status = "迟到"
	case Remind:
		status = "提醒"
	case DelayWork:
		status = "正常"
	case MonthDaily:
		status = "提醒"
	}
	return status
}

type Message struct {
	UID        uint64
	Name       string
	Date       string
	Time       string
	Type       WorkType // 1 on work, 2 off work
	Token      string
	Status     TemplateID
	Account    string
	CancelURL  string
	NotifyType NotifierType
}

type Receiver struct {
	All bool
	ID  []string
}

type Notifier interface {
	Notify(token string, msg string, receiver ...Receiver) error
	SetCancelURL(url string)
	SetAppToken(token string)
	Template(msg Message) string
}

func Valid(ntype NotifierType) bool {
	switch ntype {
	case DINGTALK, SERVERCHAN, WXPUSHER:
		return true
	}
	return false
}
