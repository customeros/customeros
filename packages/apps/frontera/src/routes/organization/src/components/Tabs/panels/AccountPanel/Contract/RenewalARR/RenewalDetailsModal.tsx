import { useCallback } from 'react';
import { useForm, useField } from 'react-inverted-form';

import { match } from 'ts-pattern';
import { twMerge } from 'tailwind-merge';
import { useDeepCompareEffect } from 'rooks';

import { cn } from '@ui/utils/cn';
import { Dot } from '@ui/media/Dot';
import { DateTimeUtils } from '@utils/date';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import {
  Currency,
  Opportunity,
  InternalStage,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import { likelihoodButtons } from '@organization/components/Tabs/panels/AccountPanel/Contract/RenewalARR/utils';
import {
  RangeSlider,
  RangeSliderThumb,
  RangeSliderTrack,
  RangeSliderFilledTrack,
} from '@ui/form/RangeSlider/RangeSlider';
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

interface RenewalDetailsProps {
  isOpen: boolean;
  data: Opportunity;
  onClose: () => void;
  currency?: string | null;
  updateOpportunityMutation: (data: Partial<Opportunity>) => void;
}

export const RenewalDetailsModal = ({
  data,
  isOpen,
  onClose,
  currency,
  updateOpportunityMutation,
}: RenewalDetailsProps) => {
  return (
    <>
      {isOpen && (
        <Modal
          onOpenChange={onClose}
          open={data?.internalStage !== InternalStage.ClosedLost && isOpen}
        >
          <ModalPortal>
            <ModalOverlay className='z-50' />
            <RenewalDetailsForm
              data={data}
              onClose={onClose}
              currency={currency || Currency.Usd}
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
  updateOpportunityMutation: (data: Partial<Opportunity>) => void;
}

const RenewalDetailsForm = ({
  data,
  onClose,
  currency,
  updateOpportunityMutation,
}: RenewalDetailsFormProps) => {
  const store = useStore();
  const users = store.users.toArray();
  const formId = `renewal-details-form-${data.id}`;
  const updatedAt = data?.updatedAt
    ? DateTimeUtils.timeAgo(data?.updatedAt)
    : null;
  const maxAmount = data.maxAmount ?? 0;
  const renewadAt = data?.renewedAt;

  const getAdjustedRate = (value: OpportunityRenewalLikelihood) => {
    return match(value)
      .with(OpportunityRenewalLikelihood.LowRenewal, () => 25)
      .with(OpportunityRenewalLikelihood.MediumRenewal, () => 50)
      .with(OpportunityRenewalLikelihood.HighRenewal, () => 100)
      .otherwise(() => 100);
  };

  const defaultValues = {
    renewalAdjustedRate: data?.renewalAdjustedRate
      ? data?.renewalAdjustedRate
      : data?.renewalLikelihood
      ? getAdjustedRate(data?.renewalLikelihood)
      : 100,
    renewalLikelihood: data?.renewalLikelihood,
    reason: data?.comments,
  };

  const updatedByUser = users?.find(
    (u) => u.id === data.renewalUpdatedByUserId,
  );
  const updatedByUserFullName = updatedByUser?.name;

  const onSubmit = useCallback(
    async (state: typeof defaultValues) => {
      const { reason, renewalLikelihood, renewalAdjustedRate } = state;

      updateOpportunityMutation({
        comments: reason,
        renewalLikelihood,
        renewalAdjustedRate,
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
        const nextRate = getAdjustedRate(action.payload.value);

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
        className='z-50 rounded-2xl '
        style={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <span className='text-lg mt-3 font-semibold'>Renewal details</span>
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
                  {updatedAt === 'today' ? `${updatedAt}` : `${updatedAt} ago`}
                </p>
              )}
            </div>

            <FormRangeSlider
              formId={formId}
              amount={maxAmount}
              currency={currency}
              renewadAt={renewadAt}
              name='renewalAdjustedRate'
            />

            {!!data.renewalLikelihood && (
              <div>
                <label htmlFor='reason' className='text-sm'>
                  <b>Reason for change</b> (optional)
                </label>
                <FormAutoresizeTextarea
                  size='sm'
                  id='reason'
                  name='reason'
                  formId={formId}
                  spellCheck='false'
                  className='pt-0 text-base'
                  placeholder={`What is the reason for updating these details`}
                />
              </div>
            )}
          </ModalBody>

          <ModalFooter className='flex p-6'>
            <Button variant='outline' onClick={onClose} className='w-full'>
              Cancel
            </Button>
            <Button
              typeof='submit'
              variant='outline'
              colorScheme='primary'
              className='ml-3 w-full'
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
      aria-describedby='likelihood-oprions-button'
      aria-disabled={value === OpportunityRenewalLikelihood.ZeroRenewal}
    >
      {likelihoodButtons.map((button, idx) => (
        <Button
          variant='outline'
          onBlur={() => onBlur?.(button.likelihood)}
          data-selected={value === button.likelihood}
          key={`${button.likelihood}-likelihood-button`}
          onClick={(e) => {
            e.preventDefault();
            onChange?.(button.likelihood);
          }}
          className={twMerge(
            idx === 0
              ? 'border-e-0 rounded-s-lg rounded-e-none !important'
              : idx === 1
              ? 'rounded-none'
              : 'border-s-0 rounded-s-none rounded-e-lg !important ',
            'w-full data-[selected=true]:bg-white !important bg-gray-50',
          )}
        >
          <div className='flex items-center gap-1'>
            <Dot className='size-2 mr-2' colorScheme={button.colorScheme} />
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
    <LikelihoodButtonGroup value={value} onBlur={onBlur} onChange={onChange} />
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
  currency = 'USD',
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
        min={0}
        step={1}
        max={100}
        value={[value]}
        className='w-full'
        onValueCommit={(values) => {
          onBlur(values[0]);
        }}
        onValueChange={(values) => {
          onChange(values[0]);
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
