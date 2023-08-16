'use client';
import React, { useRef, useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Button } from '@ui/form/Button';
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
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useQueryClient } from '@tanstack/react-query';
import { invalidateAccountDetailsQuery } from '@organization/components/Tabs/panels/AccountPanel/utils';
import { useParams } from 'next/navigation';
import { useUpdateRenewalForecastMutation } from '@organization/graphql/updateRenewalForecast.generated';
import { Box } from '@ui/layout/Box';
import CurrencyDollar from '@spaces/atoms/icons/CurrencyDollar';

export type RenewalForecastValue = {
  amount?: string | null;
  comment?: string | null;
};

interface RenewalForecastModalProps {
  isOpen: boolean;
  onClose: () => void;
  renewalForecast: RenewalForecastValue;
  name: string;
}

export const RenewalForecastModal = ({
  renewalForecast,
  isOpen,
  onClose,
  name,
}: RenewalForecastModalProps) => {
  const id = useParams()?.id as string;
  const initialRef = useRef(null);

  const [amount, setAmount] = useState<string>(renewalForecast?.amount || '');
  const [reason, setReason] = useState<string>(renewalForecast?.comment || '');
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const updateRenewalForecast = useUpdateRenewalForecastMutation(client, {
    onSuccess: () => invalidateAccountDetailsQuery(queryClient, id),
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
            {`${!renewalForecast.amount ? 'Set' : 'Update'} renewal forecast`}
          </Heading>
          <Text mt='1' fontSize='sm' fontWeight='normal'>
            {!renewalForecast.amount ? 'Setting' : 'Updating'} <b>{name}</b>{' '}
            renewal forecast will change how expected revenue is reported.
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
            } the renewal forecast?`}
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
