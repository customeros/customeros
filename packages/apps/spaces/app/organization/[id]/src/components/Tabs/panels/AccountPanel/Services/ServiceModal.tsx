'use client';
import { useRef } from 'react';

import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { Tab, Tabs, TabList, TabPanel, TabPanels } from '@ui/disclosure/Tabs';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';
import { SubscriptionServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Services/SubscriptionServiceForm';
import {
  OneTimeServiceForm,
  OneTimeServiceValue,
} from '@organization/src/components/Tabs/panels/AccountPanel/Services/OneTimeServiceForm';

export type SubscriptionServiceValue = {
  licenses?: string | null;
  description?: string | null;
  licensePrice?: string | null;
};

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
  isSubscription?: boolean;
  mode?: 'create' | 'update';
  data: SubscriptionServiceValue | OneTimeServiceValue;
}

const copy = {
  create: {
    title: 'Add a new service',
    submit: 'Add',
  },
  update: {
    subscriptionTitle: 'Update subscription service', // todo
    oneTimeTitle: 'Update one-time service', // todo
    submit: 'Update',
  },
};

export const ServiceModal = ({
  data,
  isOpen,
  onClose,
  mode = 'create',
  isSubscription,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);

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
            {mode === 'create'
              ? copy[mode].title
              : isSubscription
              ? copy.update.subscriptionTitle
              : copy.update.oneTimeTitle}
          </Heading>
        </ModalHeader>
        <ModalBody pb='0'>
          <Tabs isFitted>
            <TabList
              bg='white'
              border='1px solid'
              borderColor='gray.300'
              borderRadius='md'
            >
              <Tab
                borderTopLeftRadius='md'
                borderBottomLeftRadius='md'
                borderBottom='none'
                flex={1}
                borderRight='1px solid'
                borderRightColor='gray.300 !important'
                color='gray.500'
                bg='gray.50'
                mb={0}
                _selected={{
                  color: 'gray.500',
                  bg: 'white',
                  fontWeight: 'semibold',
                }}
              >
                Subscription
              </Tab>
              <Tab
                borderTopRightRadius='md'
                borderBottomRightRadius='md'
                borderRadius='md'
                borderBottom='none'
                flex={1}
                mb={0}
                color='gray.500'
                bg='gray.50'
                _selected={{
                  color: 'gray.500',
                  bg: 'white',
                  fontWeight: 'semibold',
                }}
              >
                One-time
              </Tab>
            </TabList>

            <TabPanels>
              <TabPanel px={0} pb={2}>
                <SubscriptionServiceFrom data={data} />
              </TabPanel>
              <TabPanel px={0} pb={2}>
                <OneTimeServiceForm data={data as OneTimeServiceValue} />
              </TabPanel>
            </TabPanels>
          </Tabs>
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
            {copy[mode].submit}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
