import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

import { MailboxTable } from './components/MailboxTable';
import { CheckoutPage } from './components/CheckoutPage';
import { CheckoutCard } from './components/CheckoutCard';
import { BaseBundleCard } from './components/BaseBundleCard';
import { EmptyMailboxes } from './components/EmptyMailboxes';
import { AddDomainsCard } from './components/AddDomainsCard';
import { UsersCard } from './components/UsersCard/UsersCard';
import { RedirectUrlCard } from './components/RedirectUrlCard';
import { ExtendedBundleCard } from './components/ExtendedBundleCard';

export const Mailboxes = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const campaign = searchParams.get('campaign');

  const hasMailboxes = store.mailboxes.value.size > 0;

  const goBuy = () => {
    if (campaign) return;
    navigate('/settings?tab=mailboxes&view=buy');
  };

  const goToMailboxes = () => {
    if (campaign) return;
    navigate('/settings?tab=mailboxes');
    store.mailboxes.resetBuyFlow();
  };

  const noOfDomains =
    store.mailboxes.baseBundle.size + store.mailboxes.extendedBundle.size;

  const showBuyFlow = searchParams.get('view') === 'buy';
  const showCheckout = searchParams.get('view') === 'checkout';
  const showMailboxes = !searchParams.get('view');

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const alertUser = (e: any) => {
      if (showMailboxes) return;
      e.preventDefault();
      e.returnValue = '';
    };

    window.addEventListener('beforeunload', alertUser);

    return () => {
      window.removeEventListener('beforeunload', alertUser);
    };
  }, [showMailboxes]);

  if (!campaign && !hasMailboxes && showMailboxes) {
    return <EmptyMailboxes onUpdate={goBuy} />;
  }

  if (showMailboxes) {
    return <MailboxTable />;
  }

  return (
    <div className='overflow-y-auto h-full'>
      <div className='grid grid-cols-2 gap-2 max-w-[800px] h-full'>
        {showBuyFlow && (
          <>
            <div className='pb-[12px] pt-[8px] px-6 flex flex-col border-r-[1px]'>
              <div className='flex items-center justify-start gap-1 mb-4'>
                <span
                  onClick={goToMailboxes}
                  className={cn(
                    'text-gray-500 font-semibold transition-colors',
                    !campaign && 'cursor-pointer hover:text-gray-700',
                  )}
                >
                  Mailboxes
                </span>
                <ChevronRight className='mt-0.5 text-gray-400 size-3' />
                <span className='font-semibold cursor-default'>Add new</span>
              </div>
              <div className='space-y-4'>
                <AddDomainsCard />
                <RedirectUrlCard />
                <UsersCard />
              </div>
            </div>
            {noOfDomains > 0 && (
              <div className='pb-[12px] pt-[8px] px-6 flex flex-col h-full border-r-[1px]'>
                <p className='mb-4 font-semibold'>Checkout</p>
                <div className='flex flex-col gap-2'>
                  <BaseBundleCard />
                  {noOfDomains >= 5 && <ExtendedBundleCard />}
                  {noOfDomains > 0 && <CheckoutCard />}
                </div>
              </div>
            )}
          </>
        )}
        {showCheckout && <CheckoutPage />}
      </div>
    </div>
  );
});
