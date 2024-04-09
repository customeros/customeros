import { useEffect } from 'react';

import { match } from 'ts-pattern';
import { useRecoilTransaction_UNSTABLE } from 'recoil';

import { Filter } from '@graphql/types';

import { getServerToAtomMapper } from './serverMappers';
import { IssueDateFilterAtom, IssueDateFilterState } from './IssueDate';
import {
  BillingCycleFilterAtom,
  BillingCycleFilterState,
} from './BillingCycle';
import {
  PaymentStatusFilterAtom,
  PaymentStatusFilterState,
} from './PaymentStatus';
import {
  InvoiceStatusFilterAtom,
  InvoiceStatusFilterState,
} from './InvoiceStatus';

export const parseRawFilters = (raw = '') => {
  if (!raw) return [];
  const filterData = JSON.parse(raw) as { filter: Filter };

  if (!filterData?.filter?.AND) return [];

  return filterData?.filter?.AND?.map((data) => {
    const property = data?.filter?.property ?? '';
    const mapToAtom = getServerToAtomMapper(property);

    if (mapToAtom) {
      return [property, mapToAtom(data)];
    }
  }).filter(Boolean) as [string, BillingCycleFilterState][];
};

export const useFilterSetter = (rawFilters?: string | null) => {
  const parsedFilters = parseRawFilters(rawFilters ?? '');

  const setFilters = useRecoilTransaction_UNSTABLE(
    ({ set }) =>
      (id: string, value: unknown) => {
        match(id)
          .with('BILLING_CYCLE', () => {
            set<BillingCycleFilterState>(
              BillingCycleFilterAtom,
              value as BillingCycleFilterState,
            );
          })
          .with('ISSUE_DATE', () => {
            set<IssueDateFilterState>(
              IssueDateFilterAtom,
              value as IssueDateFilterState,
            );
          })
          .with('PAYMENT_STATUS', () => {
            set<PaymentStatusFilterState>(
              PaymentStatusFilterAtom,
              value as PaymentStatusFilterState,
            );
          })
          .with('INVOICE_STATUS', () => {
            set<InvoiceStatusFilterState>(
              InvoiceStatusFilterAtom,
              value as InvoiceStatusFilterState,
            );
          })
          .otherwise(() => {});
      },
  );

  useEffect(() => {
    parsedFilters.forEach(([id, value]) => {
      setFilters(id, value);
    });
  }, [rawFilters]);
};
