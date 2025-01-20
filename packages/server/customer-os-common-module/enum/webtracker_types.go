package enum

type WebTrackerEvent string

const (
	WebTrackerClick    WebTrackerEvent = "click"
	WebTrackerPageExit WebTrackerEvent = "page_exit"
	WebTrackerPageView WebTrackerEvent = "page_view"
)

func (w WebTrackerEvent) String() string {
	return string(w)
}
