'use client';
import { useRef, useState } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { RenewalCycle } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon';
import { FormSelect } from '@ui/form/SyncSelect';
import { Heading } from '@ui/typography/Heading';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { FormNumberInput } from '@ui/form/NumberInput';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { CurrencyInput } from '@ui/form/CurrencyInput';
import { ClockCheck } from '@ui/media/icons/ClockCheck';
import { SelectOption } from '@shared/types/SelectOptions';
import { Certificate02 } from '@ui/media/icons/Certificate02';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { frequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

export type SubscriptionServiceValue = {
  licenses?: string | null;
  description?: string | null;
  licensePrice?: string | null;
};

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
  data: SubscriptionServiceValue;
}

export const SubscriptionServiceModal = ({
  data,
  isOpen,
  onClose,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const formId = 'TODO';
  const [billingFrequency, setBillingFrequency] = useState<
    SelectOption<RenewalCycle>
  >(frequencyOptions[2]);
  const [licensePrice, setLicensePrice] = useState<string>(
    data?.licensePrice || '',
  );
  const [licenses, setLicenses] = useState<string>(data?.licenses || '');
  const [description, setDescription] = useState<string>(
    data?.description || '',
  );

  const handleSetSubscriptionServiceData = () => {
    // todo COS-857
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
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotSingle color='primary.600' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Add a new subscription service
          </Heading>
        </ModalHeader>
        <ModalBody pb='0' mb={4}>
          <Flex gap={2} mb={2}>
            <FormSelect
              label='Billed'
              tabIndex={1}
              isLabelVisible
              name='billingFrequency'
              formId='tbd'
              value={billingFrequency}
              onChange={(d) => setBillingFrequency(d)}
              options={frequencyOptions}
              leftElement={<ClockCheck mr='3' color='gray.500' />}
            />

            <FormNumberInput
              onChange={setLicenses}
              value={`${licenses}`}
              w='full'
              placeholder='Licences'
              isLabelVisible
              label='Licences'
              min={0}
              ref={initialRef}
              leftElement={
                <Box color='gray.500'>
                  <Certificate02 height='16px' />
                </Box>
              }
              formId={formId}
              name='licences'
            />

            <CurrencyInput
              onChange={setLicensePrice}
              value={`${licensePrice}`}
              w='full'
              placeholder='Price'
              isLabelVisible
              label='Price/license'
              min={0}
              ref={initialRef}
              leftElement={
                <Box color='gray.500'>
                  <CurrencyDollar height='16px' />
                </Box>
              }
            />
          </Flex>

          <AutoresizeTextarea
            pt='0'
            id='description'
            value={description}
            label='Description (Optional)'
            isLabelVisible
            spellCheck='false'
            onChange={(e) => setDescription(e.target.value)}
            placeholder='What is this service about?'
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
            onClick={handleSetSubscriptionServiceData}
          >
            Add
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
