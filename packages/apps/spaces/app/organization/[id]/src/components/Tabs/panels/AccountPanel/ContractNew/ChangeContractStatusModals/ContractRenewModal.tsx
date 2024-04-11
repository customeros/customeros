'use client';
import React, { useRef } from 'react';
import { useForm } from 'react-inverted-form';

import { UseMutationResult } from '@tanstack/react-query';

import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { RefreshCw05 } from '@ui/media/icons/RefreshCw05';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio2';
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
  contractEnded?: string;
  serviceStarted?: string;
  nextInvoiceDate?: string;
  organizationName: string;
  committedPeriods?: number;
  onUpdateContract: UseMutationResult<
    UpdateContractMutation,
    unknown,
    Exact<{ input: ContractUpdateInput }>,
    { previousEntries: GetContractsQuery | undefined }
  >;
}

const today = new Date().toUTCString();

export enum StartContract {
  Now = 'Now',
  NextContractTerm = 'NextContractTerm',
  CustomDate = 'CustomDate',
}
export const ContractRenewsModal = ({
  isOpen,
  onClose,
  contractId,
  nextInvoiceDate,
  onUpdateContract,
  committedPeriods,
  contractEnded,
}: ContractEndModalProps) => {
  const initialRef = useRef(null);
  const [value, setValue] = React.useState(StartContract.Now);
  const formId = `contract-renew-on-form-${contractId}`;

  const timeToNextInvoice = nextInvoiceDate
    ? DateTimeUtils.format(
        nextInvoiceDate,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  const { state, setDefaultValues } = useForm<{
    endedAt?: string | Date | null;
  }>({
    formId,
    defaultValues: { endedAt: contractEnded || new Date() },
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
    if (nextValue === StartContract.Now) {
      setDefaultValues({ startedAt: today, endedAt: today });
      setValue(StartContract.Now);

      return;
    }
    if (
      nextValue === StartContract.NextContractTerm &&
      nextInvoiceDate &&
      committedPeriods
    ) {
      setDefaultValues({
        endedAt: nextInvoiceDate,
        startedAt: DateTimeUtils.addMonth(
          nextInvoiceDate,
          committedPeriods + 1,
        ),
      });
      setValue(StartContract.NextContractTerm);

      return;
    }
    if (nextValue === StartContract.CustomDate) {
      setDefaultValues({ endedAt: new Date(today) });
      setValue(StartContract.CustomDate);

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
            <RefreshCw05 className='text-error-600' />
          </FeaturedIcon>
          <h2 className='text-lg mt-4'>Renew contract:</h2>
        </ModalHeader>
        <ModalBody className='flex flex-col gap-3'>
          <RadioGroup
            value={value}
            onValueChange={handleChangeEndsOnOption}
            className='flex flex-col gap-1'
          >
            <Radio value={StartContract.Now}>
              <span className='mr-1'>Now</span>
            </Radio>
            {timeToNextInvoice && (
              <Radio value={StartContract.NextContractTerm}>
                <span className='ml-1'>
                  End of current billing period, {timeToNextInvoice}
                </span>
              </Radio>
            )}

            <Radio value={StartContract.CustomDate}>
              <div className='flex items-center max-h-6'>
                On{' '}
                {value === StartContract.CustomDate ? (
                  <div className='ml-1'>
                    <DatePickerUnderline
                      placeholder='End date'
                      defaultOpen={true}
                      // minDate={state.values.serviceStarted}
                      formId={formId}
                      name='endedAt'
                      calendarIconHidden
                      value={state.values.endedAt}
                    />
                  </div>
                ) : (
                  'custom date'
                )}
              </div>
            </Radio>
          </RadioGroup>
        </ModalBody>
        <ModalFooter p='6'>
          <Button
            variant='outline'
            colorScheme='gray'
            className='w-full'
            onClick={onClose}
          >
            Cancel
          </Button>
          <Button
            className='ml-3 w-full'
            variant='outline'
            colorScheme='error'
            onClick={handleApplyChanges}
          >
            Renew {value === StartContract.Now && 'now'}
            {value !== StartContract.Now &&
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
