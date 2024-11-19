import { useState } from 'react';

import { loadStripe, StripeElementsOptions } from '@stripe/stripe-js';
import {
  Elements,
  useStripe,
  useElements,
  PaymentElement,
} from '@stripe/react-stripe-js';

import { Input } from '@ui/form/Input';

export const CheckoutForm = () => {
  const stripe = useStripe();
  const elements = useElements();

  const [errorMessage, setErrorMessage] = useState<string | null>('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (elements == null) {
      return;
    }

    // Trigger form validation and wallet collection
    const { error: submitError } = await elements.submit();

    if (submitError) {
      // Show error to your customer
      setErrorMessage(submitError.message || 'An unknown error occurred');

      return;
    }

    // Create the PaymentIntent and obtain clientSecret from your server endpoint
    const res = await fetch('/create-intent', {
      method: 'POST',
    });

    const { client_secret: clientSecret } = await res.json();

    if (!stripe) {
      setErrorMessage('Stripe has not loaded yet.');

      return;
    }

    const { error } = await stripe.confirmPayment({
      //`Elements` instance that was used to create the Payment Element
      elements,
      clientSecret,
      confirmParams: {
        return_url: 'https://example.com/order/123/complete',
      },
    });

    if (error) {
      // This point will only be reached if there is an immediate error when
      // confirming the payment. Show error to your customer (for example, payment
      // details incomplete)
      setErrorMessage(error.message);
    } else {
      // Your customer will be redirected to your `return_url`. For some payment
      // methods like iDEAL, your customer will be redirected to an intermediate
      // site first to authorize the payment, then redirected to the `return_url`.
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit}>
        <label>Email</label>
        <Input placeholder='email' />
        <PaymentElement
          className='mt-2 flex'
          options={{
            paymentMethodOrder: ['card'],
            business: { name: 'CustomerOS' },
            terms: { card: 'always' },
            fields: {
              billingDetails: {
                email: 'auto',
                name: 'auto',
              },
            },
          }}
        />
        <button type='submit' disabled={!stripe || !elements}>
          Pay
        </button>
        {errorMessage && <div>{errorMessage}</div>}
      </form>
    </>
  );
};

const stripePromise = loadStripe('pk_test_6pRNASCoBOKtIshFeQd4XMUh');

const options: StripeElementsOptions = {
  mode: 'payment',
  amount: 1099,
  currency: 'usd',

  // Fully customizable with appearance API.
  appearance: {
    disableAnimations: true,
    rules: {
      AnimateSinglePresence: {
        display: 'none !important',
      },
    },
  },
};
export const CheckoutPage = () => (
  <Elements stripe={stripePromise} options={options as StripeElementsOptions}>
    <CheckoutForm />
  </Elements>
);
