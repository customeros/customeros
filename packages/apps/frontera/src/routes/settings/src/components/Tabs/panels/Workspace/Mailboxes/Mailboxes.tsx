import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

import { CheckoutPage } from './components/CheckoutPage';
import { EmptyMailboxes } from './components/EmptyMailboxes';
import { AddDomainsCard } from './components/AddDomainsCard';
import { UsersCard } from './components/UsersCard/UsersCard';
import { BaseBoundleCard } from './components/BaseBoundleCard';
import { CheckoutCard } from './components/CheckoutCard/CheckoutCard';
import { AdditionalDomainsCard } from './components/AdditionalDomainsCard';

export const Mailboxes = observer(() => {
  const store = useStore();
  const [isUpdated, setIsUpdated] = useState<boolean>(true);
  const [unvailbleDomains, setUnavailbleDomains] = useState<string[]>([]);
  const [checkout, setCheckout] = useState<boolean>(false);
  const [storedBrandName] = useLocalStorage<string[]>('brandName', []);
  const [selectedAdditionalDomains] = useLocalStorage<string[]>(
    'selectedAdditionalDomains',
    [],
  );
  const [searchParams] = useSearchParams();

  const mailboxesStore = store.mailboxes;

  const handleUpdate = () => {
    setIsUpdated(true);
  };

  const noOfDomains = storedBrandName.length;
  const showMainContent = searchParams.get('checkout') !== 'mailboxes';

  if (!isUpdated) {
    return <EmptyMailboxes onUpdate={handleUpdate} />;
  }

  useEffect(() => {
    const checkUnvalidDomains = async () => {
      if (checkout) {
        const response = await mailboxesStore.getMailstackCheckUnvalidDomains({
          domains: storedBrandName.concat(selectedAdditionalDomains),
        });

        setUnavailbleDomains(response || []);
      }
    };

    checkUnvalidDomains();
  }, [checkout, mailboxesStore, storedBrandName, selectedAdditionalDomains]);

  return (
    <div className='overflow-y-auto h-full'>
      <div className='grid grid-cols-2 gap-2 max-w-[800px] h-full'>
        {showMainContent ? (
          <>
            <div className='py-[10px] px-6 flex flex-col border-r-[1px]'>
              <div className='flex items-center justify-start gap-1 mb-4'>
                <span className='text-gray-500 font-semibold'>Mailboxes</span>
                <ChevronRight className='mt-0.5 text-gray-400 size-3' />
                <span className='font-semibold'>Add new</span>
              </div>
              <div className='space-y-4'>
                <AddDomainsCard />
                <UsersCard checkout={checkout} />
              </div>
            </div>
            {noOfDomains > 0 && (
              <div className='py-[10px] px-6 flex flex-col h-full border-r-[1px]'>
                <p className='mb-4 font-semibold'>Checkout</p>
                <div className='flex flex-col gap-2'>
                  <BaseBoundleCard unvalidDomains={unvailbleDomains} />
                  {noOfDomains === 5 && (
                    <AdditionalDomainsCard unvalidDomains={unvailbleDomains} />
                  )}
                  {noOfDomains > 0 && (
                    <CheckoutCard
                      unvalidDomains={unvailbleDomains}
                      onCheckoutClick={() => setCheckout(true)}
                    />
                  )}
                </div>
              </div>
            )}
          </>
        ) : (
          <CheckoutPage />
        )}
      </div>
    </div>
  );
});
