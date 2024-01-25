'use client';
import React, { useRef, useState, Fragment, useEffect } from 'react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { FeaturedIcon } from '@ui/media/Icon';
// import { Grid, GridItem } from '@ui/layout/Grid';
import { Heading } from '@ui/typography/Heading';
import { DotSingle } from '@ui/media/icons/DotSingle';
import { AutoresizeTextarea } from '@ui/form/Textarea';
// import { Invoice } from '@shared/components/Invoice/Invoice';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import {
  BilledType,
  // InvoiceLine,
  ServiceLineItem,
} from '@graphql/types';
import { useUpdateServicesMutation } from '@organization/src/graphql/updateServiceLineItems.generated';
import { ServiceItem } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServiceLineItemsModal/type';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';

import { ServiceLineItemRow } from './ServiceLineItemRow';

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  contractId: string;
  onClose: () => void;
  contractName: string;
  organizationName: string;
  serviceLineItems: Array<ServiceLineItem>;
}

interface SubscriptionServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const defaultValue = {
  name: 'Unnamed',
  quantity: 1,
  price: 1,
  billed: BilledType.Monthly,
  type: 'RECURRING',
  isDeleted: false,
};
export const ServiceLineItemsModal = ({
  isOpen,
  onClose,
  serviceLineItems,
  contractId,
}: SubscriptionServiceModalProps) => {
  const initialRef = useRef(null);
  const client = getGraphQLClient();

  const updateServices = useUpdateServicesMutation(client);
  const [services, setServices] = useState<Array<ServiceItem>>([]);
  const [_, setNote] = useState<string>('');
  useEffect(() => {
    if (isOpen) {
      if (serviceLineItems.length) {
        setServices(
          serviceLineItems
            .filter((e) => !e.endedAt)
            .map((e) => ({
              ...e,
              isDeleted: false,
              type: [
                BilledType.Quarterly,
                BilledType.Monthly,
                BilledType.Annually,
              ].includes(e.billed)
                ? 'RECURRING'
                : e.billed,
            })),
        );
      }
      if (!serviceLineItems.length) {
        setServices([defaultValue]);
      }
    }
  }, [isOpen]);
  const handleAddService = () => {
    setServices([...services, defaultValue]);
  };

  const handleUpdateService = (index: number, updatedService: ServiceItem) => {
    const updatedServices = [...services];

    updatedServices[index] = updatedService;
    setServices(updatedServices);
  };

  const handleApplyChanges = async () => {
    updateServices.mutate({
      input: {
        contractId,
        serviceLineItems: services.map((e) => ({
          serviceLineItemId: e?.id ?? '',
          billed: e.billed,
          name: e.name,
          price: e.price,
          quantity: e.quantity,
        })),
      },
    });
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      initialFocusRef={initialRef}
      size='2xl'
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent borderRadius='2xl' maxH='90vh' overflow='hidden'>
        {/*<Grid h='100%' templateColumns='1fr' gap={4} overflow='scroll'>*/}
        {/*545px 1fr*/}
        {/*<GridItem*/}
        {/*  rowSpan={1}*/}
        {/*  colSpan={1}*/}
        {/*  h='100%'*/}
        {/*  display='flex'*/}
        {/*  flexDir='column'*/}
        {/*  justifyContent='space-between'*/}
        {/*  bg='gray.25'*/}
        {/*  borderRight='1px solid'*/}
        {/*  borderColor='gray.200'*/}
        {/*  borderRadius='2xl'*/}
        {/*  // borderTopLeftRadius='2xl'*/}
        {/*  // borderBottomLeftRadius='2xl'*/}
        {/*  backgroundImage='/backgrounds/organization/circular-bg-pattern.png'*/}
        {/*  backgroundRepeat='no-repeat'*/}
        {/*  sx={{*/}
        {/*    backgroundPositionX: '1px',*/}
        {/*    backgroundPositionY: '-7px',*/}
        {/*  }}*/}
        {/*>*/}
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <DotSingle color='primary.700' />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Modify contract service line items
          </Heading>
        </ModalHeader>
        <ModalBody
          pb='0'
          display='flex'
          flexDir='column'
          flex={1}
          overflow='scroll'
        >
          <Flex
            justifyContent='space-between'
            alignItems='center'
            pr='20px'
            borderBottom='1px solid'
            borderColor='gray.300'
            pb={1}
          >
            <Text fontSize='sm' fontWeight='medium' w='30%'>
              Name
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='20%'>
              Type
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='10%'>
              Qty
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='20%'>
              Price
            </Text>
            <Text fontSize='sm' fontWeight='medium' w='20%'>
              Recurring
            </Text>
          </Flex>

          {services.map((service, index) => (
            <Fragment key={`service-line-item-${index}`}>
              <ServiceLineItemRow
                service={service}
                index={index}
                onChange={(data) => handleUpdateService(index, data)}
              />
            </Fragment>
          ))}
          <Box>
            <Button
              leftIcon={<Plus />}
              variant='ghost'
              size='sm'
              p={1}
              my={1}
              color='gray.500'
              fontWeight='regular'
              onClick={handleAddService}
            >
              New item
            </Button>
          </Box>

          <AutoresizeTextarea
            label='Note'
            isLabelVisible
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            name='contractUrl'
            textOverflow='ellipsis'
            placeholder='Paste or enter a contract link'
            onChange={(event) => setNote(event.target.value)}
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
            Apply changes
          </Button>
        </ModalFooter>
        {/*</GridItem>*/}
        {/*<GridItem pr={3}>*/}
        {/*  <Invoice*/}
        {/*    isDraft*/}
        {/*    tax={10}*/}
        {/*    note={note}*/}
        {/*    from={{*/}
        {/*      address: '',*/}
        {/*      address2: '',*/}
        {/*      city: '',*/}
        {/*      country: '',*/}
        {/*      name: '',*/}
        {/*      zip: '',*/}
        {/*      email: '',*/}
        {/*    }}*/}
        {/*    total={1400}*/}
        {/*    dueDate={new Date()}*/}
        {/*    subtotal={2002}*/}
        {/*    lines={services as unknown as InvoiceLine[]}*/}
        {/*    issueDate='10.01.2024'*/}
        {/*    billedTo={{*/}
        {/*      address: '',*/}
        {/*      address2: '',*/}
        {/*      city: '',*/}
        {/*      country: '',*/}
        {/*      name: '',*/}
        {/*      zip: '',*/}
        {/*      email: '',*/}
        {/*    }}*/}
        {/*    invoiceNumber='INV-001'*/}
        {/*  />*/}
        {/*</GridItem>*/}
        {/*</Grid>*/}
      </ModalContent>
    </Modal>
  );
};
