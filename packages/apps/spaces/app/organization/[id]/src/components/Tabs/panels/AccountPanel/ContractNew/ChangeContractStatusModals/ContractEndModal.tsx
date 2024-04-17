'use client';
import React, { useRef } from 'react';
import { useForm } from 'react-inverted-form';

import { UseMutationResult } from '@tanstack/react-query';

import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { XSquare } from '@ui/media/icons/XSquare';
import { DateTimeUtils } from '@spaces/utils/date';
import { ModalBody } from '@ui/overlay/Modal/Modal';
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
  EndOfCurrentRenewalPeriod = 'EndOfCurrentRenewalPeriod',
  CustomDate = 'CustomDate',
}
export const ContractEndModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
  renewsAt,
  nextInvoiceDate,
  onUpdateContract,
  contractEnded,
}: ContractEndModalProps) => {
  const initialRef = useRef(null);
  const [value, setValue] = React.useState(EndContract.Now);
  const formId = `contract-ends-on-form-${contractId}`;
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;

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
    if (nextValue === EndContract.Now) {
      setDefaultValues({ endedAt: today });
      setValue(EndContract.Now);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentBillingPeriod) {
      setDefaultValues({ endedAt: nextInvoiceDate });
      setValue(EndContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === EndContract.CustomDate) {
      setDefaultValues({ endedAt: new Date(today) });
      setValue(EndContract.CustomDate);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentRenewalPeriod) {
      setDefaultValues({ endedAt: renewsAt });
      setValue(EndContract.EndOfCurrentRenewalPeriod);

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
          <h2 className='text-lg mt-4'>End {organizationName}’s contract?</h2>
        </ModalHeader>
        <ModalBody className='flex flex-col gap-3'>
          <p className='text-base'>
            Ending this contract{' '}
            <span className='font-medium mr-1'>will close the renewal</span>
            and set the
            <span className='font-medium ml-1'>ARR to zero.</span>
          </p>
          <p className='text-base'>Let’s end it:</p>

          <RadioGroup
            value={value}
            onValueChange={handleChangeEndsOnOption}
            className='flex flex-col gap-1'
          >
            <Radio value={EndContract.Now}>
              <span className='mr-1'>Now</span>
            </Radio>
            {timeToNextInvoice && (
              <Radio value={EndContract.EndOfCurrentBillingPeriod}>
                <span className='ml-1'>
                  End of current billing period, {timeToNextInvoice}
                </span>
              </Radio>
            )}

            {timeToRenewal && (
              <Radio value={EndContract.EndOfCurrentRenewalPeriod}>
                <span className='ml-1'>End of renewal, {timeToRenewal}</span>
              </Radio>
            )}

            <Radio value={EndContract.CustomDate}>
              <div className='flex items-center max-h-6'>
                On{' '}
                {value === EndContract.CustomDate ? (
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
