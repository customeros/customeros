import React, { useState } from 'react';

import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { ContractStatus } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { XSquare } from '@ui/media/icons/XSquare';
import { Radio, RadioGroup } from '@ui/form/Radio/Radio';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { DatePickerUnderline2 } from '@ui/form/DatePicker/DatePickerUnderline2.tsx';
import {
  Modal,
  ModalBody,
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
  contractId: string;
  contractEnded?: string;
  serviceStarted?: string;
  organizationName: string;
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
  contractEnded,
}: ContractEndModalProps) => {
  const store = useStore();
  const contractStore = store.contracts.value.get(contractId);

  const [value, setValue] = useState(EndContract.Now);
  const renewsAt = contractStore?.value?.opportunities?.find(
    (e) => e.internalStage === 'OPEN',
  )?.renewedAt;
  const timeToRenewal = renewsAt
    ? DateTimeUtils.format(renewsAt, DateTimeUtils.dateWithAbreviatedMonth)
    : null;
  const nextInvoice = contractStore?.value?.upcomingInvoices?.[0];

  const { isModalOpen, onStatusModalClose, mode } =
    useContractModalStatusContext();
  const timeToNextInvoice = nextInvoice?.issued
    ? DateTimeUtils.format(
        nextInvoice.issued,
        DateTimeUtils.dateWithAbreviatedMonth,
      )
    : null;

  const [endedAt, setEndedAt] = useState<string | Date | null | undefined>(
    contractEnded || new Date().toString(),
  );

  const handleApplyChanges = () => {
    contractStore?.update((prev) => ({
      ...prev,
      endedAt: new Date(endedAt as string),
      contractStatus: DateTimeUtils.isFuture(endedAt as string)
        ? prev.contractStatus
        : ContractStatus.Ended,
    }));
    onStatusModalClose();
  };

  const handleChangeEndsOnOption = (nextValue: string | null) => {
    if (nextValue === EndContract.Now) {
      setEndedAt(today);
      setValue(EndContract.Now);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentBillingPeriod) {
      setEndedAt(nextInvoice?.issued);
      setValue(EndContract.EndOfCurrentBillingPeriod);

      return;
    }
    if (nextValue === EndContract.CustomDate) {
      setEndedAt(new Date(today));

      setValue(EndContract.CustomDate);

      return;
    }
    if (nextValue === EndContract.EndOfCurrentRenewalPeriod) {
      setEndedAt(renewsAt);
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
                    <DatePickerUnderline2
                      value={endedAt || new Date().toString()}
                      onChange={(e) => setEndedAt(e)}
                    />
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
          >
            End {value === EndContract.Now && 'now'}
            {value !== EndContract.Now &&
              DateTimeUtils.format(
                endedAt as string,
                DateTimeUtils.dateWithAbreviatedMonth,
              )}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
