import React from 'react';

import { Image } from '@ui/media/Image/Image';
import { Skeleton } from '@ui/feedback/Skeleton';
import { useStore } from '@shared/hooks/useStore';

import logoCustomerOs from '../../assets/logo-customeros.png';

export const LogoSection = () => {
  const store = useStore();
  const isLoading = store.globalCache?.isLoading;

  return (
    <div className='px-2 pt-2.5 h-fit mb-2 ml-3 cursor-pointer flex justify-flex-start relative'>
      {!isLoading ? (
        <Image
          width={136}
          height={30}
          alt='CustomerOS'
          className='logo-image'
          src={
            store.globalCache.value?.cdnLogoUrl ||
            store.settings.tenant.value?.logoRepositoryFileId ||
            logoCustomerOs
          }
        />
      ) : (
        <Skeleton className='w-full h-8 mr-2' />
      )}
    </div>
  );
};
