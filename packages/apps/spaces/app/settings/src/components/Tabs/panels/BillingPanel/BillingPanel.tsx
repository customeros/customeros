'use client';

import { useForm } from 'react-inverted-form';
import React, { useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';
import { useDebounce, useDeepCompareEffect } from 'rooks';
import { useUpdateTenantSettingsMutation } from '@settings/graphql/updateTenantSettings.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';
import { useCreateBillingProfileMutation } from '@settings/graphql/createTenantBillingProfile.generated';
import { useTenantUpdateBillingProfileMutation } from '@settings/graphql/updateTenantBillingProfile.generated';
import {
  TenantSettingsQuery,
  useTenantSettingsQuery,
} from '@settings/graphql/getTenantSettings.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Collapse } from '@ui/transitions/Collapse';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { SlashOctagon } from '@ui/media/icons/SlashOctagon';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';
import { DataSource, InvoiceLine, TenantBillingProfile } from '@graphql/types';

import { TenantBillingPanelDetailsForm } from './components';
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

  const tenantBillingProfileId = data?.tenantBillingProfiles?.[0]?.id ?? '';
  const queryKey = useTenantBillingProfilesQuery.getKey();
  const settingsQueryKey = useTenantSettingsQuery.getKey();

  const createBillingProfileMutation = useCreateBillingProfileMutation(client, {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });
  const { data: tenantSettingsData } = useTenantSettingsQuery(client);

  const updateTenantSettingsMutation = useUpdateTenantSettingsMutation(client, {
    onMutate: ({ input: { ...newSettings } }) => {
      queryClient.cancelQueries({ queryKey: settingsQueryKey });
      const previousEntries =
        queryClient.getQueryData<TenantSettingsQuery>(settingsQueryKey);
      queryClient.setQueryData(settingsQueryKey, {
        tenantSettings: {
          ...(previousEntries?.tenantSettings ?? {}),
          ...newSettings,
        },
      });

      return { previousSettings: previousEntries };
    },
    onError: (err, newSettings, context) => {
      queryClient.setQueryData(settingsQueryKey, context?.previousSettings);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: settingsQueryKey });
    },
  });
  const updateBillingProfileMutation = useTenantUpdateBillingProfileMutation(
    client,
    {
      onMutate: ({ input: { patch, ...restInput } }) => {
        queryClient.cancelQueries({ queryKey });

        useTenantBillingProfilesQuery.mutateCacheEntry(queryClient)(
          (cacheEntry) => {
            return produce(cacheEntry, (draft) => {
              const selectedProfile = draft?.tenantBillingProfiles?.findIndex(
                (profileId) =>
                  profileId.id === data?.tenantBillingProfiles?.[0]?.id,
              );
              if (
                selectedProfile >= 0 &&
                draft?.tenantBillingProfiles?.[selectedProfile]
              ) {
                draft.tenantBillingProfiles[selectedProfile] = {
                  ...draft.tenantBillingProfiles[selectedProfile],
                  ...(restInput as TenantBillingProfile),
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
          subtotal: 100,
          createdAt: new Date().toISOString(),
          metadata: {
            id: 'dummy-id',
            created: new Date().toISOString(),
            lastUpdated: new Date().toISOString(),
            source: DataSource.Openline,
            sourceOfTruth: DataSource.Openline,
            appSource: DataSource.Openline,
          },
          description: 'Professional tier',
          price: 50,
          quantity: 2,
          total: 100,
          taxDue: 0,
        } as InvoiceLine,
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

  const defaultValues = new TenantBillingDetailsDto(
    data?.tenantBillingProfiles?.[0] as TenantBillingProfile,
  );

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
    defaultValues,
    stateReducer: (_, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'canPayWithDirectDebitSEPA':
          case 'canPayWithDirectDebitACH':
          case 'canPayWithDirectDebitBacs':
          case 'canPayWithCard':
          case 'canPayWithPigeon': {
            updateBillingProfileMutation.mutate({
              input: {
                id: tenantBillingProfileId,
                patch: true,
                [action.payload.name]: action.payload.value,
              },
            });

            return next;
          }
          case 'country': {
            updateBillingProfileMutation.mutate({
              input: {
                id: tenantBillingProfileId,
                patch: true,
                [action.payload.name]: action.payload.value?.value,
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
    if (isFetchedAfterMount && !data?.tenantBillingProfiles?.length) {
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

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

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

  const handleToggleInvoices = () => {
    updateTenantSettingsMutation.mutate({
      input: {
        patch: true,
        billingEnabled: !tenantSettingsData?.tenantSettings?.billingEnabled,
      },
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
        <CardHeader
          px='6'
          pb='0'
          pt='4'
          as={Flex}
          alignItems='center'
          justifyContent='space-between'
        >
          <Heading as='h1' fontSize='lg' color='gray.700' pt={1}>
            <b>Billing</b>
          </Heading>

          {tenantSettingsData?.tenantSettings.billingEnabled && (
            <Menu>
              <MenuButton
                as={IconButton}
                size='xs'
                aria-label='Options'
                icon={<DotsVertical />}
                variant='outline'
                border='none'
              />
              <MenuList>
                <MenuItem
                  alignItems='center'
                  color='gray.700'
                  onClick={handleToggleInvoices}
                >
                  <SlashOctagon marginRight={1} color='gray.500' /> Disable
                  Customer billing
                </MenuItem>
              </MenuList>
            </Menu>
          )}
        </CardHeader>

        {!tenantSettingsData?.tenantSettings.billingEnabled && (
          <CardBody
            as={Flex}
            flexDir='column'
            px='6'
            w='full'
            gap={4}
            opacity={tenantSettingsData?.tenantSettings.billingEnabled ? 0 : 1}
          >
            <Text fontSize='sm'>
              Master your revenue lifecycle from contract to cash by enabling
              customer billing for your customers.
            </Text>

            <Box as='ul' pl={6} fontSize='sm'>
              <li>
                Automatically send customer invoices based on their contract
                service line items
              </li>
              <li>Let customers pay using a connected payment provider</li>
            </Box>
            <Flex alignItems='center'>
              <Button
                colorScheme='primary'
                variant='outline'
                size='sm'
                onClick={handleToggleInvoices}
              >
                Enable invoicing
              </Button>
            </Flex>
          </CardBody>
        )}

        <Collapse
          delay={{ enter: 0.2 }}
          in={tenantSettingsData?.tenantSettings.billingEnabled}
          animateOpacity
          startingHeight={0}
        >
          <TenantBillingPanelDetailsForm
            email={state.values.email}
            formId={formId}
            canPayWithCard={state.values.canPayWithCard}
            invoicingEnabled={tenantSettingsData?.tenantSettings.billingEnabled}
            canPayWithDirectDebitACH={state.values.canPayWithDirectDebitACH}
            canPayWithDirectDebitSEPA={state.values.canPayWithDirectDebitSEPA}
            canPayWithDirectDebitBacs={state.values.canPayWithDirectDebitBacs}
            onDisablePaymentMethods={handleDisablePaymentMethods}
            setIsInvoiceProviderFocused={setIsInvoiceProviderFocused}
            setIsInvoiceProviderDetailsHovered={
              setIsInvoiceProviderDetailsHovered
            }
          />
        </Collapse>
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
