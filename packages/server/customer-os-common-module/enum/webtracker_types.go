package enum

type WebSessionTimeouts int

const (
	WebSessionTimeoutPageExit WebSessionTimeouts = 5  // mins -- page exit without a following page view
	WebSessionTimeoutPageView WebSessionTimeouts = 30 // mins -- page view without a page exit
)

type WebTrackerEvent string

const (
	WebTrackerClick    WebTrackerEvent = "click"
	WebTrackerPageExit WebTrackerEvent = "page_exit"
	WebTrackerPageView WebTrackerEvent = "page_view"
)

func (w WebTrackerEvent) String() string {
	return string(w)
}
