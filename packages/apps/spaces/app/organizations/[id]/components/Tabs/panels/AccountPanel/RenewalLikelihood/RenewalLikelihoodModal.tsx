'use client';
import { useState, useEffect, useRef } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Button, ButtonGroup } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';
import { Dot } from '@ui/media/Dot';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { useQueryClient } from '@tanstack/react-query';
import { useParams } from 'next/navigation';
import { useUpdateRenewalLikelihoodMutation } from '@organization/graphql/updateRenewalLikelyhood.generated';
import {
  RenewalLikelihood,
  RenewalLikelihoodProbability,
  User,
} from '@graphql/types';
import {
  OrganizationAccountDetailsQuery,
  useOrganizationAccountDetailsQuery,
} from '@organization/graphql/getAccountPanelDetails.generated';
import { useSession } from 'next-auth/react';
import { NEW_DATE } from '@organization/components/Timeline/OrganizationTimeline';
import { useInfiniteGetTimelineQuery } from '@organization/graphql/getTimeline.generated';

interface RenewalLikelihoodModalProps {
  isOpen: boolean;
  onClose: () => void;
  renewalLikelihood: RenewalLikelihood;
  name: string;
}

export const RenewalLikelihoodModal = ({
  renewalLikelihood,
  isOpen,
  onClose,
  name,
}: RenewalLikelihoodModalProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const id = useParams()?.id as string;
  const [probability, setLikelihood] = useState<
    RenewalLikelihoodProbability | undefined | null
  >(renewalLikelihood?.probability);
  const [reason, setReason] = useState<string>(
    renewalLikelihood?.comment || '',
  );
  const { data: session } = useSession();

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const updateRenewalLikelihood = useUpdateRenewalLikelihoodMutation(client, {
    onSuccess: (data, variables) => {
      queryClient.setQueryData<OrganizationAccountDetailsQuery>(
        useOrganizationAccountDetailsQuery.getKey({ id }),
        (oldData) => {
          if (!oldData || !oldData?.organization) return;
          return {
            organization: {
              ...(oldData?.organization ?? {}),
              accountDetails: {
                ...(oldData?.organization?.accountDetails ?? {}),
                renewalLikelihood: {
                  comment: reason,
                  previousProbability: renewalLikelihood?.probability,
                  probability: probability,
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
        (oldData: any) => {
          const newEvent = {
            __typename: 'Action',
            id: `timeline-event-action-new-id-${new Date()}`,
            actionType: 'RENEWAL_LIKELIHOOD_UPDATED',
            appSource: 'customer-os-api',
            createdAt: new Date(),
            metadata: JSON.stringify({
              likelihood: probability,
              reason,
            }),
            actionCreatedBy: null,
            content: `Renewal likelihood set to ${probability} by ${session?.user?.name}`,
          };

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

          const firstPage = oldData.pages[0] ?? {};
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
    updateRenewalLikelihood.mutate({
      input: { id, probability: probability, comment: reason },
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
    <Modal isOpen={isOpen} onClose={onClose}>
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
            {`${
              !renewalLikelihood?.probability ? 'Set' : 'Update'
            } renewal likelihood`}
          </Heading>
          <Text mt='1' fontSize='sm' fontWeight='normal'>
            {!renewalLikelihood?.probability ? 'Setting' : 'Updating'}{' '}
            <b>{name}</b> renewal likelihood will change how its renewal
            estimates are calculated and actions are prioritised.
          </Text>
        </ModalHeader>
        <ModalBody as={Flex} flexDir='column' pb='0'>
          <ButtonGroup w='full' isAttached>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='success' />}
              onClick={() => setLikelihood(RenewalLikelihoodProbability.High)}
              bg={probability === 'HIGH' ? 'gray.100' : 'white'}
            >
              High
            </Button>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='warning' />}
              onClick={() => setLikelihood(RenewalLikelihoodProbability.Medium)}
              bg={probability === 'MEDIUM' ? 'gray.100' : 'white'}
            >
              Medium
            </Button>
            <Button
              w='full'
              variant='outline'
              leftIcon={<Dot colorScheme='error' />}
              onClick={() => setLikelihood(RenewalLikelihoodProbability.Low)}
              bg={probability === 'LOW' ? 'gray.100' : 'white'}
            >
              Low
            </Button>
            <Button
              variant='outline'
              w='full'
              leftIcon={<Dot />}
              onClick={() => setLikelihood(RenewalLikelihoodProbability.Zero)}
              bg={probability === 'ZERO' ? 'gray.100' : 'white'}
            >
              Zero
            </Button>
          </ButtonGroup>

          {!!probability && (
            <>
              <Text as='label' htmlFor='reason' mt='5' fontSize='sm'>
                <b>Reason for change</b> (optional)
              </Text>
              <AutoresizeTextarea
                pt='0'
                id='reason'
                value={reason}
                spellCheck='false'
                onChange={(e) => setReason(e.target.value)}
                placeholder={`What is the reason for ${
                  !renewalLikelihood?.probability ? 'setting' : 'updating'
                } the renewal likelihood?`}
              />
            </>
          )}
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
            {!renewalLikelihood?.probability ? 'Set' : 'Update'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
