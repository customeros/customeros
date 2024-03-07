import React, { ReactElement } from 'react';

import { FileX02 } from '@ui/media/icons/FileX02';
import { FileHeart02 } from '@ui/media/icons/FileHeart02';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';

export const iconsByStatus: Record<
  string,
  Record<string, string | ReactElement>
> = {
  live: {
    icon: <FileHeart02 />,
    colorScheme: 'primary',
    text: 'is now',
  },
  renewed: {
    icon: <FileCheck02 />,
    colorScheme: 'success',
    text: '',
  },
  ended: {
    icon: <FileX02 />,
    colorScheme: 'gray',
    text: 'has',
  },
};
