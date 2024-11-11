import {
  InvoiceStatus,
  ColumnViewType,
  ComparisonOperator,
} from '@shared/types/__generated__/graphql.types';

export type FilterType = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: any[];
  icon: JSX.Element;
  filterName: string;
  filterAccesor: ColumnViewType;
  filterOperators: ComparisonOperator[];
  filterType: 'text' | 'date' | 'number' | 'list';
  groupOptions?: { label: string; options: { id: string; label: string }[] }[];
};

import { File02 } from '@ui/media/icons/File02';
import { Invoice } from '@ui/media/icons/Invoice';
import { Calendar } from '@ui/media/icons/Calendar';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';

export const getFilterTypes = () => {
  const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
    [ColumnViewType.InvoicesInvoicePreview]: {
      filterType: 'text',
      filterName: 'Invoice number',
      filterAccesor: ColumnViewType.InvoicesInvoicePreview,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Invoice className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.InvoicesInvoiceNumber]: {
      filterType: 'text',
      filterName: 'Invoice number',
      filterAccesor: ColumnViewType.InvoicesInvoiceNumber,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <Invoice className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.InvoicesContract]: {
      filterType: 'text',
      filterName: 'Contract name',
      filterAccesor: ColumnViewType.InvoicesContract,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <File02 className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.InvoicesBillingCycle]: {
      filterType: 'list',
      filterName: 'Billing cycle',
      filterAccesor: ColumnViewType.InvoicesBillingCycle,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <ClockFastForward className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        { label: 'Monthly', id: 1 },
        { label: 'Quarterly', id: 3 },
        { label: 'Annually', id: 12 },
        { label: 'None', id: 0 },
      ],
    },

    [ColumnViewType.InvoicesIssueDate]: {
      filterType: 'date',
      filterName: 'Issue date',
      filterAccesor: ColumnViewType.InvoicesIssueDate,
      filterOperators: [ComparisonOperator.Gt, ComparisonOperator.Lt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },

    [ColumnViewType.InvoicesDueDate]: {
      filterType: 'date',
      filterName: 'Due date',
      filterAccesor: ColumnViewType.InvoicesDueDate,
      filterOperators: [ComparisonOperator.Gt, ComparisonOperator.Lt],
      icon: (
        <Calendar className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },
    [ColumnViewType.InvoicesAmount]: {
      filterType: 'number',
      filterName: 'Amount',
      filterAccesor: ColumnViewType.InvoicesAmount,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <CurrencyDollarCircle className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
    },

    [ColumnViewType.InvoicesInvoiceStatus]: {
      filterType: 'list',
      filterName: 'Invoice status',
      filterAccesor: ColumnViewType.InvoicesInvoiceStatus,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <ClockCheck className='group-hover:text-gray-700 text-gray-500 mb-0.5' />
      ),
      options: [
        { label: 'Out of contract', id: InvoiceStatus.OnHold },
        { label: 'Scheduled', id: InvoiceStatus.Scheduled },
        { label: 'Void', id: InvoiceStatus.Void },
        { label: 'Overdue', id: InvoiceStatus.Overdue },
        { label: 'Paid', id: InvoiceStatus.Paid },
      ],
    },
  };

  return filterTypes;
};
