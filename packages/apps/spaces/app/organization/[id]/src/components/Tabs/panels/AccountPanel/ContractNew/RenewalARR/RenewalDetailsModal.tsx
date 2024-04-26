'use client';
import { useMemo, useCallback } from 'react';
import { useForm, useField } from 'react-inverted-form';

import { twMerge } from 'tailwind-merge';
import { useDeepCompareEffect } from 'rooks';
import { UseMutationResult } from '@tanstack/react-query';

import { Dot } from '@ui/media/Dot';
import { Button } from '@ui/form/Button/Button';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { FormCurrencyInput } from '@ui/form/CurrencyInput';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { useGetUsersQuery } from '@shared/graphql/getUsers.generated';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { GetContractsQuery } from '@organization/src/graphql/getContracts.generated';
import { UpdateOpportunityRenewalMutation } from '@organization/src/graphql/updateOpportunityRenewal.generated';
import { likelihoodButtons } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/RenewalARR/utils';
import {
  Exact,
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
  updateOpportunityMutation: UpdateOpportunityMutation;
}

export const RenewalDetailsModal = ({
  data,
  isOpen,
  onClose,
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
              updateOpportunityMutation={updateOpportunityMutation}
            />
          </ModalPortal>
        </Modal>
      )}
    </>
  );
};

interface RenewalDetailsFormProps {
  data: Opportunity;
  onClose?: () => void;
  updateOpportunityMutation: UpdateOpportunityMutation;
}

const RenewalDetailsForm = ({
  data,
  onClose,
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

  const options = useMemo(() => {
    return usersData?.users?.content
      ?.filter((e) => Boolean(e.firstName) || Boolean(e.lastName))
      ?.map((o) => ({
        value: o.id,
        label: `${o.firstName} ${o.lastName}`.trim(),
      }));
  }, [usersData?.users?.content?.length]);

  const defaultValues = useMemo(
    () => ({
      renewalLikelihood: data?.renewalLikelihood,
      amount: data?.amount?.toString(),
      reason: data?.comments,
      owner: options?.find((o) => o.value === data?.owner?.id),
    }),
    [data?.renewalLikelihood, data?.amount, data?.comments, data?.owner?.id],
  );

  const onSubmit = useCallback(
    async (state: typeof defaultValues) => {
      const { owner, amount, reason, renewalLikelihood } = state;

      updateOpportunityMutation.mutate({
        input: {
          opportunityId: data.id,
          comments: reason,
          renewalLikelihood,
          ownerUserId: owner?.value,
          amount: parseFloat(amount),
        },
      });
    },
    [updateOpportunityMutation],
  );

  const { handleSubmit, setDefaultValues } = useForm({
    formId,
    defaultValues,
    onSubmit,
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
            <FormSelect
              isClearable
              name='owner'
              label='Owner'
              isLabelVisible
              formId={formId}
              isLoading={false}
              options={options}
              placeholder='Owner'
              backspaceRemovesValue
            />

            <div>
              <FormLikelihoodButtonGroup
                formId={formId}
                name='renewalLikelihood'
              />
              {data?.renewalUpdatedByUserId && (
                <p className='text-gray-500 text-xs mt-2'>Last updated by </p>
              )}
            </div>
            {data?.amount > 0 && (
              <FormCurrencyInput
                className='w-full'
                min={0}
                name='amount'
                formId={formId}
                placeholder='Amount'
                label='ARR forecast'
                leftElement={
                  <CurrencyDollar className='text-gray-500 size-4' />
                }
              />
            )}

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
              ? ' border-e-0 rounded-s-lg rounded-e-none !important'
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
            <Dot colorScheme={button.colorScheme} />
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
