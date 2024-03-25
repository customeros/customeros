'use client';

import React from 'react';

import { useConnections, useIntegrationApp } from '@integration-app/react';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';
import { BankTransferAccountList } from '@settings/components/Tabs/panels/BillingPanel/components/BankTransferAccountList';

import { Stripe } from '@ui/media/logos/Stripe';
import { Switch } from '@ui/form/Switch/Switch2';
import { ExternalSystemType } from '@graphql/types';
import { toastError } from '@ui/presentation/Toast';
import { Divider } from '@ui/presentation/Divider/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const PaymentMethods = ({
  formId,
  organizationName,
}: {
  formId: string;
  organizationName?: string | null;
}) => {
  const client = getGraphQLClient();
  const { data, isLoading } = useGetExternalSystemInstancesQuery(client);
  const iApp = useIntegrationApp();
  const { items: iConnections, refresh } = useConnections();
  const isStripeActive = !!iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');

  const handleOpenIntegrationAppModal = async () => {
    try {
      await iApp.integration('stripe').open({ showPoweredBy: false });
      await refresh();
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };

  const availablePaymentMethodTypes = data?.externalSystemInstances.find(
    (e) => e.type === ExternalSystemType.Stripe,
  )?.stripeDetails?.paymentMethodTypes;

  return (
    <>
      <div className='flex items-center'>
        <span className='text-sm text-gray-500 whitespace-nowrap mr-2'>
          Customer can pay using
        </span>
        <Divider />
      </div>
      <div className='w-full'>
        <label
          htmlFor='Stripe'
          className='flex items-center justify-between mt-0'
        >
          <span className='text-sm whitespace-nowrap'>
            <Stripe boxSize={5} mr={2} />
            Stripe
          </span>
          <Switch
            size='sm'
            isInvalid={!isLoading && !availablePaymentMethodTypes?.length}
            isChecked={isStripeActive}
            colorScheme='primary'
            onChange={handleOpenIntegrationAppModal}
          />
        </label>
        {isStripeActive && (
          <span className='text-sm capitalize text-gray-500 ml-7 line-clamp-1'>
            {availablePaymentMethodTypes?.length
              ? availablePaymentMethodTypes?.join(', ').split('_').join(' ')
              : 'No payment methods enabled in Stripe yet'}
          </span>
        )}
      </div>

      <BankTransferAccountList
        formId={formId}
        organizationName={organizationName}
      />
    </>
  );
};
