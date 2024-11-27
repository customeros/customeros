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

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { ChevronRight } from '@ui/media/icons/ChevronRight';

export const CheckoutForm = observer(() => {
  const store = useStore();
  const stripe = useStripe();
  const elements = useElements();

  const [errorMessage, setErrorMessage] = useState<string | null>('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (elements == null) {
      return;
    }

    const { error: submitError } = await elements.submit();

    if (submitError) {
      setErrorMessage(submitError.message || 'An unknown error occurred');

      return;
    }

    const res = await store.mailboxes.getPaymentIntent();

    if (!stripe) {
      setErrorMessage('Stripe has not loaded yet.');

      return;
    }

    if (!res?.clientSecret) {
      // show some error about payment not processable
      return;
    }

    const { error, paymentIntent } = await stripe.confirmPayment({
      elements,
      clientSecret: res?.clientSecret,
      redirect: 'if_required',
      confirmParams: {
        return_url: 'http://localhost:5173/hello',
        receipt_email: 'acalinica@customeros.ai',
      },
    });

    console.log(paymentIntent);

    if (error) {
      store.ui.toastError(
        'Could not process your payment',
        'stripe-processing',
      );
      // This point will only be reached if there is an immediate error when
      // confirming the payment. Show error to your customer (for example, payment
      // details incomplete)
      // setErrorMessage(error?.message);
    } else {
      // Your customer will be redirected to your `return_url`. For some payment
      // methods like iDEAL, your customer will be redirected to an intermediate
      // site first to authorize the payment, then redirected to the `return_url`.
    }
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
          className='w-full mt-4'
          isDisabled={!stripe || !elements}
        >
          Pay
        </Button>
        {errorMessage && <div>{errorMessage}</div>}
      </form>
    </>
  );
});

const stripePromise = loadStripe(
  'pk_test_51NmzLnEVwE7CWhpkM1aC51Y9MDX4FwryNWDfwotBBAodIGkashnVV0HoRAmArnpiOvPjhgbH1IdjXKaeHxLF0BiG00DeAezb3D',
);

const options: StripeElementsOptions = {
  mode: 'payment',
  amount: 1099,
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

export const CheckoutPage = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const goToAddNew = () => {
    const params = new URLSearchParams(searchParams);

    params.delete('checkout');
    setSearchParams(params);
  };

  return (
    <div className='py-2 px-4 w-[full] border-r-[1px]'>
      <div className='flex items-center justify-start gap-1 mb-4'>
        <span
          onClick={() => navigate('/settings?tab=mailboxes')}
          className='font-semibold text-gray-500 hover:text-gray-700 hover:cursor-pointer'
        >
          Mailboxes
        </span>
        <ChevronRight className='mt-0.5 text-gray-400 size-3' />
        <span
          onClick={goToAddNew}
          className='font-semibold text-gray-500 hover:text-gray-700 hover:cursor-pointer'
        >
          Add new
        </span>
        <ChevronRight className='mt-0.5 text-gray-400 size-3' />
        <span className='font-semibold'>Pay</span>
      </div>
      <Elements
        stripe={stripePromise}
        options={options as StripeElementsOptions}
      >
        <CheckoutForm />
      </Elements>
    </div>
  );
};
