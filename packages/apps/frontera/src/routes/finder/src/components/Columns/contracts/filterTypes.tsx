import {
  ColumnViewType,
  ContractStatus,
  ComparisonOperator,
  OpportunityRenewalLikelihood,
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

import { type RootStore } from '@store/root';

import { Key01 } from '@ui/media/icons/Key01';
import { File02 } from '@ui/media/icons/File02';
import { Calendar } from '@ui/media/icons/Calendar';
import { Activity } from '@ui/media/icons/Activity';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { Calculator } from '@ui/media/icons/Calculator';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { CurrencyDollarCircle } from '@ui/media/icons/CurrencyDollarCircle';

export const getFilterTypes = (store?: RootStore) => {
  const filterTypes: Partial<Record<ColumnViewType, FilterType>> = {
    [ColumnViewType.ContractsName]: {
      filterType: 'text',
      filterName: 'Contract name',
      filterAccesor: ColumnViewType.ContractsName,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <File02 className='group-hover:text-gray-700 text-gray-500' />,
    },
    [ColumnViewType.ContractsEnded]: {
      filterType: 'date',
      filterName: 'Ended',
      filterAccesor: ColumnViewType.ContractsEnded,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: <Calendar className='group-hover:text-gray-700 text-gray-500' />,
    },
    [ColumnViewType.ContractsPeriod]: {
      filterType: 'list',
      filterName: 'Period',
      filterAccesor: ColumnViewType.ContractsPeriod,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <ClockFastForward className='group-hover:text-gray-700 text-gray-500' />
      ),
      options: [
        { id: '1', label: 'Monthly' },
        { id: '3', label: 'Quarterly' },
        { id: '12', label: 'Annually' },
      ],
    },

    [ColumnViewType.ContractsCurrency]: {
      filterType: 'list',
      filterName: 'Currency',
      filterAccesor: ColumnViewType.ContractsCurrency,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: (
        <CurrencyDollarCircle className='group-hover:text-gray-700 text-gray-500' />
      ),
      options: [
        { id: 'USD', label: 'USD' },
        { id: 'EUR', label: 'EUR' },
        { id: 'GBP', label: 'GBP' },
      ],
    },

    [ColumnViewType.ContractsStatus]: {
      filterType: 'list',
      filterName: 'Status',
      filterAccesor: ColumnViewType.ContractsStatus,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <ClockCheck className='group-hover:text-gray-700 text-gray-500' />,
      options: [
        { label: 'Draft', id: ContractStatus.Draft },
        { label: 'Ended', id: ContractStatus.Ended },
        { label: 'Live', id: ContractStatus.Live },
        { label: 'Out of Contract', id: ContractStatus.OutOfContract },
        { label: 'Scheduled', id: ContractStatus.Scheduled },
      ],
    },
    [ColumnViewType.ContractsRenewal]: {
      filterType: 'list',
      filterName: 'Renewal',
      filterAccesor: ColumnViewType.ContractsRenewal,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Calendar className='grou-hover:text-gray-700 text-gray-500' />,
      options: [
        { label: 'Auto-renews', id: 'Auto-renews' },
        { label: 'Not auto-renewing', id: 'Not auto-renewing' },
      ],
    },
    [ColumnViewType.ContractsLtv]: {
      filterType: 'number',
      filterName: 'LTV',
      filterAccesor: ColumnViewType.ContractsLtv,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: (
        <CurrencyDollarCircle className='group-hover:text-gray-700 text-gray-500' />
      ),
    },
    [ColumnViewType.ContractsRenewalDate]: {
      filterType: 'date',
      filterName: 'Renewal date',
      filterAccesor: ColumnViewType.ContractsRenewalDate,
      filterOperators: [ComparisonOperator.Lt, ComparisonOperator.Gt],
      icon: <Calendar className='group-hover:text-gray-700 text-gray-500' />,
    },
    [ColumnViewType.ContractsForecastArr]: {
      filterType: 'number',
      filterName: 'ARR forecast',
      filterAccesor: ColumnViewType.ContractsForecastArr,
      filterOperators: [
        ComparisonOperator.Gt,
        ComparisonOperator.Lt,
        ComparisonOperator.Eq,
        ComparisonOperator.NotEqual,
      ],
      icon: <Calculator className='group-hover:text-gray-700 text-gray-500' />,
    },
    [ColumnViewType.ContractsHealth]: {
      filterType: 'list',
      filterName: 'Health',
      filterAccesor: ColumnViewType.ContractsHealth,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Activity className='group-hover:text-gray-700 text-gray-500' />,
      options: [
        { id: OpportunityRenewalLikelihood.HighRenewal, label: 'High' },
        { id: OpportunityRenewalLikelihood.MediumRenewal, label: 'Medium' },
        { id: OpportunityRenewalLikelihood.LowRenewal, label: 'Low' },
        { id: OpportunityRenewalLikelihood.ZeroRenewal, label: 'Zero' },
      ],
    },
    [ColumnViewType.ContractsOwner]: {
      filterType: 'list',
      filterName: 'Owner',
      filterAccesor: ColumnViewType.ContractsOwner,
      filterOperators: [
        ComparisonOperator.Contains,
        ComparisonOperator.NotContains,
        ComparisonOperator.IsEmpty,
        ComparisonOperator.IsNotEmpty,
      ],
      icon: <Key01 className='group-hover:text-gray-700 text-gray-500' />,
      options: store?.users.toArray().map((user) => ({
        id: user?.id,
        label: user?.name,
        avatar: user?.value?.profilePhotoUrl,
      })),
    },
  };

  return filterTypes;
};
