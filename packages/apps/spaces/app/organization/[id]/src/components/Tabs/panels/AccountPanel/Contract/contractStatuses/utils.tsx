import React from 'react';

import { Text } from '@ui/typography/Text';
import { ContractStatus } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { XSquare } from '@ui/media/icons/XSquare';
import { DotLive } from '@ui/media/icons/DotLive';
import { DateTimeUtils } from '@spaces/utils/date';

export const contractOptionIcon: Record<ContractStatus, JSX.Element | null> = {
  [ContractStatus.Draft]: <Edit03 color='gray.500' boxSize='inherit' />,
  [ContractStatus.Ended]: <XSquare color='gray.500' boxSize='inherit' />,
  [ContractStatus.Live]: <DotLive color='inherit' boxSize='inherit' />,
  [ContractStatus.Undefined]: null,
};

export const confirmationModalDataByStatus = {
  [ContractStatus.Live]: {
    title: 'Make this a draft contract',
    description: (orgName: string, serviceStartDate: string) => (
      <>
        To make
        <Text as='span' fontWeight='medium' mx={1}>
          {orgName}
        </Text>
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
        <Text as='span' fontWeight='medium' mx={1}>
          {orgName}
        </Text>
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
        <Text as='span' fontWeight='medium' mx={1}>
          {orgName}
        </Text>
        contract will close the renewal and set the ARR to zero.
        <br />
        Let’s end it on:
      </>
    ),
    submit: 'Ended contract',
    colorScheme: 'error',
  },
};
