'use client';
import { useParams } from 'next/navigation';
import { useRef, useState, useEffect } from 'react';

import { useSession } from 'next-auth/react';
import { useQueryClient } from '@tanstack/react-query';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { User, RenewalLikelihoodProbability } from '@graphql/types';
import { NEW_DATE } from '@organization/src/components/Timeline/OrganizationTimeline';
import { useInfiniteGetTimelineQuery } from '@organization/src/graphql/getTimeline.generated';
import { invalidateAccountDetailsQuery } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { useUpdateRenewalForecastMutation } from '@organization/src/graphql/updateRenewalForecast.generated';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';
import {
  OrganizationAccountDetailsQuery,
  useOrganizationAccountDetailsQuery,
} from '@organization/src/graphql/getAccountPanelDetails.generated';

export type RenewalForecastValue = {
  amount?: string | null;
  comment?: string | null;
};

interface RenewalForecastModalProps {
  name: string;
  isOpen: boolean;
  onClose: () => void;
  renewalForecast: RenewalForecastValue;
  renewalProbability?: RenewalLikelihoodProbability | null;
}

export const RenewalForecastModal = ({
  renewalForecast,
  renewalProbability,
  isOpen,
  onClose,
  name,
}: RenewalForecastModalProps) => {
  const id = useParams()?.id as string;
  const initialRef = useRef(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const { data: session } = useSession();

  const [amount, setAmount] = useState<string>(renewalForecast?.amount || '');
  const [reason, setReason] = useState<string>(renewalForecast?.comment || '');
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const updateRenewalForecast = useUpdateRenewalForecastMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
      queryClient.setQueryData<OrganizationAccountDetailsQuery>(
        useOrganizationAccountDetailsQuery.getKey({ id }),
        (oldData) => {
          if (!oldData || !oldData?.organization) return;

          return {
            organization: {
              ...(oldData?.organization ?? {}),
              accountDetails: {
                ...(oldData?.organization?.accountDetails ?? {}),
                renewalForecast: {
                  comment: reason,
                  amount: amount as unknown as number,
                  potentialAmount: null,
                  updatedAt: new Date(),
                  updatedBy: [session?.user] as unknown as User,
                },
              },
            },
          };
        },
      );

      queryClient.setQueryData(
        useInfiniteGetTimelineQuery.getKey({
          organizationId: id,
          from: NEW_DATE,
          size: 50,
        }),
        (oldData) => {
          const newEvent = {
            __typename: 'Action',
            id: `timeline-event-action-new-id-${new Date()}`,
            actionType: 'RENEWAL_FORECAST_UPDATED',
            appSource: 'customer-os-api',
            createdAt: new Date(),
            metadata: JSON.stringify({
              likelihood: renewalProbability,
              reason: reason,
            }),
            actionCreatedBy: null,
            content: `ARR forecast set to $${amount} by ${session?.user?.name}`,
          };

          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          if (!oldData || !oldData.pages?.length) {
            return {
              pages: [
                {
                  organization: {
                    id,
                    timelineEventsTotalCount: 1,
                    timelineEvents: [newEvent],
                  },
                },
              ],
            };
          }

          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          const firstPage = oldData.pages[0] ?? {};
          // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
          const pages = oldData.pages?.slice(1);

          const firstPageWithEvent = {
            ...firstPage,
            organization: {
              ...firstPage?.organization,
              timelineEvents: [
                ...(firstPage?.organization?.timelineEvents ?? []),
                newEvent,
              ],
              timelineEventsTotalCount:
                (firstPage?.organization?.timelineEventsTotalCount ?? 0) + 1,
            },
          };

          return {
            ...oldData,
            pages: [firstPageWithEvent, ...pages],
          };
        },
      );
    },
  });

  const handleSet = () => {
    updateRenewalForecast.mutate({
      input: {
        id,
        amount: (amount as unknown as number) || null,
        comment: reason,
      },
    });
    onClose();
  };

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Modal isOpen={isOpen} onClose={onClose} initialFocusRef={initialRef}>
      <ModalOverlay />
      <ModalContent
        borderRadius='2xl'
        backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
        backgroundRepeat='no-repeat'
        sx={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='warning'>
            <Icons.AlertTriangle />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            {`${!renewalForecast.amount ? 'Set' : 'Update'} ARR forecast`}
          </Heading>
          <Text mt='1' fontSize='sm' fontWeight='normal'>
            {!renewalForecast.amount ? 'Setting' : 'Updating'} <b>{name}</b> ARR
            forecast will change how expected revenue is reported.
          </Text>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' pb='0'>
          <CurrencyInput
            onChange={setAmount}
            value={`${amount}`}
            w='full'
            placeholder='Amount'
            label='Amount'
            min={0}
            ref={initialRef}
            leftElement={
              <Box color='gray.500'>
                <CurrencyDollar height='16px' />
              </Box>
            }
          />

          <Text as='label' htmlFor='reason' mt='4' fontSize='sm'>
            <b>Reason for change</b> (optional)
          </Text>
          <AutoresizeTextarea
            pt='0'
            id='reason'
            value={reason}
            spellCheck='false'
            onChange={(e) => setReason(e.target.value)}
            placeholder={`What is the reason for ${
              !renewalForecast.amount ? 'setting' : 'updating'
            } the ARR forecast?`}
          />
        </ModalBody>
        <ModalFooter p='6'>
          <Button variant='outline' w='full' onClick={onClose}>
            Cancel
          </Button>
          <Button
            ml='3'
            w='full'
            variant='outline'
            colorScheme='primary'
            onClick={handleSet}
          >
            {!renewalForecast.amount ? 'Set' : 'Update'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
