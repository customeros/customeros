'use client';
import React, { useRef } from 'react';

import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { ModalBody } from '@ui/overlay/Modal';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { XSquare } from '@ui/media/icons/XSquare';
import { Radio, RadioGroup } from '@ui/form/Radio';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractQuery } from '@organization/src/graphql/getContract.generated';
import {
  Modal,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  organizationName?: string;
}

export const ContractEndModal = ({
  isOpen,
  onClose,
  contractId,
  organizationName,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const formId = `billing-details-form-${contractId}`;
  const client = getGraphQLClient();
  const [value, setValue] = React.useState('1');

  const { data } = useGetContractQuery(
    client,
    {
      id: contractId,
    },
    {
      enabled: isOpen && !!contractId,
      refetchOnMount: true,
    },
  );

  const handleApplyChanges = () => {};

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='md'
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl'>
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='error'>
            <XSquare color='error.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            End{' '}
            {data?.contract?.organizationLegalName ||
              organizationName ||
              "Unnamed's "}{' '}
            contract?
          </Heading>
        </ModalHeader>
        <ModalBody gap={4}>
          <Text>
            Ending this contract <b>will close the renewal</b> and set the{' '}
            <b>ARR to zero.</b>
          </Text>
          <Text>Let’s end it on:</Text>

          <RadioGroup
            value={value}
            onChange={setValue}
            flexDir='column'
            display='flex'
          >
            <Radio value={'1'} colorScheme='primary'>
              Now
            </Radio>
            <Radio value={'2'} colorScheme='primary'>
              End of current billing period, 12 Mar 2024
            </Radio>
            <Radio value={'3'} colorScheme='primary'>
              On a custom date
            </Radio>
          </RadioGroup>

          {/*<DatePicker*/}
          {/*  name='endDate'*/}
          {/*  formId='todo'*/}
          {/*  label='Contract end date'*/}
          {/*  placeholder='Contract end date'*/}
          {/*/>*/}
          <FormAutoresizeTextarea
            name=''
            formId=''
            label='Reason for change (optional)'
            isLabelVisible
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
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
            loadingText='Applying changes...'
            onClick={handleApplyChanges}
          >
            End contract
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
