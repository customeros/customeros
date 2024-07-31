import { useState, useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { produce } from 'immer';
import { observer } from 'mobx-react-lite';
import { useQueryClient } from '@tanstack/react-query';
import { useDebounce, useDeepCompareEffect } from 'rooks';
import { useUpdateTenantSettingsMutation } from '@settings/graphql/updateTenantSettings.generated';
import { useTenantBillingProfilesQuery } from '@settings/graphql/getTenantBillingProfiles.generated';
import { BillingPanelInvoice } from '@settings/components/Tabs/panels/BillingPanel/BillingPanelInvoice';
import { useCreateBillingProfileMutation } from '@settings/graphql/createTenantBillingProfile.generated';
import { useTenantUpdateBillingProfileMutation } from '@settings/graphql/updateTenantBillingProfile.generated';
import {
  TenantSettingsQuery,
  useTenantSettingsQuery,
} from '@settings/graphql/getTenantSettings.generated';
import { BankTransferSelectionContextProvider } from '@settings/components/Tabs/panels/BillingPanel/context/BankTransferSelectionContext';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { SlashOctagon } from '@ui/media/icons/SlashOctagon';
import { validateEmail } from '@shared/util/emailValidation';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  TenantBillingProfile,
  TenantBillingProfileUpdateInput,
} from '@graphql/types';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

import { TenantBillingPanelDetailsForm } from './components';
import { TenantBillingDetailsDto } from './TenantBillingProfile.dto';

export const BillingPanel = observer(() => {
  const store = useStore();
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
  const { open: isOpen, onOpen, onClose } = useDisclosure();

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
    onError: (_err, _newSettings, context) => {
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
          ({ ...cacheEntry }) => {
            return produce({ ...cacheEntry }, (draft) => {
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
        const bccError = validateEmail(state.values.sendInvoicesBcc);

        return {
          ...next,
          fields: {
            ...next.fields,
            sendInvoicesBcc: {
              ...next.fields.sendInvoicesBcc,
              meta: {
                ...next.fields.sendInvoicesBcc.meta,
                hasError: !!bccError,
              },
              error: bccError ?? '',
            },
          },
        };
      };

      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'check':
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
                region: '',
                [action.payload.name]: action.payload.value?.value,
              },
            });

            return next;
          }
          case 'sendInvoicesBcc':
          case 'sendInvoicesFrom':
          case 'vatNumber':
          case 'legalName':
          case 'addressLine1':
          case 'addressLine2':
          case 'addressLine3':
          case 'region':
          case 'zip':

          case 'locality': {
            handleUpdateData.cancel();
            handleUpdateData({
              [action.payload.name]: action.payload.value,
            });

            return next;
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
          case 'region':
          case 'zip':

          case 'locality': {
            handleUpdateData.flush();

            return next;
          }
          case 'sendInvoicesFrom':

          case 'sendInvoicesBcc': {
            return getStateAfterValidation();
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
          check: true,
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
  const billingEnabledStyle = tenantSettingsData?.tenantSettings.billingEnabled
    ? 'opacity-0'
    : 'opacity-100';

  return (
    <div className='flex'>
      <BankTransferSelectionContextProvider>
        <div className='flex-1 w-full h-[100vh] bg-gray-25 flex-col shadow-none max-w-[400px] min-w-[400px] border-r border-gray-300 overflow-y-scroll pr-0 '>
          <div className='flex items-center justify-between px-6 pb-0 pt-4'>
            <h1 className='text-lg text-gray-700 pt-1'>
              <b>Billing</b>
            </h1>

            {tenantSettingsData?.tenantSettings.billingEnabled && (
              <Menu>
                <MenuButton>
                  <IconButton
                    size='xs'
                    variant='ghost'
                    colorScheme='gray'
                    aria-label='Options'
                    icon={<DotsVertical />}
                  />
                </MenuButton>
                <MenuList>
                  <MenuItem
                    onClick={handleToggleInvoices}
                    className='flex items-center justify-center'
                  >
                    <SlashOctagon className='mr-2 text-gray-500' /> Disable
                    Customer billing
                  </MenuItem>
                </MenuList>
              </Menu>
            )}
          </div>

          {!tenantSettingsData?.tenantSettings.billingEnabled && (
            <div
              className={cn(
                billingEnabledStyle,
                'flex flex-col px-6 w-full gap-4',
              )}
            >
              <span className='text-sm'>
                Master your revenue lifecycle from contract to cash by enabling
                customer billing for your customers.
              </span>

              <ul className='pl-6 text-sm'>
                <li className='list-disc'>
                  Automatically send customer invoices based on their contract
                  service line items
                </li>
                <li className='list-disc'>
                  Let customers pay using a connected payment provider
                </li>
              </ul>
              <div className='items-center'>
                <Button
                  size='sm'
                  variant='outline'
                  colorScheme='primary'
                  isDisabled={store.demoMode}
                  onClick={handleToggleInvoices}
                >
                  Enable Customer billing
                </Button>
              </div>
            </div>
          )}

          {tenantSettingsData?.tenantSettings.billingEnabled && (
            <TenantBillingPanelDetailsForm
              formId={formId}
              country={state.values?.country}
              legalName={state.values?.legalName}
              setIsInvoiceProviderFocused={setIsInvoiceProviderFocused}
              setIsInvoiceProviderDetailsHovered={
                setIsInvoiceProviderDetailsHovered
              }
            />
          )}
        </div>
        <BillingPanelInvoice
          values={state.values}
          isInvoiceProviderFocused={isInvoiceProviderFocused}
          isInvoiceProviderDetailsHovered={isInvoiceProviderDetailsHovered}
        />
      </BankTransferSelectionContextProvider>

      <ConfirmDeleteDialog
        isOpen={isOpen}
        hideCloseButton
        onClose={onClose}
        confirmButtonLabel='Disable'
        label='Disable Customer billing?'
        onConfirm={handleDisableBillingDetails}
        icon={<SlashOctagon color='error.600' />}
        isLoading={updateTenantSettingsMutation.isPending}
        body='Disabling Customer billing will stop the sending of invoices, and prevent customers from being able to pay.'
      />
    </div>
  );
});
