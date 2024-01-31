'use client';

import { useForm } from 'react-inverted-form';
import React, { useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';
import { useCreateBillingProfileMutation } from '@settings/graphql/createTenantBillingProfile.generated';
import { useTenantUpdateBillingProfileMutation } from '@settings/graphql/updateTenantBillingProfile.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { TenantBillingProfile } from '@graphql/types';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  TenantBillingDetails,
  TenantBillingDetailsDto,
} from './TenantBillingProfile.dto';

export const BillingPanel = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const { data, isFetchedAfterMount } = useTenantBillingProfilesQuery(client);
  const [isInvoiceProviderFocused, setIsInvoiceProviderFocused] =
    useState<boolean>(false);
  const [isInvoiceProviderDetailsHovered, setIsInvoiceProviderDetailsHovered] =
    useState<boolean>(false);

  const [
    isDomesticBankingDetailsSectionFocused,
    setisDomesticBankingDetailsSectionFocused,
  ] = useState<boolean>(false);
  const [
    isDomesticBankingDetailsSectionHovered,
    setisDomesticBankingDetailsSectionHovered,
  ] = useState<boolean>(false);
  const [
    isInternationalBankingDetailsSectionFocused,
    setisInternationalBankingDetailsSectionFocused,
  ] = useState<boolean>(false);
  const [
    isInternationalBankingDetailsSectionHovered,
    setisInternationalBankingDetailsSectionHovered,
  ] = useState<boolean>(false);
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
        addressLine: '29 Maple Lane',
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
  }, 500);
  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues: newDefaults,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        if (action.payload.name === 'country') {
          const payload = TenantBillingDetailsDto.toPayload(next.values);
          updateBillingProfileMutation.mutate({
            input: {
              id: tenantBillingProfileId,
              ...payload,
            },
          });

          return {
            ...next,
          };
        }
        if (action.payload.name !== 'country') {
          handleUpdateData({
            ...next.values,
          });

          return {
            ...next,
          };
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
        input: {},
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
          <FormInput
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
              label='Billing address line 2'
              name='addressLine2'
              placeholder='Address line 2'
              formId={formId}
              onFocus={() => setIsInvoiceProviderFocused(true)}
              onBlur={() => setIsInvoiceProviderFocused(false)}
            />

            <Flex>
              <FormInput
                label='Billing address locality'
                name='locality'
                placeholder='City'
                formId={formId}
                onFocus={() => setIsInvoiceProviderFocused(true)}
                onBlur={() => setIsInvoiceProviderFocused(false)}
              />
              <FormInput
                label='Billing address zip/Postal code'
                name='zip'
                placeholder='ZIP/Potal code'
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
          </Flex>

          <FormAutoresizeTextarea
            label='Domestic banking details'
            isLabelVisible
            name='domesticPaymentsBankInfo'
            formId={formId}
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            onMouseEnter={() => setisDomesticBankingDetailsSectionHovered(true)}
            onMouseLeave={() =>
              setisDomesticBankingDetailsSectionHovered(false)
            }
            onFocus={() => setisDomesticBankingDetailsSectionFocused(true)}
            onBlur={() => setisDomesticBankingDetailsSectionFocused(false)}
          />
          <FormAutoresizeTextarea
            label='International banking details'
            isLabelVisible
            name='internationalPaymentsBankInfo'
            formId={formId}
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            onMouseEnter={() =>
              setisInternationalBankingDetailsSectionHovered(true)
            }
            onMouseLeave={() =>
              setisInternationalBankingDetailsSectionHovered(false)
            }
            onFocus={() => setisInternationalBankingDetailsSectionFocused(true)}
            onBlur={() => setisInternationalBankingDetailsSectionFocused(false)}
          />
        </CardBody>
      </Card>
      <Box borderRight='1px solid' borderColor='gray.300' maxH='100vh'>
        <Invoice
          isInvoiceProviderFocused={
            isInvoiceProviderFocused ||
            (isInvoiceProviderDetailsHovered &&
              !isInternationalBankingDetailsSectionFocused &&
              !isDomesticBankingDetailsSectionFocused)
          }
          isDomesticBankingDetailsSectionFocused={
            isDomesticBankingDetailsSectionFocused ||
            (isDomesticBankingDetailsSectionHovered &&
              !isInvoiceProviderFocused &&
              !isInternationalBankingDetailsSectionFocused)
          }
          isInternationalBankingDetailsSectionFocused={
            isInternationalBankingDetailsSectionFocused ||
            (isInternationalBankingDetailsSectionHovered &&
              !isInvoiceProviderFocused &&
              !isDomesticBankingDetailsSectionFocused)
          }
          from={{
            addressLine: state.values.addressLine1 ?? '',
            addressLine2: state.values.addressLine2 ?? '',
            locality: state.values.locality ?? '',
            zip: state.values.zip ?? '',
            country: state?.values?.country?.label ?? '',
            email: '',
            name: state.values?.legalName ?? '',
          }}
          {...invoicePreviewStaticData}
          domesticBankingDetails={state?.values?.domesticPaymentsBankInfo}
          internationalBankingDetails={
            state?.values?.internationalPaymentsBankInfo
          }
        />
      </Box>
    </Flex>
  );
};
