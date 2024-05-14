import { Clock } from '@ui/media/icons/Clock';
import { ContractStatus } from '@graphql/types';
import { Edit03 } from '@ui/media/icons/Edit03';
import { XSquare } from '@ui/media/icons/XSquare';
import { DotLive } from '@ui/media/icons/DotLive';
import { PauseCircle } from '@ui/media/icons/PauseCircle';

export const contractOptionIcon: Record<ContractStatus, JSX.Element | null> = {
  [ContractStatus.Draft]: <Edit03 className='text-gray-500 size-3' />,
  [ContractStatus.Ended]: <XSquare className='text-gray-500 size-3' />,
  [ContractStatus.Live]: <DotLive className='size-3' />,
  [ContractStatus.OutOfContract]: (
    <PauseCircle className='text-warning-500 size-3' />
  ),
  [ContractStatus.Scheduled]: <Clock className='text-primary-600 size-3' />,
  [ContractStatus.Undefined]: null,
};
