'use client';

import React, { useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { useDebounce } from 'rooks';
import { useQueryClient } from '@tanstack/react-query';
import { VatInput } from '@settings/components/Tabs/panels/BillingPanel/VatInput';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';
import { useCreateBillingProfileMutation } from '@settings/graphql/createTenantBillingProfile.generated';
import { useTenantUpdateBillingProfileMutation } from '@settings/graphql/updateTenantBillingProfile.generated';

import { Flex } from '@ui/layout/Flex';
import { Eu } from '@ui/media/logos/Eu';
import { Us } from '@ui/media/logos/Us';
import { Gb } from '@ui/media/logos/Gb';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { Divider } from '@ui/presentation/Divider';
import { TenantBillingProfile } from '@graphql/types';
import { FormSwitch } from '@ui/form/Switch/FromSwitch';
import { countryOptions } from '@shared/util/countryOptions';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import {
  TenantBillingDetails,
  TenantBillingDetailsDto,
} from './TenantBillingProfile.dto';
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const TenantBillingPanelDetailsForm = ({
  setIsInvoiceProviderDetailsHovered,
  setIsInvoiceProviderFocused,
}: {
  setIsInvoiceProviderFocused: (newState: boolean) => void;
  setIsInvoiceProviderDetailsHovered: (newState: boolean) => void;
}) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();

  const { data, isFetchedAfterMount } = useTenantBillingProfilesQuery(client);

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

  return (
    <Flex flexDir='column' gap={4}>
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
            !!state.values.email?.length && !emailRegex.test(state.values.email)
          }
          onFocus={() => setIsInvoiceProviderFocused(true)}
          onBlur={() => setIsInvoiceProviderFocused(false)}
        />
      </Flex>

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
    </Flex>
  );
};
