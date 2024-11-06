import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Download04 } from '@ui/media/icons/Download04';

export const SystemUpdateNotification = observer(() => {
  const store = useStore();

  const handleReload = () => {
    window.location.reload();
  };

  if (!store.ui.isSystemNotificationOpen) return null;

  return (
    <div className='m-[10px] p-1 flex flex-col items-start gap-1 border rounded-[4px] border-gray-100 bg-gray-50'>
      <div className='flex items-center gap-1 flex-wrap'>
        <Download04 className='text-gray-500' />
        <p className='font-medium text-sm'>New version available</p>
      </div>
      <p className='text-xs'>
        Get an improved version of the app by reloading it.
      </p>
      <Button
        size='xxs'
        variant='ghost'
        colorScheme='primary'
        onClick={handleReload}
        className='ml-auto mt-2'
      >
        Reload app
      </Button>
    </div>
  );
});
