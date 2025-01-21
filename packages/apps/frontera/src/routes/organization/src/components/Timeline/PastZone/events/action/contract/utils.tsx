import { ReactElement } from 'react';

import { FileX02 } from '@ui/media/icons/FileX02';
import { PauseCircle } from '@ui/media/icons/PauseCircle';
import { FileHeart02 } from '@ui/media/icons/FileHeart02';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';

export const iconsByStatus: Record<
  string,
  Record<string, string | ReactElement>
> = {
  live: {
    icon: <FileHeart02 className='text-primary-600' />,
    text: 'is now',
  },
  renewed: {
    icon: <FileCheck02 className='text-success-500' />,
    text: '',
  },
  ended: {
    icon: <FileX02 className='text-gray-500' />,
    text: 'has',
  },
  out_of_contract: {
    icon: <PauseCircle className='text-warning-500' />,
    text: 'has',
  },
};
