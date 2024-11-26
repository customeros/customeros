import { useSearchParams } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Button } from '@ui/form/Button/Button';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { CheckCircle } from '@ui/media/icons/CheckCircle';
import { ChevronRight } from '@ui/media/icons/ChevronRight';
import { Card, CardFooter, CardContent } from '@ui/presentation/Card/Card';

const formatNumberWithComma = (num: number): string => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
};

interface CheckoutCardProps {
  unvalidDomains: string[];
  onCheckoutClick: (isCheckedOut: boolean) => void;
}

export const CheckoutCard = ({
  onCheckoutClick,
  unvalidDomains,
}: CheckoutCardProps) => {
  const [searchParams, setSearchParams] = useSearchParams();

  const [storedBrandName, _setStoredBrandName] = useLocalStorage<string[]>(
    'brandName',
    [],
  );

  const [storedUserName, _setStoredUserName] = useLocalStorage<string[]>(
    'userName',
    [],
  );

  const [selectedAdditionalDomains] = useLocalStorage<string[]>(
    'selectedAdditionalDomains',
    [],
  );

  const handlePaymentView = async () => {
    onCheckoutClick(true);

    if (storedUserName.length >= 1 && unvalidDomains?.length === 0) {
      const params = new URLSearchParams(searchParams.toString() ?? '');

      params.set('checkout', 'mailboxes');
      setSearchParams(params.toString());
    }
  };

  const noOfMailboxes =
    storedBrandName.length +
    selectedAdditionalDomains.length * storedUserName.length;

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
        <CardFooter className='flex flex-col p-0 mt-3 items-center justify-center'>
          {unvalidDomains.length > 0 && (
            <div className='mb-2 bg-error-50 w-full flex items-center gap-2 rounded-lg py-1 px-2'>
              <DotSingle className='text-error-500 size-6' />
              <span className='text-error-700 text-sm'>
                1 of your domains are unavailable
              </span>
            </div>
          )}

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
