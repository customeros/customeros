import { useParams } from 'react-router-dom';
import { useForm } from 'react-inverted-form';
import { SelectInstance } from 'react-select';
import { components, OptionProps } from 'react-select';
import { useRef, useEffect, FormEvent, useCallback } from 'react';

import set from 'lodash/set';
import { produce } from 'immer';
import { match } from 'ts-pattern';
import { useQueryClient } from '@tanstack/react-query';

import { SelectOption } from '@ui/utils/types';
import { Button } from '@ui/form/Button/Button';
import { Flag04 } from '@ui/media/icons/Flag04';
import { toastError } from '@ui/presentation/Toast';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { OnboardingStatus, OnboardingDetails } from '@graphql/types';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
import { useUpdateOnboardingStatusMutation } from '@organization/graphql/updateOnboardingStatus.generated';
import {
  OrganizationQuery,
  useOrganizationQuery,
} from '@organization/graphql/organization.generated';
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

import { options } from './util';
import {
  OnboardingStatusDto,
  OnboardingStatusForm,
} from './OboardingStatus.dto';

interface OnboardingStatusModalProps {
  isOpen: boolean;
  onClose: () => void;
  data?: OnboardingDetails | null;
  onFetching?: (status: boolean) => void;
}

const formId = 'onboarding-status-update-form';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getIconcolorScheme: any = (status: OnboardingStatus) =>
  match(status)
    .returnType<string>()
    .with(
      OnboardingStatus.Successful,
      OnboardingStatus.OnTrack,
      OnboardingStatus.Done,
      () => 'success',
    )
    .with(OnboardingStatus.Late, OnboardingStatus.Stuck, () => 'warning')
    .otherwise(() => 'gray');

const getIcon = (status: OnboardingStatus) => {
  const color = `${getIconcolorScheme(status)}.500`;

  return match(status)
    .with(OnboardingStatus.Successful, () => (
      <Trophy01 className='mr-3' color={color} />
    ))
    .otherwise(() => <Flag04 color={color} className='mr-3' />);
};

export const OnboardingStatusModal = ({
  data,
  isOpen,
  onClose,
  onFetching,
}: OnboardingStatusModalProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const id = useParams()?.id as string;
  const timeout = useRef<NodeJS.Timeout>();
  const queryKey = useOrganizationQuery.getKey({ id });
  const initialFocusRef = useRef<SelectInstance>(null);

  const [timelineMeta] = useTimelineMeta();
  const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
    timelineMeta.getTimelineVariables,
  );

  const updateOnboardingStatus = useUpdateOnboardingStatusMutation(client, {
    onMutate: ({ input }) => {
      onFetching?.(true);
      queryClient.cancelQueries({ queryKey });

      const { previousEntries } = useOrganizationQuery.mutateCacheEntry(
        queryClient,
        { id },
      )((currentCache) => {
        return produce(currentCache, (draft) => {
          if (!draft) return;
          set(draft, 'organization.accountDetails.onboarding', {
            ...input,
            updatedAt: new Date().toISOString(),
          });
        });
      });

      return { previousEntries };
    },
    onError: (_, __, context) => {
      queryClient.setQueryData<OrganizationQuery>(
        queryKey,
        context?.previousEntries,
      );
      toastError(
        `We couldn't update the onboarding status`,
        `${id}-onboarding-status-update-error`,
      );
      onFetching?.(false);
    },
    onSettled: () => {
      onClose();
      if (timeout.current) clearTimeout(timeout.current);

      timeout.current = setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
        queryClient.invalidateQueries({ queryKey: timelineQueryKey });
        onFetching?.(false);
      }, 500);
    },
  });

  const onSubmit = useCallback(
    async (values: OnboardingStatusForm) => {
      updateOnboardingStatus.mutate({
        input: OnboardingStatusDto.toPayload({ id, ...values }),
      });
    },
    [updateOnboardingStatus.mutate],
  );

  const defaultValues = OnboardingStatusDto.toForm(data);
  const { state, handleSubmit, setDefaultValues } =
    useForm<OnboardingStatusForm>({
      formId,
      defaultValues,
      onSubmit,
      stateReducer: (_, action, next) => {
        if (action.type === 'HAS_SUBMITTED') {
          return { ...next, values: { ...next.values, comments: '' } };
        }

        return next;
      },
    });

  useEffect(() => {
    if (isOpen) {
      initialFocusRef.current?.focus();
    }
  }, [isOpen]);

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, []);

  return (
    <Modal open={isOpen} onOpenChange={onClose}>
      <ModalPortal>
        <ModalOverlay>
          <ModalContent
            className='rounded-2xl'
            style={{
              backgroundPositionX: '1px',
              backgroundPositionY: '-7px',
              backgroundImage: `url('/backgrounds/organization/circular-bg-pattern.png')`,
              backgroundRepeat: 'no-repeat',
            }}
          >
            <ModalCloseButton />
            <ModalHeader>
              <FeaturedIcon
                size='lg'
                colorScheme={getIconcolorScheme(
                  state?.values?.status?.value ??
                    OnboardingStatus.NotApplicable,
                )}
                className='ml-[12px] mt-[5px]'
              >
                {state?.values?.status?.value ===
                OnboardingStatus.Successful ? (
                  <Trophy01 />
                ) : (
                  <Flag04 />
                )}
              </FeaturedIcon>
              <h2 className='text-lg mt-6'>Update onboarding status</h2>
            </ModalHeader>
            <ModalBody className='gap-4 flex flex-col'>
              <FormSelect
                name='status'
                label='Status'
                isLabelVisible
                formId={formId}
                options={options}
                ref={initialFocusRef}
                openMenuOnFocus={true}
                isDisabled={updateOnboardingStatus.isPending}
                components={{ Option }}
              />
              {defaultValues.status.value !== state?.values?.status?.value && (
                <div>
                  <label className='text-sm' htmlFor='reason'>
                    <b>Reason for change</b> (optional)
                  </label>
                  <FormAutoresizeTextarea
                    formId={formId}
                    name='comments'
                    spellCheck='false'
                    disabled={updateOnboardingStatus.isPending}
                    placeholder={`What is the reason for changing the onboarding status?`}
                  />
                </div>
              )}
            </ModalBody>
            <ModalFooter className='p-6 flex '>
              <Button
                className='w-full'
                variant='outline'
                onClick={onClose}
                isDisabled={updateOnboardingStatus.isPending}
              >
                Cancel
              </Button>
              <Button
                className='w-full ml-3'
                onClick={() => handleSubmit({} as FormEvent<HTMLFormElement>)}
                variant='outline'
                colorScheme='primary'
                isLoading={updateOnboardingStatus.isPending}
              >
                Update status
              </Button>
            </ModalFooter>
          </ModalContent>
        </ModalOverlay>
      </ModalPortal>
    </Modal>
  );
};

export const Option = ({
  data,
  children,
  ...rest
}: OptionProps<SelectOption<OnboardingStatus>>) => {
  const icon = getIcon(data.value);

  return (
    <components.Option data={data} {...rest}>
      {icon}
      {children}
    </components.Option>
  );
};
