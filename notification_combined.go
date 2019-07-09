package crawler

func NewCombinedNotification(notification ...Notification) *combinedNotification {
	return &combinedNotification{notification}
}

type combinedNotification struct {
	notifications []Notification
}

func (n *combinedNotification) Send(msg string) {
	for _, i := range n.notifications {
		i.Send(msg)
	}
}
