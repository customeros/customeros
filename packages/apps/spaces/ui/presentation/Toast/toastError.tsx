import { toast } from 'react-toastify';
import { XClose } from '@ui/media/icons/XClose';
import ExclamationWaves from '@spaces/atoms/icons/ExclamationWaves';

export const toastError = (text: string, id: string) => {
  return toast.error(text, {
    toastId: id,
    icon: ExclamationWaves,
    closeButton: ({ closeToast }) => (
      <div onClick={closeToast}>
        <XClose height='30px' width='30px' color='#F04438' />
      </div>
    ),
  });
};
