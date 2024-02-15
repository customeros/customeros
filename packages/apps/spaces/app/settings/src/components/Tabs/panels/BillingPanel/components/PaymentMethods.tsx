'use client';

import React, { useEffect } from 'react';

import { useConnections, useIntegrationApp } from '@integration-app/react';

import { Flex } from '@ui/layout/Flex';
import { Gb } from '@ui/media/logos/Gb';
import { Us } from '@ui/media/logos/Us';
import { Eu } from '@ui/media/logos/Eu';
import { Text } from '@ui/typography/Text';
import { Stripe } from '@ui/media/logos/Stripe';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { toastError } from '@ui/presentation/Toast';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';

export const PaymentMethods = ({
  onResetPaymentMethods,
  canPayWithCard,
  canPayWithDirectDebitSEPA,
  canPayWithDirectDebitACH,
  canPayWithDirectDebitBacs,
}: {
  canPayWithCard?: boolean | null;
  onResetPaymentMethods: () => void;
  canPayWithDirectDebitACH?: boolean | null;
  canPayWithDirectDebitSEPA?: boolean | null;
  canPayWithDirectDebitBacs?: boolean | null;
}) => {
  const formId = 'tenant-billing-profile-form';

  const iApp = useIntegrationApp();
  const { items: iConnections, refresh, loading } = useConnections();
  const isStripeActive = iConnections
    .map((item) => item.integration?.key)
    .find((e) => e === 'stripe');

  const handleStripe = async (onChange: () => void) => {
    if (isStripeActive) {
      onChange();
    }
    if (!isStripeActive) {
      try {
        await iApp.integration('stripe').open({ showPoweredBy: false });
        await refresh();

        onChange();
        // continue
      } catch (err) {
        toastError('Integration failed', 'get-intergration-data');
      }
    }
  };

  const handleOpenIntegrationAppModal = async () => {
    try {
      await iApp.integration('stripe').open({ showPoweredBy: false });
      await refresh();
      // continue
    } catch (err) {
      toastError('Integration failed', 'get-intergration-data');
    }
  };
  useEffect(() => {
    if (!isStripeActive && !loading && iConnections?.length > 0) {
      onResetPaymentMethods();
    }
  }, [iConnections, isStripeActive, loading]);

  return (
    <>
      <Flex position='relative' alignItems='center'>
        <Text color='gray.500' fontSize='xs' whiteSpace='nowrap' mr={2}>
          Customer can pay using
        </Text>
        <Divider background='gray.200' />
      </Flex>

      <FormSwitch
        name='canPayWithCard'
        formId={formId}
        size='sm'
        onChangeCallback={handleStripe}
        leftElement={
          canPayWithCard && (
            <IconButton
              variant='ghost'
              aria-label='open integration app modal'
              size='xs'
              border='1px solid'
              borderColor='gray.300'
              borderRadius='50%'
              padding={1}
              mr={2}
              onClick={handleOpenIntegrationAppModal}
              icon={<Stripe boxSize={3} />}
            />
          )
        }
        label={
          <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>
            Credit or Debit cards
          </Text>
        }
      />

      <Flex flexDir='column' gap={2}>
        <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>
          Direct debit via
        </Text>
        <FormSwitch
          name='canPayWithDirectDebitSEPA'
          formId={formId}
          size='sm'
          onChangeCallback={handleStripe}
          leftElement={
            canPayWithDirectDebitSEPA && (
              <IconButton
                variant='ghost'
                aria-label='open integration app modal'
                size='xs'
                border='1px solid'
                borderColor='gray.300'
                borderRadius='50%'
                padding={1}
                mr={2}
                onClick={handleOpenIntegrationAppModal}
                icon={<Stripe boxSize={3} />}
              />
            )
          }
          label={
            <Text
              fontSize='sm'
              fontWeight='medium'
              whiteSpace='nowrap'
              as='label'
            >
              <Eu mr={2} />
              SEPA
            </Text>
          }
        />
        <FormSwitch
          name='canPayWithDirectDebitACH'
          formId={formId}
          size='sm'
          onChangeCallback={handleStripe}
          leftElement={
            canPayWithDirectDebitACH && (
              <IconButton
                variant='ghost'
                aria-label='open integration app modal'
                size='xs'
                border='1px solid'
                borderColor='gray.300'
                borderRadius='50%'
                padding={1}
                mr={2}
                onClick={handleOpenIntegrationAppModal}
                icon={<Stripe boxSize={3} />}
              />
            )
          }
          label={
            <Text
              fontSize='sm'
              fontWeight='medium'
              whiteSpace='nowrap'
              as='label'
            >
              <Us mr={2} />
              ACH
            </Text>
          }
        />

        <FormSwitch
          name='canPayWithDirectDebitBacs'
          formId={formId}
          size='sm'
          onChangeCallback={handleStripe}
          leftElement={
            canPayWithDirectDebitBacs && (
              <IconButton
                variant='ghost'
                aria-label='open integration app modal'
                size='xs'
                border='1px solid'
                borderColor='gray.300'
                borderRadius='50%'
                padding={1}
                mr={2}
                onClick={handleOpenIntegrationAppModal}
                icon={<Stripe boxSize={3} />}
              />
            )
          }
          label={
            <Text
              fontSize='sm'
              fontWeight='medium'
              whiteSpace='nowrap'
              as='label'
            >
              <Gb mr={2} />
              Bacs
            </Text>
          }
        />
      </Flex>
      {/*<Flex justifyContent='space-between' alignItems='center'>*/}
      {/*  <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>*/}
      {/*    Bank transfer*/}
      {/*  </Text>*/}
      {/*  <Switch size='sm' />*/}
      {/*</Flex>*/}
      {/*<FormSwitch*/}
      {/*  name='canPayWithPigeon'*/}
      {/*  formId={formId}*/}
      {/*  size='sm'*/}
      {/*  label={*/}
      {/*    <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>*/}
      {/*      Carrier pigeon*/}
      {/*    </Text>*/}
      {/*  }*/}
      {/*/>*/}
    </>
  );
};
