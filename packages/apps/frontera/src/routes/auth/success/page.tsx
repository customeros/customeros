import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';

export const SuccessPage = observer(() => {
  const navigate = useNavigate();
  const store = useStore();

  useEffect(() => {
    if (store.session.isHydrated && store.session.isAuthenticated) {
      const originPath = new URLSearchParams(window.location.search).get(
        'origin',
      );

      store.session.fetchSession({
        onSuccess: () =>
          setTimeout(() => navigate(originPath ?? '/organizations'), 500),
      });
    }
  }, [store.session.isHydrated, store.session.isAuthenticated]);

  return (
    <div
      className={cn(
        'absolute bg-white flex flex-col items-center justify-center top-0 right-0 bottom-0 left-0 z-10 opacity-100 transition-opacity duration-500',
        store.session.isBootstrapped && 'opacity-0',
      )}
    >
      <Spinner label='loading' className='text-gray-300 fill-gray-500' />
    </div>
  );
});
