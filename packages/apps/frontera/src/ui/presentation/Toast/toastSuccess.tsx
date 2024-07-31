import { toast } from 'react-toastify';

import { XClose } from '@ui/media/icons/XClose';
import { IconButton } from '@ui/form/IconButton';

import CheckWaves from './assets/CheckWaves';

export const toastSuccess = (text: string, id: string) => {
  return toast.success(text, {
    toastId: id,
    icon: CheckWaves,
    closeButton: ({ closeToast }) => (
      <IconButton
        variant='ghost'
        aria-label='Close'
        onClick={closeToast}
        colorScheme='success'
        icon={<XClose className='size-5' />}
      />
    ),
  });
};
