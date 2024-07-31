import { useConnections, useIntegrationApp } from '@integration-app/react';
import { paymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/utils';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';
import { BankTransferAccountList } from '@settings/components/Tabs/panels/BillingPanel/components/BankTransferAccountList';

import { Stripe } from '@ui/media/logos/Stripe';
import { Switch } from '@ui/form/Switch/Switch';
import { ExternalSystemType } from '@graphql/types';
import { toastError } from '@ui/presentation/Toast';
import { Divider } from '@ui/presentation/Divider/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const PaymentMethods = ({
  formId,
  legalName,
}: {
  formId: string;
  legalName?: string | null;
}) => {
  const client = getGraphQLClient();
  const { data } = useGetExternalSystemInstancesQuery(client);
  const iApp = useIntegrationApp();
  const { items: iConnections, refresh } = useConnections();
  const isStripeActive = !!iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');

  const handleOpenIntegrationAppModal = async () => {
    try {
      await iApp.integration('stripe').open({ showPoweredBy: false });
      await refresh();
      await iApp
        .flowInstance({
          flowKey: 'stripe-default-flow-v1',
          integrationKey: 'stripe',
          autoCreate: false,
        })
        .run({
          nodeKey: 'manual-sync-payment-methods',
        });
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };

  const availablePaymentMethodTypes = data?.externalSystemInstances.find(
    (e) => e.type === ExternalSystemType.Stripe,
  )?.stripeDetails?.paymentMethodTypes;

  return (
    <>
      <div className='flex relative items-center'>
        <span className='text-sm whitespace-nowrap mr-2 text-gray-500'>
          Customer can pay using
        </span>
        <Divider />
      </div>
      <div className='w-full'>
        <label className='flex items-center justify-between m-0'>
          <span className='text-sm whitespace-nowrap'>
            <Stripe className='size-5 mr-2' />
            Stripe
          </span>
          <Switch
            size='sm'
            colorScheme='primary'
            isChecked={isStripeActive}
            onChange={handleOpenIntegrationAppModal}
            isInvalid={!availablePaymentMethodTypes?.length}
          />
        </label>
        {isStripeActive && (
          <span className='capitalize text-gray-500 text-sm ml-7 line-clamp-1'>
            {availablePaymentMethodTypes?.length
              ? availablePaymentMethodTypes
                  ?.map((e) => paymentMethods?.[e])
                  .filter(Boolean)
                  .join(', ')
              : 'No payment methods enabled in Stripe yet'}
          </span>
        )}
      </div>

      <BankTransferAccountList formId={formId} legalName={legalName} />
    </>
  );
};
