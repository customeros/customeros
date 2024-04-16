import React from 'react';

import { Clock } from '@ui/media/icons/Clock';
import { ContractStatus } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { XSquare } from '@ui/media/icons/XSquare';
import { DotLive } from '@ui/media/icons/DotLive';
import { PauseCircle } from '@ui/media/icons/PauseCircle';

export const contractOptionIcon: Record<ContractStatus, JSX.Element | null> = {
  [ContractStatus.Draft]: <Edit03 color='gray.500' boxSize='3' />,
  [ContractStatus.Ended]: <XSquare color='gray.500' boxSize='3' />,
  [ContractStatus.Live]: <DotLive color='inherit' boxSize='3' />,
  [ContractStatus.OutOfContract]: (
    <PauseCircle color='warning.500' boxSize='inherit' />
  ),
  [ContractStatus.Scheduled]: <Clock color='primary.600' boxSize='3' />,
  [ContractStatus.Undefined]: null,
};
