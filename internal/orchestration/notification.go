package orchestration

const (
	SYSTEM_NOTIFICATION_EVENT = EventType("system_notification_event")
)

var notifChan chan *Notification = make(chan *Notification, 5)

type Notification struct {
	Content string
	Type    string
}

func NewNotification(notifType, content string) *Notification {
	return &Notification{
		Content: content,
		Type:    notifType,
	}
}

func SendErrorNotification(content string) {
	go func(ch chan<- *Notification) {
		ch <- &Notification{
			Type:    "error",
			Content: content,
		}
	}(notifChan)
}

func SendInfoNotification(content string) {
	go func(ch chan<- *Notification) {
		ch <- &Notification{
			Type:    "info",
			Content: content,
		}
	}(notifChan)
}

func SendWarningNotification(content string) {
	go func(ch chan<- *Notification) {
		ch <- &Notification{
			Type:    "warning",
			Content: content,
		}
	}(notifChan)
}

func SendSuccessNotification(content string) {
	go func(ch chan<- *Notification) {
		ch <- &Notification{
			Type:    "success",
			Content: content,
		}
	}(notifChan)
}

type NotificationMonitor struct {
	EventType EventType
}

func start(doneChan <-chan struct{}, eventChan chan<- *EventChannelData, notificationChan <-chan *Notification) {
	for {
		select {
		case <-doneChan:
			return
		case notif := <-notificationChan:
			eventChan <- ToEventData(SYSTEM_NOTIFICATION_EVENT, notif.Type, map[string]string{
				"Type":    notif.Type,
				"Content": notif.Content,
			})
		}
	}
}

func startNotificationMonitor(sm *SubscriptionManager, stopChan <-chan struct{}) {
	eventChan := make(chan *EventChannelData)
	go start(stopChan, eventChan, notifChan)

	for {
		select {
		case <-stopChan:
			return
		case event := <-eventChan:
			sm.mu.Lock()
			for _, listenerChan := range sm.listeners[SYSTEM_NOTIFICATION_EVENT] {
				go func(ch chan *EventChannelData) {
					ch <- event
				}(listenerChan)
			}
			sm.mu.Unlock()
		}
	}
}
