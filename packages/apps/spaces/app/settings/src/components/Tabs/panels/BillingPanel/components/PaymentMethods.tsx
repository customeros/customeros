'use client';

import React from 'react';

import { useConnections, useIntegrationApp } from '@integration-app/react';
import { paymentMethods } from '@settings/components/Tabs/panels/BillingPanel/components/utils';
import { useGetExternalSystemInstancesQuery } from '@settings/graphql/getExternalSystemInstances.generated';
import { BankTransferAccountList } from '@settings/components/Tabs/panels/BillingPanel/components/BankTransferAccountList';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Switch } from '@ui/form/Switch';
import { FormLabel } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Stripe } from '@ui/media/logos/Stripe';
import { Divider } from '@ui/presentation/Divider';
import { ExternalSystemType } from '@graphql/types';
import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

export const PaymentMethods = ({
  formId,
  organizationName,
}: {
  formId: string;
  organizationName?: string | null;
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
      <Flex position='relative' alignItems='center'>
        <Text fontSize='sm' whiteSpace='nowrap' mr={2} color='gray.500'>
          Customer can pay using
        </Text>
        <Divider background='gray.200' />
      </Flex>
      <Box w='full'>
        <FormLabel
          display='flex'
          alignItems='center'
          justifyContent='space-between'
          m={0}
        >
          <Text fontSize='sm' whiteSpace='nowrap'>
            <Stripe boxSize={5} mr={2} />
            Stripe
          </Text>
          <Switch
            size='sm'
            isInvalid={!availablePaymentMethodTypes?.length}
            isChecked={isStripeActive}
            colorScheme='primary'
            onChange={handleOpenIntegrationAppModal}
          />
        </FormLabel>
        {isStripeActive && (
          <Text
            textTransform='capitalize'
            color='gray.500'
            fontSize='sm'
            ml={7}
            noOfLines={1}
          >
            {availablePaymentMethodTypes?.length
              ? availablePaymentMethodTypes
                  ?.map((e) => paymentMethods?.[e])
                  .filter(Boolean)
                  .join(', ')
              : 'No payment methods enabled in Stripe yet'}
          </Text>
        )}
      </Box>

      <BankTransferAccountList
        formId={formId}
        organizationName={organizationName}
      />
    </>
  );
};
