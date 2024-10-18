import { useParams } from 'react-router-dom';
import { SelectInstance } from 'react-select';
import { useRef, useState, useEffect } from 'react';
import { components, OptionProps } from 'react-select';

import set from 'lodash/set';
import { match } from 'ts-pattern';
import { observer } from 'mobx-react-lite';
import { useQueryClient } from '@tanstack/react-query';

import { SelectOption } from '@ui/utils/types';
import { Button } from '@ui/form/Button/Button';
import { Flag04 } from '@ui/media/icons/Flag04';
import { Select } from '@ui/form/Select/Select';
import { useStore } from '@shared/hooks/useStore';
import { OnboardingStatus } from '@graphql/types';
import { Trophy01 } from '@ui/media/icons/Trophy01';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { useTimelineMeta } from '@organization/components/Timeline/state';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';
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

interface OnboardingStatusModalProps {
  isOpen: boolean;
  onClose: () => void;
  onFetching?: (status: boolean) => void;
}

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
      <Trophy01 color={color} className='mr-3' />
    ))
    .otherwise(() => <Flag04 color={color} className='mr-3' />);
};

export const OnboardingStatusModal = observer(
  ({ isOpen, onClose }: OnboardingStatusModalProps) => {
    const store = useStore();
    const id = useParams()?.id as string;
    const organization = store.organizations.value.get(id);
    const onboardingDetails = organization?.value?.accountDetails?.onboarding;

    const [initialStatus, setInitialStatus] = useState<OnboardingStatus>(
      OnboardingStatus.NotApplicable,
    );

    const [comments, setComments] = useState<string>('');

    const queryClient = useQueryClient();
    const initialFocusRef = useRef<SelectInstance>(null);

    const [timelineMeta] = useTimelineMeta();
    const timelineQueryKey = useInfiniteGetTimelineQuery.getKey(
      timelineMeta.getTimelineVariables,
    );

    const handleSubmit = () => {
      if (!organization) return;

      set(organization.value, 'accountDetails.onboarding.comments', comments);
      set(
        organization.value,
        'accountDetails.onboarding.status',
        initialStatus,
      );

      organization.commit();

      onClose();

      setTimeout(() => {
        queryClient.invalidateQueries({
          queryKey: timelineQueryKey,
        });
      }, 1000);
    };

    useEffect(() => {
      if (isOpen) {
        initialFocusRef.current?.focus();
      }
    }, [isOpen]);

    useEffect(() => {
      setInitialStatus(
        onboardingDetails?.status ?? OnboardingStatus.NotApplicable,
      );
    }, []);

    return (
      <Modal open={isOpen} onOpenChange={onClose}>
        <ModalPortal>
          <ModalOverlay>
            <ModalContent className='rounded-2xl'>
              <ModalCloseButton />
              <ModalHeader>
                <h2 className='text-lg mt-6'>Update onboarding status</h2>
              </ModalHeader>
              <ModalBody className='gap-4 flex flex-col'>
                <div>
                  <label htmlFor='status' className='text-md'>
                    Status
                  </label>
                  <Select
                    id='status'
                    name='status'
                    options={options}
                    ref={initialFocusRef}
                    openMenuOnFocus={true}
                    components={{ Option }}
                    value={options.find((o) => o.value === initialStatus)}
                    onChange={(value) => {
                      setInitialStatus(value.value);
                    }}
                  />
                </div>
                {initialStatus !== onboardingDetails?.status && (
                  <div>
                    <label htmlFor='comments' className='text-sm'>
                      <b>Reason for change</b> (optional)
                    </label>
                    <AutoresizeTextarea
                      id='comments'
                      name='comments'
                      spellCheck='false'
                      value={comments ?? ''}
                      onChange={(e) => setComments(e.target.value)}
                      placeholder={`What is the reason for changing the onboarding status?`}
                    />
                  </div>
                )}
              </ModalBody>
              <ModalFooter className='p-6 flex '>
                <Button variant='outline' onClick={onClose} className='w-full'>
                  Cancel
                </Button>
                <Button
                  variant='outline'
                  colorScheme='primary'
                  onClick={handleSubmit}
                  className='w-full ml-3'
                >
                  Update status
                </Button>
              </ModalFooter>
            </ModalContent>
          </ModalOverlay>
        </ModalPortal>
      </Modal>
    );
  },
);

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
