import { toast } from 'react-toastify';
import CheckWaves from '@spaces/atoms/icons/CheckWaves';
import Times from '@spaces/atoms/icons/Times';

export const toastSuccess = (text: string, id: string) => {
  return toast.success(text, {
    toastId: id,
    icon: CheckWaves,
    closeButton: ({ closeToast }) => (
      <div onClick={closeToast}>
        <Times height={30} width={30} color='#17B26A' />
      </div>
    ),
  });
};
