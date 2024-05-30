package enum

type InvoiceStatus string

const (
	InvoiceStatusNone        InvoiceStatus = ""
	InvoiceStatusInitialized InvoiceStatus = "INITIALIZED"
	InvoiceStatusDue         InvoiceStatus = "DUE"
	InvoiceStatusPaid        InvoiceStatus = "PAID"
	InvoiceStatusVoid        InvoiceStatus = "VOID"
	InvoiceStatusScheduled   InvoiceStatus = "SCHEDULED"
	InvoiceStatusOverdue     InvoiceStatus = "OVERDUE"
	InvoiceStatusOnHold      InvoiceStatus = "ON_HOLD"
	InvoiceStatusEmpty       InvoiceStatus = "EMPTY"
)

var AllInvoiceStatuses = []InvoiceStatus{
	InvoiceStatusNone,
	InvoiceStatusInitialized,
	InvoiceStatusDue,
	InvoiceStatusPaid,
	InvoiceStatusVoid,
	InvoiceStatusScheduled,
	InvoiceStatusOverdue,
	InvoiceStatusOnHold,
	InvoiceStatusEmpty,
}

func DecodeInvoiceStatus(s string) InvoiceStatus {
	if IsValidInvoiceStatus(s) {
		return InvoiceStatus(s)
	}
	return InvoiceStatusNone
}

func IsValidInvoiceStatus(s string) bool {
	for _, ms := range AllInvoiceStatuses {
		if ms == InvoiceStatus(s) {
			return true
		}
	}
	return false
}

func (e InvoiceStatus) String() string {
	return string(e)
}
