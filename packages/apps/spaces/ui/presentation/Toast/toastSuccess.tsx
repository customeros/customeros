import { toast } from 'react-toastify';
import CheckWaves from '@spaces/atoms/icons/CheckWaves';
import { XClose } from '@ui/media/icons/XClose';

export const toastSuccess = (text: string, id: string) => {
  return toast.success(text, {
    toastId: id,
    icon: CheckWaves,
    closeButton: ({ closeToast }) => (
      <div onClick={closeToast}>
        <XClose height='30px' width='30px' color='#17B26A' />
      </div>
    ),
  });
};
