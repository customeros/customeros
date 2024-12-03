import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

const formatNumberWithComma = (num: number): string => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const CheckoutCard = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

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
  const total = (199.99 + store.mailboxes.extendedBundle.size * 18.99).toFixed(
    2,
  );

  return (
    <>
      <Card className='py-2 px-3 bg-white mt-2'>
        <CardContent className='p-0'>
          <div className='flex items-center gap-1 bg-gray-100 rounded-lg py-1 px-2 leading-4'>
            <CheckCircle className='size-7 text-gray-500 mr-2' />
            <p className='text-sm'>
              With{' '}
              <span className='font-medium'>{`${noOfMailboxes} mailboxes`}</span>{' '}
              you can send up to
              <span className='font-medium'> {`${noOfEmails} emails`}</span> per
              month
            </p>
          </div>
        </CardContent>
        <CardFooter className='flex flex-col p-0 mt-3 items-center justify-center'>
          {store.mailboxes.invalidDomains.length > 0 && (
            <div className='mb-2 bg-error-50 w-full flex items-center gap-2 rounded-lg py-1 px-2'>
              <DotSingle className='text-error-500 size-6' />
              <span className='text-error-700 text-sm'>
                {store.mailboxes.invalidDomains.length} of your domains are
                unavailable
              </span>
            </div>
          )}
          {!!store.mailboxes.invalidBaseBundle && (
            <div className='mb-2 bg-error-50 w-full flex items-center gap-2 rounded-lg py-1 px-2'>
              <DotSingle className='text-error-500 size-6' />
              <span className='text-error-700 text-sm'>
                {store.mailboxes.invalidBaseBundle}
              </span>
            </div>
          )}

          <Button
            className='w-full'
            colorScheme='primary'
            rightIcon={<ChevronRight />}
            isLoading={store.mailboxes.isLoading}
            onClick={() => {
              handlePaymentView();
            }}
          >
            {`Checkout: $${total}/month`}
          </Button>
        </CardFooter>
      </Card>
    </>
  );
});
