import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { Sale03 } from '@ui/media/icons/Sale03';
import { useStore } from '@shared/hooks/useStore';
import { Inbox01 } from '@ui/media/icons/Inbox01';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

const formatNumberWithComma = (num: number): string => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

interface CheckoutCardProps {
  showTotal?: boolean;
  hideCheckoutButton?: boolean;
}

export const CheckoutCard = observer(
  ({ hideCheckoutButton, showTotal }: CheckoutCardProps) => {
    const store = useStore();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const { open, onOpen, onClose } = useDisclosure();
    const hasMailboxes = store.mailboxes.value.size > 0;

    const handlePaymentView = async () => {
      const campaign = searchParams.get('campaign');

      store.mailboxes.validateBuy({
        onSuccess: () => {
          navigate(
            '/settings?tab=mailboxes&view=checkout' +
              (campaign ? `&campaign=${campaign}` : ''),
          );
        },
      });
    };

    const noOfMailboxes = store.mailboxes.mailboxesCount;
    const noOfEmails = formatNumberWithComma(noOfMailboxes * 1200);
    const total = (
      (hasMailboxes ? 0 : 199.99) +
      store.mailboxes.extendedBundle.size * 18.99
    ).toFixed(2);

    return (
      <>
        <Card
          className={cn(
            'pt-2 pb-2 px-3 bg-white mt-2',
            hideCheckoutButton && 'pb-0',
          )}
        >
          <CardContent className='p-0'>
            <div className='mb-2 bg-gray-100 w-full flex items-center gap-2 rounded-lg py-1 px-2 leading-4'>
              <Inbox01 className='ml-1 text-gray-500 mr-2 size-5 min-w-5' />
              {store.mailboxes.usernamesCount === 0 ? (
                <span className='text-gray-700 text-sm'>
                  Add usernames to see how many emails you can send
                  <InfoCircle
                    onClick={onOpen}
                    className='size-3 ml-1 text-gray-500 cursor-pointer hover:text-gray-700'
                  />
                </span>
              ) : (
                <p className='text-sm'>
                  With{' '}
                  <span className='font-medium'>{`${noOfMailboxes} mailboxes`}</span>{' '}
                  you can send up to
                  <span className='font-medium'>
                    {' '}
                    {`${noOfEmails} emails`}
                  </span>{' '}
                  per month
                  <InfoCircle
                    onClick={onOpen}
                    className='size-3 ml-1 text-gray-500 cursor-pointer hover:text-gray-700'
                  />
                </p>
              )}
            </div>

            {showTotal && (
              <div className='mb-2 bg-success-50 w-full flex items-center gap-2 rounded-lg py-2 px-2 leading-4'>
                <Sale03 className='ml-1 text-success-500 mr-2 size-5 min-w-5' />
                <span className='text-success-700 text-sm'>
                  <span className='font-medium'>Total:</span>{' '}
                  {`$${total}/month`}
                </span>
              </div>
            )}
          </CardContent>
          <CardFooter
            className={cn(
              'flex flex-col p-0 mt-3 items-center justify-center',
              hideCheckoutButton && 'mt-2',
            )}
          >
            {store.mailboxes.invalidDomains.length > 0 && (
              <div className='mb-2 bg-error-50 w-full flex items-center gap-2 rounded-lg py-1 px-2 leading-4'>
                <DotSingle className='text-error-500 size-6 mr-2' />
                <span className='text-error-700 text-sm'>
                  {store.mailboxes.invalidDomains.length} of your domains are
                  unavailable
                </span>
              </div>
            )}

            {!!store.mailboxes.invalidBaseBundle && (
              <div className='mb-2 bg-error-50 w-full flex items-center gap-2 rounded-lg py-1 px-2 leading-4'>
                <DotSingle className='text-error-500 size-6 mr-1' />
                <span className='text-error-700 text-sm'>
                  {store.mailboxes.invalidBaseBundle}
                </span>
              </div>
            )}

            {!hideCheckoutButton && (
              <Button
                className='w-full'
                colorScheme='primary'
                rightIcon={<ChevronRight />}
                isLoading={store.mailboxes.isLoading}
                dataTest='settings-mailboxes-checkout'
                onClick={() => {
                  handlePaymentView();
                }}
              >
                {`Checkout: $${total}/month`}
              </Button>
            )}
          </CardFooter>
        </Card>

        <InfoDialog
          isOpen={open}
          onClose={onClose}
          onConfirm={onClose}
          confirmButtonLabel='Got it'
          label='Mailbox best practices'
          body={
            <div className='space-y-4'>
              <p className='text-sm'>
                To maintain deliverability and avoid spam filters, we auto-warm
                and rotate mailboxes.
              </p>
              <p className='text-sm'>
                We start with 3 emails per day and gradually increase the
                volume, limiting each mailbox to send up to 40 emails per day.
              </p>
              <p className='text-sm'>
                Using two mailboxes per domain has proven especially effective.
              </p>
            </div>
          }
        />
      </>
    );
  },
);
