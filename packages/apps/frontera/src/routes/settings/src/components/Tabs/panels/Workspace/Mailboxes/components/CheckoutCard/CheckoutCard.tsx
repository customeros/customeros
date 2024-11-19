import { useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Button } from '@ui/form/Button/Button';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

const formatNumberWithComma = (num: number): string => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

export const CheckoutCard = () => {
  const [searchParams, setSearchParams] = useSearchParams();

  const [storedBrandName, _setStoredBrandName] = useLocalStorage<string[]>(
    'brandName',
    [],
  );

  const [storeUserName, _setStoreUserName] = useLocalStorage<string[]>(
    'userName',
    [],
  );

  const [selectedAdditionalDomains] = useLocalStorage<string[]>(
    'selectedAdditionalDomains',
    [],
  );

  const handlePaymentView = () => {
    const params = new URLSearchParams(searchParams.toString() ?? '');

    params.set('checkout', 'mailboxes');
    setSearchParams(params.toString());
  };

  const noOfMailboxes =
    storedBrandName.length +
    selectedAdditionalDomains.length * storeUserName.length;

  const noOfEmails = formatNumberWithComma(noOfMailboxes * 1200);

  const total = (199.99 + selectedAdditionalDomains.length * 18.99).toFixed(2);

  return (
    <>
      <Card className='py-2 px-3 bg-white mt-2'>
        <CardContent className='p-0'>
          <div className='flex items-center gap-1 bg-gray-50 rounded-lg py-1 px-2 leading-4'>
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
        <CardFooter className='p-0 mt-3 items-center justify-center'>
          <Button
            className='w-full'
            colorScheme='primary'
            rightIcon={<ChevronRight />}
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
};
