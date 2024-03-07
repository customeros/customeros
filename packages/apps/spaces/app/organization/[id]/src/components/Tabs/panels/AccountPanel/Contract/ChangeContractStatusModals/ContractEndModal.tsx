'use client';
import React, { useRef } from 'react';
import { useForm } from 'react-inverted-form';

import { UseMutationResult } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { ModalBody } from '@ui/overlay/Modal';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { XSquare } from '@ui/media/icons/XSquare';
import { DateTimeUtils } from '@spaces/utils/date';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { Exact, ContractUpdateInput } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateContractMutation } from '@organization/src/graphql/updateContract.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

interface ContractEndModalProps {
  isOpen: boolean;
  renewsAt?: string;
  contractId: string;
  onClose: () => void;
  organizationName: string;
  serviceStartedAt?: string;
  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

const today = new Date().toUTCString();

export enum EndContract {
  Now = 'Now',
  EndOfCurrentBillingPeriod = 'EndOfCurrentBillingPeriod',
  CustomDate = 'CustomDate',
}
export const ContractEndModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  renewsAt,
  onUpdateContract,
}: ContractEndModalProps) => {
  const initialRef = useRef(null);
  const [value, setValue] = React.useState(EndContract.Now);
  const formId = `contract-ends-on-form-${contractId}`;

  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

  const { state, setDefaultValues } = useForm<{
    endedAt?: string | Date | null;
  }>({
    formId,
    defaultValues: { endedAt: new Date() },
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  const handleApplyChanges = () => {
    onUpdateContract.mutate(
      {
        input: {
          contractId,
          patch: true,
          endedAt: state.values.endedAt,
        },
      },
      {
        onSuccess: () => {
          onClose();
        },
      },
    );
  };

  const handleChangeEndsOnOption = (nextValue: string | null) => {
    if (nextValue === EndContract.Now) {
      setDefaultValues({ endedAt: today });
      setValue(EndContract.Now);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentBillingPeriod) {
      setDefaultValues({ endedAt: renewsAt });
      setValue(EndContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === EndContract.CustomDate) {
      setDefaultValues({ endedAt: new Date(today) });
      setValue(EndContract.CustomDate);

      return;
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='md'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='error'>
            <XSquare color='error.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            End {organizationName}’s contract?
          </Heading>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' gap={4}>
          <Text>
            Ending this contract{' '}
            <Text fontWeight='medium' as='span' mr={1}>
              will close the renewal
            </Text>
            and set the
            <Text fontWeight='medium' as='span' ml={1}>
              ARR to zero.
            </Text>
          </Text>
          <Text>Let’s end it on:</Text>

          <RadioGroup
            value={value}
            onChange={handleChangeEndsOnOption}
            flexDir='column'
            display='flex'
          >
            <Radio value={EndContract.Now} colorScheme='primary'>
              Now
            </Radio>
            <Radio
              value={EndContract.EndOfCurrentBillingPeriod}
              colorScheme='primary'
              display={renewsAt ? 'flex' : 'none'}
            >
              End of current billing period, {timeToRenewal}
            </Radio>
            <Radio value={EndContract.CustomDate} colorScheme='primary'>
              <Flex alignItems='center'>
                On a{' '}
                {value === EndContract.CustomDate ? (
                  <Box ml={1}>
                    <DatePickerUnderline
                      placeholder='End date'
                      // minDate={state.values.serviceStartedAt}
                      formId={formId}
                      name='endedAt'
                      calendarIconHidden
                      value={state.values.endedAt}
                    />
                  </Box>
                ) : (
                  'custom date'
                )}
              </Flex>
            </Radio>
          </RadioGroup>
        </ModalBody>
        <ModalFooter p='6'>
          <Button variant='outline' w='full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            variant='outline'
            colorScheme='error'
            loadingText='Applying changes...'
            onClick={handleApplyChanges}
          >
            End {value === EndContract.Now && 'now'}
            {value !== EndContract.Now &&
              DateTimeUtils.format(
                state?.values?.endedAt as string,
                DateTimeUtils.dateWithAbreviatedMonth,
              )}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
