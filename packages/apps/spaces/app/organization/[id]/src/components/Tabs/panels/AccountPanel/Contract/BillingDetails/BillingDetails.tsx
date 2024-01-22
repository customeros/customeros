'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import React, { useRef, useMemo, useState, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Contract } from '@graphql/types';
import { FormInput } from '@ui/form/Input';
import { FeaturedIcon } from '@ui/media/Icon';
import { File02 } from '@ui/media/icons/File02';
import { Grid, GridItem } from '@ui/layout/Grid';
import { Heading } from '@ui/typography/Heading';
import { FormSelect } from '@ui/form/SyncSelect';
import { toastError } from '@ui/presentation/Toast';
import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useUpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

import { countryOptions, currencyOptions } from './utils';
import { BillingDetailsDto, BillingDetailsForm } from './BillingDetails.dto';

interface SubscriptionServiceModalProps {
  data: Contract;
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  organizationName: string;
}

export const BillingDetails = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  data,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const formId = `billing-details-form-${contractId}`;
  const [isBillingDetailsFocused, setIsBillingDetailsFocused] =
    useState<boolean>(false);
  const id = useParams()?.id as string;

  const [isBillingDetailsHovered, setIsBillingDetailsHovered] =
    useState<boolean>(false);
  const queryKey = useGetContractsQuery.getKey({ id });

  const queryClient = useQueryClient();
  const client = getGraphQLClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const updateContract = useUpdateContractMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      queryClient.setQueryData<GetContractsQuery>(queryKey, (currentCache) => {
        return produce(currentCache, (draft) => {
          const previousContracts = draft?.['organization']?.['contracts'];
          const updatedContractIndex = previousContracts?.findIndex(
            (contract) => contract.id === contractId,
          );
          if (draft?.['organization']?.['contracts']) {
            draft['organization']['contracts']?.map((contractData, index) => {
              if (index !== updatedContractIndex) {
                return contractData;
              }

              return { ...input };
            });
          }
        });
      });
      const previousEntries =
        queryClient.getQueryData<GetContractsQuery>(queryKey);

      return { previousEntries };
    },
    onError: (error, { input }, context) => {
      queryClient.setQueryData<GetContractsQuery>(
        queryKey,
        context?.previousEntries,
      );

      toastError(
        'Failed to update billing details',
        `update-contract-error-${error}`,
      );
    },
    onSuccess: () => {
      onClose();
    },
    onSettled: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
      }, 800);
    },
  });

  const defaultValues = new BillingDetailsDto({
    ...data,
    organizationLegalName: data.organizationLegalName || organizationName,
  } as BillingDetailsForm);

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    if (isOpen) {
      const newDefaults = new BillingDetailsDto({
        ...data,
        organizationLegalName: data.organizationLegalName || organizationName,
      } as BillingDetailsForm);

      setDefaultValues(newDefaults);
    }
  }, [isOpen]);

  const handleApplyChanges = () => {
    const payload = BillingDetailsDto.toPayload(state.values);

    updateContract.mutate({
      input: {
        contractId,
        ...payload,
      },
    });
  };

  const fromDummyData = useMemo(
    () => ({
      addressLine: '29 Maple Lane',
      addressLine2: 'Springfield, Haven County',
      locality: 'San Francisco',
      zip: '89302',
      country: 'United States',
      email: 'invoices@acme.com',
      name: 'Acme Corp.',
    }),
    [],
  );

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='4xl'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <Grid h='100%' templateColumns='356px 1fr'>
          <GridItem
            rowSpan={1}
            colSpan={1}
            h='100%'
            display='flex'
            flexDir='column'
            justifyContent='space-between'
            bg='gray.25'
            borderRight='1px solid'
            borderColor='gray.200'
            borderTopLeftRadius='2xl'
            borderBottomLeftRadius='2xl'
            backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
            backgroundRepeat='no-repeat'
            sx={{
              backgroundPositionX: '1px',
              backgroundPositionY: '-7px',
            }}
          >
            <ModalHeader>
              <FeaturedIcon size='lg' colorScheme='primary'>
                <File02 color='primary.600' />
              </FeaturedIcon>
              <Heading fontSize='lg' mt='4'>
                {organizationName} contract details
              </Heading>
            </ModalHeader>
            <ModalBody pb='0' gap={4} display='flex' flexDir='column' flex={1}>
              <FormInput
                label='Link to contract'
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
                }}
                formId={formId}
                name='contractUrl'
                textOverflow='ellipsis'
                placeholder='Paste or enter a contract link'
              />

              <FormInput
                label='Organization legal name'
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
                }}
                onMouseEnter={() => setIsBillingDetailsHovered(true)}
                onMouseLeave={() => setIsBillingDetailsHovered(false)}
                onFocus={() => setIsBillingDetailsFocused(true)}
                onBlur={() => setIsBillingDetailsFocused(false)}
                formId={formId}
                name='organizationLegalName'
                textOverflow='ellipsis'
                placeholder='Organization legal name'
              />

              <Flex
                flexDir='column'
                onMouseEnter={() => setIsBillingDetailsHovered(true)}
                onMouseLeave={() => setIsBillingDetailsHovered(false)}
              >
                <FormInput
                  label='Billing address'
                  isLabelVisible
                  labelProps={{
                    fontSize: 'sm',
                    mb: 0,
                    fontWeight: 'semibold',
                  }}
                  formId={formId}
                  name='addressLine1'
                  textOverflow='ellipsis'
                  placeholder='Address line 1'
                  onFocus={() => setIsBillingDetailsFocused(true)}
                  onBlur={() => setIsBillingDetailsFocused(false)}
                />
                <FormInput
                  label='Address line 2'
                  formId={formId}
                  name='addressLine2'
                  textOverflow='ellipsis'
                  placeholder='Address line 2'
                  onFocus={() => setIsBillingDetailsFocused(true)}
                  onBlur={() => setIsBillingDetailsFocused(false)}
                />
                <Flex>
                  <FormInput
                    label='City'
                    formId={formId}
                    name='locality'
                    textOverflow='ellipsis'
                    placeholder='City'
                    onFocus={() => setIsBillingDetailsFocused(true)}
                    onBlur={() => setIsBillingDetailsFocused(false)}
                  />
                  <FormInput
                    label='ZIP/Postal code'
                    formId={formId}
                    name='zip'
                    textOverflow='ellipsis'
                    placeholder='ZIP/Potal code'
                    onFocus={() => setIsBillingDetailsFocused(true)}
                    onBlur={() => setIsBillingDetailsFocused(false)}
                  />
                </Flex>
                <FormSelect
                  label='Country'
                  placeholder='Country'
                  name='country'
                  formId={formId}
                  options={countryOptions}
                  onFocus={() => setIsBillingDetailsFocused(true)}
                  onBlur={() => setIsBillingDetailsFocused(false)}
                />
              </Flex>
              <FormInput
                label='Send invoice to'
                isLabelVisible
                labelProps={{
                  fontSize: 'sm',
                  mb: 0,
                  fontWeight: 'semibold',
                }}
                formId={formId}
                name='invoiceEmail'
                textOverflow='ellipsis'
                placeholder='Email'
                onMouseEnter={() => setIsBillingDetailsHovered(true)}
                onMouseLeave={() => setIsBillingDetailsHovered(false)}
                onFocus={() => setIsBillingDetailsFocused(true)}
                onBlur={() => setIsBillingDetailsFocused(false)}
              />
              <FormSelect
                label='Currency'
                placeholder='currency'
                isLabelVisible
                name='currency'
                formId={formId}
                options={currencyOptions}
              />
            </ModalBody>
            <ModalFooter p='6'>
              <Button variant='outline' w='full' onClick={onClose}>
                Cancel
              </Button>
              <Button
                ml='3'
                w='full'
                variant='outline'
                colorScheme='primary'
                loadingText='Applying changes...'
                onClick={handleApplyChanges}
              >
                Done
              </Button>
            </ModalFooter>
          </GridItem>
          <GridItem>
            <Box width='100%'>
              <Invoice
                isBilledToFocused={
                  isBillingDetailsFocused || isBillingDetailsHovered
                }
                tax={0}
                note={''}
                from={fromDummyData}
                total={0}
                dueDate={new Date().toISOString()}
                subtotal={0}
                issueDate={new Date().toISOString()}
                billedTo={{
                  addressLine: state.values.addressLine1 ?? '',
                  addressLine2: state.values.addressLine2 ?? '',
                  locality: state.values.locality ?? '',
                  zip: state.values.zip ?? '',
                  country: state?.values?.country?.label ?? '',
                  email: state.values.invoiceEmail ?? '',
                  name: state.values?.organizationLegalName ?? '',
                }}
                status='Preview'
                invoiceNumber='INV-003'
                lines={[]}
              />
            </Box>
          </GridItem>
        </Grid>
      </ModalContent>
    </Modal>
  );
};
