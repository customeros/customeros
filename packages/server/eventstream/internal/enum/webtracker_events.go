package enum

type WebTrackerEvent string

const (
	WebTrackerPageExit WebTrackerEvent = "page_exit"
	WebTrackerPageView WebTrackerEvent = "page_view"
	WebTrackerClick    WebTrackerEvent = "click"
	WebTrackerIdentify WebTrackerEvent = "identify"
)

func (w WebTrackerEvent) String() string {
	return string(w)
}

func IsValidWebTrackerEvent(event string) bool {
	switch event {
	case WebTrackerPageExit.String():
		return true
	case WebTrackerPageView.String():
		return true
	case WebTrackerClick.String():
		return true
	case WebTrackerIdentify.String():
		return true
	}
	return false
}
