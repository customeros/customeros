import { useForm } from 'react-inverted-form';
import React, { useRef, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { UseDisclosureReturn } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Divider } from '@ui/presentation/Divider';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { Card, CardBody, CardFooter, CardHeader } from '@ui/presentation/Card';
import { ContractDTO } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Contract.dto';
import { ServiceModal } from '@organization/src/components/Tabs/panels/AccountPanel/Services/ServiceModal';
import { ServicesList } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServicesList';
import { ContractStatusDropdown } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ContractStatusDropdown';

interface ContractCardProps {
  data?: null; // todo when BE contract is available
  name?: string;
  serviceModal: UseDisclosureReturn;
}
export const ContractCard = ({
  data,
  serviceModal,
  name = '',
}: ContractCardProps) => {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const formId = 'contractForm';

  const defaultValues = ContractDTO.toForm(data);

  const { setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        return next;
      }

      return next;
    },
  });

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [
    defaultValues.contractSigned?.toISOString(),
    defaultValues.contractRenews?.toISOString(),
    defaultValues.contractEnds?.toISOString(),
    defaultValues.serviceStarts?.toISOString(),
  ]);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Card
      px='4'
      py='3'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      border='1px solid'
      borderColor='gray.200'
      bg='gray.50'
    >
      <CardHeader
        as={Flex}
        p='0'
        pb={3}
        w='full'
        justifyContent='space-between'
      >
        <Heading size='sm' color='gray.700' noOfLines={1}>
          {name} ContractCard
        </Heading>
        <Flex>
          <File02 color='gray.400' mr={4} />
          {/*<Check color='gray.400' />*/}
          <ContractStatusDropdown />
        </Flex>
      </CardHeader>

      <CardBody as={Flex} p='0' flexDir='column' w='full'>
        <Flex gap='4' mb={2}>
          <DatePicker
            label='Contract signed'
            placeholder='Signed date'
            formId={formId}
            name='contractSigned'
            calendarIconHidden
            inset='120% auto auto 0px'
          />
          <DatePicker
            label='Contract ends'
            placeholder='End date'
            formId={formId}
            name='contractEnds'
            calendarIconHidden
          />
        </Flex>
        <Flex gap='4'>
          <DatePicker
            label='Service starts'
            placeholder='Start date'
            formId={formId}
            name='serviceStarts'
            calendarIconHidden
            inset='120% auto auto 0px'
          />
          <DatePicker
            label='Contract renews'
            placeholder='Renewal date'
            formId={formId}
            name='contractRenews'
            calendarIconHidden
          />
        </Flex>
        <Divider my='2' />
      </CardBody>

      <CardFooter p='0' w='full'>
        <Flex w='full' alignItems='center' justifyContent='space-between'>
          <Text fontWeight='semibold' fontSize='sm'>
            No services
          </Text>
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Add service'
            color='gray.400'
            onClick={() => serviceModal.onOpen()}
            icon={<Plus boxSize='4' />}
          />
        </Flex>
        <ServicesList />
      </CardFooter>
      <ServiceModal
        isOpen={serviceModal.isOpen}
        onClose={serviceModal.onClose}
        data={{}}
      />
    </Card>
  );
};
