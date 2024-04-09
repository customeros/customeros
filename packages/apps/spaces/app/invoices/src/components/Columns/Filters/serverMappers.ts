import { mapIssueDateToAtom } from './IssueDate';
import { mapBillingCycleToAtom } from './BillingCycle';
import { mapInvoiceStatusToAtom } from './InvoiceStatus';
import { mapPaymentStatusToAtom } from './PaymentStatus';

const serverMappers = {
  ISSUE_DATE: mapIssueDateToAtom,
  BILLING_CYCLE: mapBillingCycleToAtom,
  PAYMENT_STATUS: mapPaymentStatusToAtom,
  INVOICE_STATUS: mapInvoiceStatusToAtom,
};

export const getServerToAtomMapper = (property: string) =>
  serverMappers[property as keyof typeof serverMappers];
