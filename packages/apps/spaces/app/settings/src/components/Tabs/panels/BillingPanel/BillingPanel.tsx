'use client';

import { useForm } from 'react-inverted-form';
import React, { useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { VatInput } from '@settings/components/Tabs/panels/BillingPanel/VatInput';
import { PaymentMethods } from '@settings/components/Tabs/panels/BillingPanel/PaymentMethods';
import { LogoUploader } from '@settings/components/LogoUploadComponent/LogoUploader';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';
import { useCreateBillingProfileMutation } from '@settings/graphql/createTenantBillingProfile.generated';
import { useTenantUpdateBillingProfileMutation } from '@settings/graphql/updateTenantBillingProfile.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { TenantBillingProfile } from '@graphql/types';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  TenantBillingDetails,
  TenantBillingDetailsDto,
} from './TenantBillingProfile.dto';
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const BillingPanel = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const { data, isFetchedAfterMount } = useTenantBillingProfilesQuery(client);
  const [isInvoiceProviderFocused, setIsInvoiceProviderFocused] =
    useState<boolean>(false);
  const [isInvoiceProviderDetailsHovered, setIsInvoiceProviderDetailsHovered] =
    useState<boolean>(false);

  const tenantBillingProfileId = data?.tenantBillingProfiles?.[0]?.id ?? '';
  const queryKey = useTenantBillingProfilesQuery.getKey();

  const createBillingProfileMutation = useCreateBillingProfileMutation(client, {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });
  const updateBillingProfileMutation = useTenantUpdateBillingProfileMutation(
    client,
    {
      onMutate: ({ input }) => {
        queryClient.cancelQueries({ queryKey });

        useTenantBillingProfilesQuery.mutateCacheEntry(queryClient)(
          (cacheEntry) => {
            return produce(cacheEntry, (draft) => {
              const selectedProfile = draft?.tenantBillingProfiles?.findIndex(
                (profileId) =>
                  profileId.id === data?.tenantBillingProfiles?.[0]?.id,
              );

              if (
                selectedProfile &&
                draft?.tenantBillingProfiles?.[selectedProfile]
              ) {
                draft.tenantBillingProfiles[selectedProfile] = {
                  ...draft.tenantBillingProfiles[selectedProfile],
                  ...(input as TenantBillingProfile),
                };
              }
            });
          },
        );
      },
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey });
      },
    },
  );
  const formId = 'tenant-billing-profile-form';
  const invoicePreviewStaticData = useMemo(
    () => ({
      status: 'Preview',
      invoiceNumber: 'INV-003',
      lines: [
        {
          amount: 100,
          createdAt: new Date().toISOString(),
          id: 'dummy-id',
          name: 'Professional tier',
          price: 50,
          quantity: 2,
          totalAmount: 100,
          vat: 0,
        },
      ],
      tax: 0,
      note: '',
      total: 100,
      dueDate: new Date().toISOString(),
      subtotal: 100,
      issueDate: new Date().toISOString(),
      billedTo: {
        addressLine1: '29 Maple Lane',
        addressLine2: 'Springfield, Haven County',
        locality: 'San Francisco',
        zip: '89302',
        country: 'United States',
        email: 'invoices@acme.com',
        name: 'Acme Corp.',
      },
    }),
    [],
  );

  const newDefaults = new TenantBillingDetailsDto();

  const handleUpdateData = useDebounce((d: TenantBillingDetails) => {
    const payload = TenantBillingDetailsDto.toPayload(d);
    updateBillingProfileMutation.mutate({
      input: {
        id: tenantBillingProfileId,
        ...payload,
      },
    });
  }, 2500);
  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues: newDefaults,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'country':
          case 'canPayWithDirectDebitSEPA':
          case 'canPayWithDirectDebitACH':
          case 'canPayWithDirectDebitBacs':
          case 'canPayWithCard':
          case 'canPayWithPigeon': {
            const payload = TenantBillingDetailsDto.toPayload(next.values);
            updateBillingProfileMutation.mutate({
              input: {
                id: tenantBillingProfileId,
                ...payload,
              },
            });

            return next;
          }
          case 'vatNumber':
          case 'sendInvoicesFrom':
          case 'organizationLegalName':
          case 'addressLine1':
          case 'addressLine2':
          case 'addressLine3':
          case 'zip':
          case 'locality': {
            handleUpdateData({
              ...next.values,
            });

            return next;
          }
          default:
            return next;
        }
      }
      if (action.type === 'FIELD_BLUR') {
        switch (action.payload.name) {
          case 'vatNumber':
          case 'sendInvoicesFrom':
          case 'organizationLegalName':
          case 'addressLine1':
          case 'addressLine2':
          case 'addressLine3':
          case 'zip':
          case 'locality': {
            handleUpdateData.flush();

            return next;
          }
          default:
            return next;
        }
      }

      return next;
    },
  });

  useEffect(() => {
    return handleUpdateData.flush();
  }, []);

  useEffect(() => {
    if (isFetchedAfterMount && !data?.tenantBillingProfiles.length) {
      createBillingProfileMutation.mutate({
        input: {
          canPayWithDirectDebitACH: false,
          canPayWithDirectDebitSEPA: false,
          canPayWithDirectDebitBacs: false,
          canPayWithCard: false,
          canPayWithPigeon: false,
          sendInvoicesFrom: '',
          vatNumber: '',
        },
      });
    }
  }, [isFetchedAfterMount, data]);

  useEffect(() => {
    if (
      isFetchedAfterMount &&
      !!data?.tenantBillingProfiles.length &&
      data?.tenantBillingProfiles?.[0]
    ) {
      const newDefaults = new TenantBillingDetailsDto(
        data?.tenantBillingProfiles?.[0] as TenantBillingProfile,
      );
      setDefaultValues(newDefaults);
    }
  }, [isFetchedAfterMount, data]);

  const handleDisablePaymentMethods = () => {
    updateBillingProfileMutation.mutate({
      input: {
        id: tenantBillingProfileId,
        canPayWithDirectDebitACH: false,
        canPayWithDirectDebitSEPA: false,
        canPayWithDirectDebitBacs: false,
        canPayWithCard: false,
        canPayWithPigeon: false,
        patch: true,
      },
    });
    const newDefaults = new TenantBillingDetailsDto(
      data?.tenantBillingProfiles?.[0] as TenantBillingProfile,
    );
    setDefaultValues({
      ...newDefaults,
      canPayWithDirectDebitACH: false,
      canPayWithDirectDebitSEPA: false,
      canPayWithDirectDebitBacs: false,
      canPayWithCard: false,
      canPayWithPigeon: false,
    });
  };

  return (
    <Flex>
      <Card
        flex='1'
        w='full'
        h='100vh'
        bg='#FCFCFC'
        flexDirection='column'
        boxShadow='none'
        background='gray.25'
        maxW={400}
        borderRight='1px solid'
        borderColor='gray.300'
        overflowY='scroll'
        borderRadius='none'
      >
        <CardHeader px='6' pb='0' pt='4'>
          <Heading as='h1' fontSize='lg' color='gray.700' pt={1}>
            <b>Billing</b>
          </Heading>
        </CardHeader>
        <CardBody as={Flex} flexDir='column' px='6' w='full' gap={4}>
          <LogoUploader />
          <FormInput
            autoComplete='off'
            label='Organization legal name'
            placeholder='Legal name'
            isLabelVisible
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            name='legalName'
            formId={formId}
            onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
            onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
            onFocus={() => setIsInvoiceProviderFocused(true)}
            onBlur={() => setIsInvoiceProviderFocused(false)}
          />
          <Flex
            flexDir='column'
            onMouseEnter={() => setIsInvoiceProviderDetailsHovered(true)}
            onMouseLeave={() => setIsInvoiceProviderDetailsHovered(false)}
          >
            <FormInput
              autoComplete='off'
              label='Billing address'
              placeholder='Address line 1'
              isLabelVisible
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                fontWeight: 'semibold',
              }}
              name='addressLine1'
              formId={formId}
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />
            <FormInput
              autoComplete='off'
              label='Billing address line 2'
              name='addressLine2'
              placeholder='Address line 2'
              formId={formId}
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />

            <Flex gap={2}>
              <FormInput
                autoComplete='off'
                label='Billing address locality'
                name='locality'
                placeholder='City'
                formId={formId}
                onFocus={() => setIsInvoiceProviderFocused(true)}
                onBlur={() => setIsInvoiceProviderFocused(false)}
              />
              <FormInput
                autoComplete='off'
                label='Billing address zip/Postal code'
                name='zip'
                placeholder='ZIP/Postal code'
                formId={formId}
                onFocus={() => setIsInvoiceProviderFocused(true)}
                onBlur={() => setIsInvoiceProviderFocused(false)}
              />
            </Flex>
            <FormSelect
              name='country'
              placeholder='Country'
              formId={formId}
              options={countryOptions}
            />
            <VatInput
              formId={formId}
              name='vatNumber'
              autoComplete='off'
              label='VAT number'
              isLabelVisible
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                mt: 4,
                fontWeight: 'semibold',
              }}
              textOverflow='ellipsis'
              placeholder='VAT number'
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />

            <FormInput
              autoComplete='off'
              label='Send invoice from'
              isLabelVisible
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                mt: 4,
                fontWeight: 'semibold',
              }}
              formId={formId}
              name='email'
              textOverflow='ellipsis'
              placeholder='Email'
              type='email'
              isInvalid={
                !!state.values.email?.length &&
                !emailRegex.test(state.values.email)
              }
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />
          </Flex>

          <PaymentMethods
            canPayWithCard={state.values.canPayWithCard}
            canPayWithDirectDebitACH={state.values.canPayWithDirectDebitACH}
            canPayWithDirectDebitSEPA={state.values.canPayWithDirectDebitSEPA}
            canPayWithDirectDebitBacs={state.values.canPayWithDirectDebitBacs}
            onResetPaymentMethods={handleDisablePaymentMethods}
          />
          {/*<Flex justifyContent='space-between' alignItems='center'>*/}
          {/*  <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>*/}
          {/*    Bank transfer*/}
          {/*  </Text>*/}
          {/*  <Switch size='sm' />*/}
          {/*</Flex>*/}
          <FormSwitch
            name='canPayWithPigeon'
            formId={formId}
            size='sm'
            label={
              <Text fontSize='sm' fontWeight='semibold' whiteSpace='nowrap'>
                Carrier pigeon
              </Text>
            }
          />
        </CardBody>
      </Card>
      <Box borderRight='1px solid' borderColor='gray.300' maxH='100vh'>
        <Invoice
          isInvoiceProviderFocused={
            isInvoiceProviderFocused || isInvoiceProviderDetailsHovered
          }
          from={{
            addressLine1: state.values.addressLine1 ?? '',
            addressLine2: state.values.addressLine2 ?? '',
            locality: state.values.locality ?? '',
            zip: state.values.zip ?? '',
            country: state?.values?.country?.label ?? '',
            email: state?.values?.email ?? '',
            name: state.values?.legalName ?? '',
            vatNumber: state.values?.vatNumber ?? '',
          }}
          {...invoicePreviewStaticData}
        />
      </Box>
    </Flex>
  );
};
