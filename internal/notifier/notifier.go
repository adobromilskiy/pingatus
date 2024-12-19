package notifier

type Notifier interface {
	Send(msg string)
}
