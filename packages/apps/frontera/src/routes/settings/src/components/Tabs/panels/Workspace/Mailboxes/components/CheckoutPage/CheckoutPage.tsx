import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { loadStripe, StripeElementsOptions } from '@stripe/stripe-js';
import {
  Elements,
  useStripe,
  useElements,
  PaymentElement,
} from '@stripe/react-stripe-js';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

export const CheckoutForm = observer(() => {
  const store = useStore();
  const stripe = useStripe();
  const elements = useElements();
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string>('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (elements == null) {
      return;
    }

    const { error: submitError } = await elements.submit();

    if (submitError) {
      setErrorMessage(
        submitError.message || 'Looks like we’re having trouble right now',
      );

      return;
    }

    const res = await store.mailboxes.getPaymentIntent();

    if (!stripe) {
      store.ui.toastError('Stripe has not loaded yet', 'stripe-not-loaded-yet');

      return;
    }

    if (!res?.clientSecret) {
      store.ui.toastError(
        'We were unable to start this payment',
        'missing-client-secret-stripe-error',
      );

      return;
    }

    setIsLoading(true);

    const { error, paymentIntent } = await stripe.confirmPayment({
      elements,
      clientSecret: res?.clientSecret,
      redirect: 'if_required',
      confirmParams: {
        return_url: `${window.location.origin}/settings?tab=mailboxes`,
      },
    });

    if (error) {
      store.ui.toastError(
        `We couldn't process your payment`,
        'stripe-processing',
      );
    } else {
      await store.mailboxes.buyDomains(paymentIntent.id);
      store.mailboxes.resetBuyFlow();
      store.ui.toastSuccess('Mailboxes bought', 'mailbox-buy-success');
      navigate('/settings?tab=mailboxes');
    }

    setIsLoading(false);
  };

  return (
    <>
      <form onSubmit={handleSubmit}>
        <PaymentElement
          options={{
            business: { name: 'CustomerOS' },
            terms: { card: 'always' },
          }}
        />
        <Button
          typeof='submit'
          variant='solid'
          colorScheme='blue'
          isLoading={isLoading}
          className='w-full mt-4'
          loadingText='Processing...'
          isDisabled={!stripe || !elements}
        >
          Pay
        </Button>
        {errorMessage && <div>{errorMessage}</div>}
      </form>
    </>
  );
});

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLIC_KEY);

const options: StripeElementsOptions = {
  mode: 'payment',
  currency: 'usd',
  appearance: {
    disableAnimations: true,
    rules: {
      '.Input': {
        width: '600px',
      },
    },
  },
};

export const CheckoutPage = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const campaign = searchParams.get('campaign');
  const campaignParam = campaign ? `&campaign=${campaign}` : '';

  return (
    <div className='py-2 px-6 w-[full] border-r-[1px]'>
      <div className='flex items-center justify-start gap-1 mb-4'>
        <span
          className={cn(
            'font-semibold text-gray-500',
            !campaign && 'hover:cursor-pointer hover:text-gray-700',
          )}
          onClick={() => {
            if (campaign) return;
            navigate('/settings?tab=mailboxes');
            store.mailboxes.resetBuyFlow();
          }}
        >
          Mailboxes
        </span>
        <ChevronRight className='mt-0.5 text-gray-400 size-3' />
        <span
          className='font-semibold text-gray-500 hover:text-gray-700 hover:cursor-pointer'
          onClick={() =>
            navigate('/settings?tab=mailboxes&view=buy' + campaignParam)
          }
        >
          Add new
        </span>
        <ChevronRight className='mt-0.5 text-gray-400 size-3' />
        <span className='font-semibold'>Pay</span>
      </div>
      <Elements
        stripe={stripePromise}
        options={
          {
            ...options,
            amount: store.mailboxes.totalAmount,
          } as StripeElementsOptions
        }
      >
        <CheckoutForm />
      </Elements>
    </div>
  );
});
