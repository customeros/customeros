import { useMemo } from 'react';

export const useTablePlaceholder = (tableViewName?: string) => {
  return useMemo(() => {
    switch (tableViewName) {
      case 'Targets':
      case 'All orgs':
        return { multi: 'organizations', single: 'organization' };
      case 'Customers':
        return { multi: 'customers', single: 'customer' };
      case 'All Contacts':
        return { multi: 'contacts', single: 'contacts' };
      case 'Leads':
        return { multi: 'leads', single: 'lead' };
      case 'Churn':
        return { multi: 'churned', single: 'churned' };
      case 'Past':
      case 'Upcoming':
        return { multi: 'invoices', single: 'invoice' };
      default:
        return { multi: 'organizations', single: 'organization' };
    }
  }, [tableViewName]);
};
