'use client';
import { useRef, useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { Button } from '@ui/form/Button';
import { FeaturedIcon } from '@ui/media/Icon';
import { ServiceLineItem } from '@graphql/types';
import { Heading } from '@ui/typography/Heading';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Tab, Tabs, TabList, TabPanel, TabPanels } from '@ui/disclosure/Tabs';
import { useCreateServiceMutation } from '@organization/src/graphql/createService.generated';
import { billedTypeOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { OneTimeServiceForm } from '@organization/src/components/Tabs/panels/AccountPanel/Services/OneTimeServiceForm';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';
import {
  ServiceDTO,
  ServiceForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Services/Service.dto';
import { SubscriptionServiceFrom } from '@organization/src/components/Tabs/panels/AccountPanel/Services/SubscriptionServiceForm';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  data?: ServiceLineItem;
  isSubscription?: boolean;
  mode?: 'create' | 'update';
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
  contractId,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const client = getGraphQLClient();
  const formId = `${mode}-service-item`;
  const createService = useCreateServiceMutation(client, {
    onSuccess: () => {
      onClose();
    },
  });

  const defaultValues = ServiceDTO.toForm(data);
  const { setDefaultValues, state } = useForm<ServiceForm>({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [
    defaultValues.renewalCycle,
    defaultValues.billed,
    defaultValues.appSource,
    defaultValues.quantity,
    defaultValues.serviceStartedAt,
    defaultValues.externalReference,
  ]);
  const handleSetSubscriptionServiceData = () => {
    createService.mutate({
      input: { ...ServiceDTO.toPayload(state.values, contractId) },
    });
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
                onClick={() => {
                  setDefaultValues({
                    ...defaultValues,
                  });
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
                onClick={() => {
                  setDefaultValues({
                    ...defaultValues,
                    billed: billedTypeOptions[0],
                  });
                }}
              >
                One-time
              </Tab>
            </TabList>

            <TabPanels>
              <TabPanel px={0} pb={2}>
                <SubscriptionServiceFrom formId={formId} />
              </TabPanel>
              <TabPanel px={0} pb={2}>
                <OneTimeServiceForm formId={formId} />
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
