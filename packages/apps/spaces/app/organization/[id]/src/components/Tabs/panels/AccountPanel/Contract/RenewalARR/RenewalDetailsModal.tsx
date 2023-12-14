'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import { useRef, useMemo, useEffect } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { Dot } from '@ui/media/Dot';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { toastError } from '@ui/presentation/Toast';
import { Button, ButtonGroup } from '@ui/form/Button';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FormCurrencyInput } from '@ui/form/CurrencyInput';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { FormElement, FormElementProps } from '@ui/form/FormElement';
import { useGetUsersQuery } from '@organizations/graphql/getUsers.generated';
import {
  Opportunity,
  InternalStage,
  OpportunityRenewalLikelihood,
} from '@graphql/types';
import {
  GetContractsQuery,
  useGetContractsQuery,
} from '@organization/src/graphql/getContracts.generated';
import { useUpdateOpportunityRenewalMutation } from '@organization/src/graphql/updateOpportunityRenewal.generated';
import {
  getButtonStyles,
  likelihoodButtons,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalARR/utils';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

interface RenewalDetailsProps {
  isOpen: boolean;
  data: Opportunity;
  onClose: () => void;
}

export const RenewalDetailsModal = ({
  data,
  isOpen,
  onClose,
}: RenewalDetailsProps) => {
  const orgId = useParams()?.id as string;
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const formId = `renewal-details-form-${data.id}`;
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const { data: usersData } = useGetUsersQuery(client, {
    pagination: {
      limit: 50,
      page: 1,
    },
  });

  const getContractsQueryKey = useGetContractsQuery.getKey({
    id: orgId,
  });

  const updateOpportunityMutation = useUpdateOpportunityRenewalMutation(
    client,
    {
      onMutate: ({ input }) => {
        queryClient.cancelQueries(getContractsQueryKey);

        queryClient.setQueryData<GetContractsQuery>(
          getContractsQueryKey,
          (currentCache) => {
            if (!currentCache || !currentCache?.organization) return;

            return produce(currentCache, (draft) => {
              if (draft?.['organization']?.['contracts']) {
                draft['organization']['contracts']?.map(
                  (contractData, index) => {
                    return (contractData.opportunities ?? []).map(
                      (opportunity) => {
                        const { opportunityId, ...rest } = input;
                        if ((opportunity as Opportunity).id === opportunityId) {
                          return {
                            ...opportunity,
                            ...rest,
                            renewalUpdatedByUserAt: new Date().toISOString(),
                          };
                        }

                        return opportunity;
                      },
                    );
                  },
                );
              }
            });
          },
        );
        const previousEntries =
          queryClient.getQueryData<GetContractsQuery>(getContractsQueryKey);

        return { previousEntries };
      },
      onError: (_, __, context) => {
        queryClient.setQueryData<GetContractsQuery>(
          getContractsQueryKey,
          context?.previousEntries,
        );
        toastError(
          'Failed to update renewal details',
          'update-renewal-details-error',
        );
      },
      onSettled: () => {
        onClose();

        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }
        timeoutRef.current = setTimeout(() => {
          queryClient.invalidateQueries(getContractsQueryKey);
        }, 1000);
      },
    },
  );

  const options = useMemo(() => {
    return usersData?.users?.content
      ?.filter((e) => Boolean(e.firstName) || Boolean(e.lastName))
      ?.map((o) => ({
        value: o.id,
        label: `${o.firstName} ${o.lastName}`.trim(),
      }));
  }, [usersData?.users?.content?.length]);

  const defaultValues = {
    renewalLikelihood: data?.renewalLikelihood,
    amount: data?.amount?.toString(),
    reason: data?.comments,
    owner: options?.find((o) => o.value === data?.owner?.id),
  };

  const { state, handleSubmit } = useForm({
    formId,
    defaultValues,
    onSubmit: async ({ amount, owner, reason, renewalLikelihood }) => {
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
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Modal
      isOpen={data?.internalStage !== InternalStage.ClosedLost && isOpen}
      onClose={onClose}
    >
      <ModalOverlay />
      <ModalContent
        as='form'
        borderRadius='2xl'
        onSubmit={handleSubmit}
        backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
        backgroundRepeat='no-repeat'
        sx={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <ClockFastForward />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Renewal details
          </Heading>
        </ModalHeader>
        <ModalBody pb='0' gap={4} as={Flex} flexDir='column'>
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
              <Text color='gray.500' fontSize='xs' mt={2}>
                Last updated by{' '}
              </Text>
            )}
          </div>
          {data?.amount > 0 && (
            <FormCurrencyInput
              min={0}
              w='full'
              name='amount'
              formId={formId}
              placeholder='Amount'
              label='ARR forecast'
              isLabelVisible
              leftElement={
                <Box color='gray.500'>
                  <CurrencyDollar height='16px' />
                </Box>
              }
            />
          )}

          {!!state.values.renewalLikelihood && (
            <div>
              <Text as='label' htmlFor='reason' fontSize='sm'>
                <b>Reason for change</b> (optional)
              </Text>
              <FormAutoresizeTextarea
                pt='0'
                formId={formId}
                id='reason'
                name='reason'
                spellCheck='false'
                placeholder={`What is the reason for updating these details`}
              />
            </div>
          )}
        </ModalBody>
        <ModalFooter p='6'>
          <Button variant='outline' w='full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            type='submit'
            variant='outline'
            colorScheme='primary'
          >
            Update
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
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
    <ButtonGroup
      w='full'
      isAttached
      isDisabled={value === OpportunityRenewalLikelihood.ZeroRenewal}
      aria-describedby='likelihood-oprions-button'
    >
      {likelihoodButtons.map((button) => (
        <Button
          key={`${button.likelihood}-likelihood-button`}
          variant='outline'
          leftIcon={<Dot colorScheme={button.colorScheme} />}
          onBlur={() => onBlur?.(button.likelihood)}
          onClick={() => onChange?.(button.likelihood)}
          sx={{
            ...getButtonStyles(value, button.likelihood),
          }}
        >
          {button.label}
        </Button>
      ))}
    </ButtonGroup>
  );
};

const FormLikelihoodButtonGroup = (props: FormElementProps) => {
  return (
    <FormElement {...props}>
      <LikelihoodButtonGroup />
    </FormElement>
  );
};
