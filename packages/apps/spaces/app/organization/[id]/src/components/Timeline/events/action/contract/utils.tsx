import React, { ReactElement } from 'react';

import { FolderCheck } from '@ui/media/icons/FolderCheck';
import { FolderClosed } from '@ui/media/icons/FolderClosed';

export const iconsByStatus: Record<
  string,
  Record<string, string | ReactElement>
> = {
  live: {
    icon: <FolderCheck />,
    colorScheme: 'primary',
    text: 'is now',
  },
  renewed: {
    icon: <FolderCheck />,
    colorScheme: 'success',
    text: '',
  },
  ended: {
    icon: <FolderClosed />,
    colorScheme: 'gray',
    text: 'has',
  },
};
