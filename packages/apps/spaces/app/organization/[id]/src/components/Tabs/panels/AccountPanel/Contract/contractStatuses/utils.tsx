import React from 'react';

import { Clock } from '@ui/media/icons/Clock';
import { ContractStatus } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { XSquare } from '@ui/media/icons/XSquare';
import { DotLive } from '@ui/media/icons/DotLive';
import { DateTimeUtils } from '@spaces/utils/date';
import { PauseCircle } from '@ui/media/icons/PauseCircle';

export const contractOptionIcon: Record<ContractStatus, JSX.Element | null> = {
  [ContractStatus.Draft]: <Edit03 className='text-gray-500' />,
  [ContractStatus.Ended]: <XSquare className='text-gray-500' />,
  [ContractStatus.Live]: <DotLive />,
  [ContractStatus.OutOfContract]: <PauseCircle className='text-warning-500' />,
  [ContractStatus.Scheduled]: <Clock className='text-primary-600 size-3' />,
  [ContractStatus.Undefined]: null,
};

export const confirmationModalDataByStatus = {
  [ContractStatus.Live]: {
    title: 'Make this a draft contract',
    description: (orgName: string, serviceStartDate: string) => (
      <>
        To make
        <span className='font-medium mx-1'>{orgName}</span>
        contract a draft, we will clear the service start date (currently{' '}
        {serviceStartDate}).
      </>
    ),
    submit: 'Make live',
    colorScheme: 'primary',
  },
  [ContractStatus.Draft]: {
    title: 'Make this contract live',
    description: (orgName: string, serviceStartDate: string) => (
      <>
        Congrats! Let’s make
        <span className='font-medium mx-1'>{orgName}</span>
        contract live on: contract a draft, we will clear the service start date
        (currently{' '}
        {DateTimeUtils.format(
          serviceStartDate,
          DateTimeUtils.dateWithAbreviatedMonth,
        )}
        ).
      </>
    ),
    submit: 'Make draft',
    colorScheme: 'gray',
  },
  [ContractStatus.Ended]: {
    title: 'End this contract?',
    description: (orgName: string) => (
      <>
        Ending
        <span className='font-medium mx-1'>{orgName}</span>
        contract will close the renewal and set the ARR to zero.
        <br />
        Let’s end it on:
      </>
    ),
    submit: 'Ended contract',
    colorScheme: 'error',
  },
};
