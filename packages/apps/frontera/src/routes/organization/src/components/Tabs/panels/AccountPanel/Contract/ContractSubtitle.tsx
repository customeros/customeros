import { useMemo } from 'react';

import { toZonedTime } from 'date-fns-tz';
import { observer } from 'mobx-react-lite';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Contract, ContractStatus } from '@graphql/types';

export function getCommittedPeriodLabel(months: string | number) {
  if (`${months}` === '1') {
    return 'Monthly';
  }
  if (`${months}` === '3') {
    return 'Quarterly';
  }

  if (`${months}` === '12') {
    return 'Annual';
  }

  return `${months}-month`;
}

export const ContractSubtitle = observer(({ id }: { id: string }) => {
  const { contracts } = useStore();
  const data = contracts.value.get(id)?.value as Contract;
  const serviceStarted = data?.serviceStarted
    ? toZonedTime(data?.serviceStarted, 'UTC').toUTCString()
    : null;

  const contractEnded = data?.contractEnded
    ? toZonedTime(data?.contractEnded, 'UTC').toUTCString()
    : null;

  const renewalCalculatedDate = useMemo(() => {
    if (!serviceStarted) return null;
    const parsed = data?.committedPeriodInMonths
      ? parseFloat(data?.committedPeriodInMonths)
      : 1;

    return DateTimeUtils.addMonth(serviceStarted, parsed).toString();
  }, [data?.serviceStarted, data?.committedPeriodInMonths]);
  const renewalDate = renewalCalculatedDate
    ? DateTimeUtils.format(
        renewalCalculatedDate,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  const hasStartedService =
    serviceStarted && !DateTimeUtils.isFuture(serviceStarted);

  const serviceStartDate =
    serviceStarted && DateTimeUtils.isFuture(serviceStarted)
      ? DateTimeUtils.format(
          serviceStarted,
          DateTimeUtils.dateWithAbreviatedMonth,
        )
      : null;

  const endDate = contractEnded
    ? DateTimeUtils.format(contractEnded, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  const renewalPeriod = getCommittedPeriodLabel(data?.committedPeriodInMonths);

  const isJustCreated =
    DateTimeUtils.differenceInMins(
      data.metadata.lastUpdated,
      data.metadata.created,
    ) === 0;

  if (isJustCreated && !serviceStartDate) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        Monthly contract{' '}
        {data?.autoRenew ? 'auto-renewing' : 'not auto-renewing'}
      </p>
    );
  }
  if (
    !hasStartedService &&
    !serviceStartDate &&
    data?.committedPeriodInMonths
  ) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod ? `${renewalPeriod} contract` : 'Contract'} starting ...
        <Button
          className='underline ml-1 p-0 font-normal text-sm text-gray-500 hover:text-gray-500 focus:text-gray-500'
          variant='link'
          size='xs'
        >
          Edit contract
        </Button>
      </p>
    );
  }
  if (
    !hasStartedService &&
    serviceStartDate &&
    renewalPeriod &&
    data?.contractStatus === ContractStatus.Draft
  ) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod} contract{' '}
        {data?.autoRenew ? 'auto-renewing' : 'not auto-renewing'}
      </p>
    );
  }

  if (!hasStartedService && serviceStartDate && data?.committedPeriodInMonths) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod ? `${renewalPeriod} contract` : 'Contract'} set to go
        live on {serviceStartDate}
      </p>
    );
  }
  if (hasStartedService && endDate) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod ? `${renewalPeriod} contract` : 'Contract'}{' '}
        {contractEnded && DateTimeUtils.isFuture(contractEnded)
          ? 'ending'
          : 'ended on'}{' '}
        {endDate}
      </p>
    );
  }
  if (hasStartedService && renewalDate && data?.autoRenew) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod ? `${renewalPeriod} contract` : 'Contract'} auto-renewing{' '}
        {renewalDate}
      </p>
    );
  }

  if (hasStartedService && !data?.autoRenew) {
    return (
      <p className='font-normal shadow-none text-sm  text-gray-500 focus:text-gray-500 hover:text-gray-500 hover:no-underline focus:no-underline'>
        {renewalPeriod ? `${renewalPeriod} contract` : 'Contract'} until{' '}
        {renewalDate}, not auto-renewing
      </p>
    );
  }

  return null;
});
