import { useMemo } from 'react';

export const useTablePlaceholder = (tableViewName?: string) => {
  return useMemo(() => {
    switch (tableViewName) {
      case 'Targets':
      case 'All orgs':
        return { multi: 'orgs', single: 'org' };
      case 'Customers':
        return { multi: 'customers', single: 'customer' };
      case 'All Contacts':
      case 'Contacts':
        return { multi: 'contacts', single: 'contacts' };
      case 'Leads':
        return { multi: 'leads', single: 'lead' };
      case 'Churn':
        return { multi: 'churned', single: 'churned' };
      case 'Contracts':
        return { multi: 'contracts', single: 'contract' };
      case 'Opportunities':
        return { multi: 'opps', single: 'opp' };
      case 'Sequences':
        return { multi: 'sequences', single: 'sequence' };
      case 'Past':
      case 'Upcoming':
        return { multi: 'invoices', single: 'invoice' };
      default:
        return { multi: 'orgs', single: 'org' };
    }
  }, [tableViewName]);
};
