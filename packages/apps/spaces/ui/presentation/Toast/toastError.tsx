import { toast } from 'react-toastify';
import Times from '@spaces/atoms/icons/Times';
import ExclamationWaves from '@spaces/atoms/icons/ExclamationWaves';

export const toastError = (text: string, id: string) => {
  return toast.error(text, {
    toastId: id,
    icon: ExclamationWaves,
    closeButton: ({ closeToast }) => (
      <div onClick={closeToast}>
        <Times height={30} width={30} color='#F04438' />
      </div>
    ),
  });
};
