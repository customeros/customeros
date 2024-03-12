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
  validateEmail,
  validateEmailLocalPart,
} from '@settings/components/Tabs/panels/BillingPanel/utils';
import {
  TenantSettingsQuery,
  useTenantSettingsQuery,
} from '@settings/graphql/getTenantSettings.generated';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { SlashOctagon } from '@ui/media/icons/SlashOctagon';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { Card, CardBody, CardHeader } from '@ui/layout/Card';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';
import {
  DataSource,
  InvoiceLine,
  TenantBillingProfile,
  TenantBillingProfileUpdateInput,
} from '@graphql/types';

import { TenantBillingPanelDetailsForm } from './components';
import { TenantBillingDetailsDto } from './TenantBillingProfile.dto';

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
  const { isOpen, onOpen, onClose } = useDisclosure();

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

  const defaultValues = new TenantBillingDetailsDto({
    ...data?.tenantBillingProfiles?.[0],
    baseCurrency: tenantSettingsData?.tenantSettings?.baseCurrency,
  } as TenantBillingProfile & { baseCurrency: string });

  const handleUpdateData = useDebounce(
    (d: Partial<TenantBillingProfileUpdateInput>) => {
      updateBillingProfileMutation.mutate({
        input: {
          id: tenantBillingProfileId,
          patch: true,
          ...d,
        },
      });
    },
    2500,
  );

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      const getStateAfterValidation = () => {
        return produce(next, (draft) => {
          const sendInvoiceFromError = validateEmailLocalPart(
            draft.values.sendInvoicesFrom,
          );
          const bccError = validateEmail(draft.values.sendInvoicesBcc);
          // we do it like this so that if the email is valid, we reset the states.
          draft.fields.sendInvoicesFrom.meta.hasError = !!sendInvoiceFromError;
          draft.fields.sendInvoicesFrom.error = sendInvoiceFromError ?? '';

          draft.fields.sendInvoicesBcc.meta.hasError = !!bccError;
          draft.fields.sendInvoicesBcc.error = bccError ?? '';
        });
      };
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'canPayWithBankTransfer':
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
          case 'sendInvoicesBcc':
          case 'vatNumber':
          case 'legalName':
          case 'addressLine1':
          case 'addressLine2':
          case 'addressLine3':
          case 'zip':
          case 'locality': {
            handleUpdateData.cancel();
            handleUpdateData({
              [action.payload.name]: action.payload.value,
            });

            return next;
          }

          case 'sendInvoicesFrom': {
            handleUpdateData.cancel();

            handleUpdateData({
              [action.payload
                .name]: `${action.payload.value}@invoices.customeros.ai`,
            });

            return getStateAfterValidation();
          }
          case 'baseCurrency': {
            updateTenantSettingsMutation.mutate({
              input: {
                patch: true,
                baseCurrency: action.payload.value?.value,
              },
            });

            return next;
          }
          default:
            return next;
        }
      }

      if (action.type === 'FIELD_BLUR') {
        setIsInvoiceProviderFocused(false);
        switch (action.payload.name) {
          case 'vatNumber':
          case 'legalName':
          case 'addressLine1':
          case 'addressLine2':
          case 'addressLine3':
          case 'zip':
          case 'locality': {
            handleUpdateData.flush();

            return next;
          }
          case 'sendInvoicesFrom': {
            const formattedEmail = (action.payload?.value || '')
              ?.trim()
              .split(' ')
              .join('-');
            if (!formattedEmail?.length && state.values?.legalName?.length) {
              handleUpdateData.cancel();
              const newEmail = `${state.values.legalName
                .split(' ')
                .join('-')
                .toLowerCase()}@invoices.customeros.ai`;

              updateBillingProfileMutation.mutate({
                input: {
                  id: tenantBillingProfileId,
                  patch: true,
                  sendInvoicesFrom: newEmail,
                },
              });

              return {
                ...next,
                values: {
                  ...next.values,
                  sendInvoicesFrom: `${state.values.legalName
                    .split(' ')
                    .join('-')
                    .toLowerCase()}`,
                },
              };
            } else {
              handleUpdateData.flush();
            }

            return {
              ...getStateAfterValidation(),
              values: {
                ...next.values,
                sendInvoicesFrom: formattedEmail,
              },
            };
          }
          default:
            return next;
        }
      }

      if (action.type === 'SET_DEFAULT_VALUES') {
        return getStateAfterValidation();
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
          canPayWithBankTransfer: true,
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

  const handleDisableBillingDetails = () => {
    updateTenantSettingsMutation.mutate(
      {
        input: {
          patch: true,
          billingEnabled: false,
        },
      },
      {
        onSuccess: onClose,
      },
    );
  };

  const handleToggleInvoices = () => {
    if (!tenantSettingsData?.tenantSettings?.billingEnabled) {
      updateTenantSettingsMutation.mutate({
        input: {
          patch: true,
          billingEnabled: true,
        },
      });

      return;
    }
    onOpen();
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
        minW={400}
        borderRight='1px solid'
        borderColor='gray.300'
        overflowY='scroll'
        borderRadius='none'
        pr={0}
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

        {tenantSettingsData?.tenantSettings.billingEnabled && (
          <TenantBillingPanelDetailsForm
            formId={formId}
            setIsInvoiceProviderFocused={setIsInvoiceProviderFocused}
            setIsInvoiceProviderDetailsHovered={
              setIsInvoiceProviderDetailsHovered
            }
          />
        )}
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
            email: state?.values?.sendInvoicesFrom ?? '',
            name: state.values?.legalName ?? '',
            vatNumber: state.values?.vatNumber ?? '',
          }}
          {...invoicePreviewStaticData}
        />
      </Box>
      <ConfirmDeleteDialog
        label='Disable Customer billing?'
        icon={<SlashOctagon color='error.600' />}
        body='Disabling Customer billing will stop the sending of invoices, and prevent customers from being able to pay.'
        confirmButtonLabel='Disable'
        isOpen={isOpen}
        onClose={onClose}
        onConfirm={handleDisableBillingDetails}
        isLoading={updateTenantSettingsMutation.isPending}
        hideCloseButton
      />
    </Flex>
  );
};
