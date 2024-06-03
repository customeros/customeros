import React, { useRef } from 'react';
import { useParams } from 'react-router-dom';
import { useForm } from 'react-inverted-form';

import { useQueryClient } from '@tanstack/react-query';
import { UseMutationResult } from '@tanstack/react-query';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { XSquare } from '@ui/media/icons/XSquare';
import { ModalBody } from '@ui/overlay/Modal/Modal';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { Exact, ContractUpdateInput } from '@graphql/types';
import { DatePickerUnderline } from '@ui/form/DatePicker/DatePickerUnderline';
import { UpdateContractMutation } from '@organization/graphql/updateContract.generated';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/graphql/getContracts.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal/Modal';
import {
  ContractStatusModalMode,
  useContractModalStatusContext,
} from '@organization/components/Tabs/panels/AccountPanel/context/ContractStatusModalsContext';

interface ContractEndModalProps {
  renewsAt?: string;
  contractId: string;
  contractEnded?: string;
  serviceStarted?: string;
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
  contractId,
  organizationName,
  renewsAt,
  onUpdateContract,
  contractEnded,
}: ContractEndModalProps) => {
  const queryClient = useQueryClient();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const id = useParams()?.id as string;
  const queryKey = useGetContractsQuery.getKey({ id });

  const [value, setValue] = React.useState(EndContract.Now);
  const formId = `contract-ends-on-form-${contractId}`;
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;
  const { isModalOpen, onStatusModalClose, mode, nextInvoice } =
    useContractModalStatusContext();
  const timeToNextInvoice = nextInvoice?.issued
    ? DateTimeUtils.format(
        nextInvoice.issued,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  const { state, setDefaultValues } = useForm<{
    endedAt?: string | Date | null;
  }>({
    formId,
    defaultValues: { endedAt: contractEnded || new Date() },
    stateReducer: (_, _action, next) => {
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
          onStatusModalClose();
        },
        onSettled: () => {
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
          }

          timeoutRef.current = setTimeout(() => {
            queryClient.invalidateQueries({ queryKey });
            queryClient.invalidateQueries({
              queryKey: ['GetTimeline.infinite'],
            });
          }, 1000);
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
      setDefaultValues({ endedAt: nextInvoice?.issued });
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
      open={isModalOpen && mode === ContractStatusModalMode.End}
      onOpenChange={onStatusModalClose}
    >
      <ModalOverlay className='z-50' />
      <ModalContent className='rounded-2xl z-[999]'>
        <ModalHeader className='pb-3'>
          <FeaturedIcon
            size='lg'
            colorScheme='error'
            className='mt-3 mb-6 ml-[10px] '
          >
            <XSquare className='text-error-600' />
          </FeaturedIcon>
          <h2 className='text-lg mt-2 font-semibold'>
            End
            {organizationName}’s contract?
          </h2>
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
            className='flex flex-col gap-1 text-base'
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
                    <DatePickerUnderline formId={formId} name='endedAt' />
                  </div>
                ) : (
                  'custom date'
                )}
              </div>
            </Radio>
          </RadioGroup>
        </ModalBody>
        <ModalFooter className='p-6 flex'>
          <Button
            variant='outline'
            colorScheme='gray'
            className='w-full'
            onClick={onStatusModalClose}
            size='lg'
          >
            Cancel
          </Button>
          <Button
            className='ml-3 w-full'
            size='lg'
            variant='outline'
            colorScheme='error'
            onClick={handleApplyChanges}
            loadingText='Saving...'
            isLoading={onUpdateContract.isPending}
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
