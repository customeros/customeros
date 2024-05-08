'use client';
import { useMemo, useCallback } from 'react';
import { useForm, useField } from 'react-inverted-form';

import { match } from 'ts-pattern';
import { twMerge } from 'tailwind-merge';
import { useDeepCompareEffect } from 'rooks';
import { UseMutationResult } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Dot } from '@ui/media/Dot';
import { Button } from '@ui/form/Button/Button';
import { DateTimeUtils } from '@spaces/utils/date';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateOpportunityRenewalMutation } from '@organization/src/graphql/updateOpportunityRenewal.generated';
import { likelihoodButtons } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/RenewalARR/utils';
import {
  RangeSlider,
  RangeSliderTrack,
  RangeSliderThumb,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider';
import {
  Exact,
  Currency,
  Opportunity,
  InternalStage,
  OpportunityRenewalLikelihood,
  OpportunityRenewalUpdateInput,
} from '@graphql/types';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalPortal,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal/Modal';

type UpdateOpportunityMutation = UseMutationResult<
  UpdateOpportunityRenewalMutation,
  unknown,
  Exact<{ input: OpportunityRenewalUpdateInput }>,
  { previousEntries: GetContractsQuery | undefined }
>;
interface RenewalDetailsProps {
  isOpen: boolean;
  data: Opportunity;
  onClose: () => void;
  currency?: string | null;
  updateOpportunityMutation: UpdateOpportunityMutation;
}

export const RenewalDetailsModal = ({
  data,
  isOpen,
  onClose,
  currency = Currency.Usd,
  updateOpportunityMutation,
}: RenewalDetailsProps) => {
  return (
    <>
      {isOpen && (
        <Modal
          open={data?.internalStage !== InternalStage.ClosedLost && isOpen}
          onOpenChange={onClose}
        >
          <ModalPortal>
            <ModalOverlay />
            <RenewalDetailsForm
              data={data}
              onClose={onClose}
              currency={currency as string}
              updateOpportunityMutation={updateOpportunityMutation}
            />
          </ModalPortal>
        </Modal>
      )}
    </>
  );
};

interface RenewalDetailsFormProps {
  currency: string;
  data: Opportunity;
  onClose?: () => void;
  updateOpportunityMutation: UpdateOpportunityMutation;
}

const RenewalDetailsForm = ({
  data,
  onClose,
  currency,
  updateOpportunityMutation,
}: RenewalDetailsFormProps) => {
  const client = getGraphQLClient();
  const formId = `renewal-details-form-${data.id}`;

  const { data: usersData } = useGetUsersQuery(client, {
    pagination: {
      limit: 50,
      page: 1,
    },
  });
  const updatedAt = data?.updatedAt
    ? DateTimeUtils.timeAgo(data?.updatedAt)
    : null;
  const maxAmount = data.maxAmount ?? 0;
  const renewadAt = data?.renewedAt;

  const defaultValues = useMemo(
    () => ({
      renewalAdjustedRate: data?.renewalAdjustedRate ?? 50,
      renewalLikelihood: data?.renewalLikelihood,
      reason: data?.comments,
    }),
    [data?.renewalLikelihood, data?.amount, data?.comments],
  );

  const updatedByUser = usersData?.users.content?.find(
    (u) => u.id === data.renewalUpdatedByUserId,
  );
  const updatedByUserFullName =
    updatedByUser?.name ||
    [updatedByUser?.firstName, updatedByUser?.lastName]
      .filter(Boolean)
      .join(' ');

  const onSubmit = useCallback(
    async (state: typeof defaultValues) => {
      const { reason, renewalLikelihood, renewalAdjustedRate } = state;

      updateOpportunityMutation.mutate({
        input: {
          opportunityId: data.id,
          comments: reason,
          renewalLikelihood,
          renewalAdjustedRate,
        },
      });
    },
    [updateOpportunityMutation],
  );

  const { handleSubmit, setDefaultValues } = useForm({
    formId,
    defaultValues,
    onSubmit,
    stateReducer: (_state, action, next) => {
      if (
        action.type === 'FIELD_CHANGE' &&
        action.payload.name === 'renewalLikelihood'
      ) {
        const nextRate = match(action.payload.value)
          .with(OpportunityRenewalLikelihood.LowRenewal, () => 25)
          .with(OpportunityRenewalLikelihood.MediumRenewal, () => 50)
          .with(OpportunityRenewalLikelihood.HighRenewal, () => 100)
          .otherwise(() => 100);

        return {
          ...next,
          values: {
            ...next.values,
            renewalAdjustedRate: nextRate,
          },
        };
      }

      return next;
    },
  });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

  return (
    <>
      <ModalContent
        className='rounded-2xl bg-[url(/backgrounds/organization/circular-bg-pattern.png)] bg-no-repeat'
        style={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon
            size='lg'
            colorScheme='primary'
            className='ml-[12px] mt-1 mb-[31px]'
          >
            <ClockFastForward />
          </FeaturedIcon>
          <span className='text-lg mt-3'>Renewal details</span>
        </ModalHeader>
        <form onSubmit={(v) => handleSubmit(v)}>
          <ModalBody className='pb-0 gap-4 flex flex-col'>
            <div>
              <FormLikelihoodButtonGroup
                formId={formId}
                name='renewalLikelihood'
              />

              {updatedAt && (
                <p className='text-gray-500 text-xs mt-2'>
                  Last updated{' '}
                  {updatedByUserFullName
                    ? `by ${updatedByUserFullName}`
                    : 'automatically'}{' '}
                  {updatedAt} ago
                </p>
              )}
            </div>

            <FormRangeSlider
              formId={formId}
              currency={currency}
              name='renewalAdjustedRate'
              amount={maxAmount}
              renewadAt={renewadAt}
            />

            {!!data.renewalLikelihood && (
              <div>
                <label className='text-sm' htmlFor='reason'>
                  <b>Reason for change</b> (optional)
                </label>
                <FormAutoresizeTextarea
                  className='pt-0 text-base'
                  size='sm'
                  formId={formId}
                  id='reason'
                  name='reason'
                  spellCheck='false'
                  placeholder={`What is the reason for updating these details`}
                />
              </div>
            )}
          </ModalBody>

          <ModalFooter className='flex p-6'>
            <Button
              variant='outline'
              className='w-full'
              onClick={onClose}
              isDisabled={updateOpportunityMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              className='ml-3 w-full'
              variant='outline'
              colorScheme='primary'
              isLoading={updateOpportunityMutation.isPending}
              typeof='submit'
              loadingText='Updating...'
              spinner={
                <Spinner
                  label='Updating...'
                  className='text-primary-500 fill-primary-700 size-4'
                />
              }
            >
              Update
            </Button>
          </ModalFooter>
        </form>
      </ModalContent>
    </>
  );
};

interface LikelihoodButtonGroupProps {
  value?: OpportunityRenewalLikelihood | null;
  onBlur?: (value: OpportunityRenewalLikelihood) => void;
  onChange?: (value: OpportunityRenewalLikelihood) => void;
}

const LikelihoodButtonGroup = ({
  value,
  onBlur,
  onChange,
}: LikelihoodButtonGroupProps) => {
  return (
    <div
      className='inline-flex w-full'
      aria-disabled={value === OpportunityRenewalLikelihood.ZeroRenewal}
      aria-describedby='likelihood-oprions-button'
    >
      {likelihoodButtons.map((button, idx) => (
        <Button
          key={`${button.likelihood}-likelihood-button`}
          variant='outline'
          className={twMerge(
            idx === 0
              ? 'border-e-0 rounded-s-lg rounded-e-none !important'
              : idx === 1
              ? 'rounded-none'
              : 'border-s-0 rounded-s-none rounded-e-lg !important',
            'w-full data-[selected=true]:bg-gray-50 !important',
          )}
          onBlur={() => onBlur?.(button.likelihood)}
          onClick={(e) => {
            e.preventDefault();
            onChange?.(button.likelihood);
          }}
          data-selected={value === button.likelihood}
        >
          <div className='flex items-center gap-1'>
            <Dot colorScheme={button.colorScheme} className='size-2 mr-2' />
            {button.label}
          </div>
        </Button>
      ))}
    </div>
  );
};

interface FormLikelihoodButtonGroupProps {
  name: string;
  formId: string;
}

const FormLikelihoodButtonGroup = ({
  name,
  formId,
}: FormLikelihoodButtonGroupProps) => {
  const { getInputProps } = useField(name, formId);
  const { value, onChange, onBlur } = getInputProps();

  return (
    <LikelihoodButtonGroup value={value} onChange={onChange} onBlur={onBlur} />
  );
};

interface FormRangeSliderProps {
  name: string;
  formId: string;
  amount?: number;
  currency: string;
  renewadAt?: string;
}

const FormRangeSlider = ({
  name,
  formId,
  amount = 0,
  currency,
  renewadAt,
}: FormRangeSliderProps) => {
  const { getInputProps } = useField(name, formId);
  const { value, onChange, onBlur, ...rest } = getInputProps();
  const defaultFormattedAmount = formatCurrency(amount, 2, currency);
  const formattedNewAmount = formatCurrency(
    amount * (value / 100),
    2,
    currency,
  );
  const formattedRenewedAt = renewadAt
    ? DateTimeUtils.format(renewadAt, DateTimeUtils.dateWithAbreviatedMonth)
    : undefined;

  const trackStyle = cn('h-0.5 transition-colors', {
    'bg-orangeDark-700': value <= 25,
    'bg-yellow-400': value > 25 && value < 75,
    'bg-greenLight-400': value >= 75,
  });

  const thumbStyle = cn('ring-1 transition-colors shadow-md cursor-pointer', {
    'ring-orangeDark-700': value <= 25,
    'ring-yellow-400': value > 25 && value < 75,
    'ring-greenLight-400': value >= 75,
  });

  return (
    <div>
      <div className='flex items-center justify-between mb-3'>
        <p className='font-medium text-base'>
          Renewal ARR{' '}
          {formattedRenewedAt && (
            <span className='text-gray-400 font-normal text-sm'>
              on {formattedRenewedAt}
            </span>
          )}
        </p>

        <p className='text-base font-medium'>
          {formattedNewAmount !== defaultFormattedAmount && (
            <span className='text-sm text-gray-400 font-normal'>
              <s>{defaultFormattedAmount}</s>
            </span>
          )}{' '}
          {formattedNewAmount}
        </p>
      </div>
      <RangeSlider
        step={1}
        min={0}
        max={100}
        value={[value]}
        className='w-full'
        onValueChange={(values) => {
          onChange(values[0]);
        }}
        onValueCommit={(values) => {
          onBlur(values[0]);
        }}
        {...rest}
      >
        <RangeSliderTrack className='bg-gray-400 h-0.5'>
          <RangeSliderFilledTrack className={trackStyle} />
        </RangeSliderTrack>
        <RangeSliderThumb className={thumbStyle} />
      </RangeSlider>
    </div>
  );
};
