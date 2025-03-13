package enum

type WebTrackerEvent string

const (
	WebTrackerPageExit WebTrackerEvent = "page_exit"
	WebTrackerPageView WebTrackerEvent = "page_view"
	WebTrackerClick    WebTrackerEvent = "click"
)

func (w WebTrackerEvent) String() string {
	return string(w)
}
