import { useForm } from 'react-inverted-form';
import React, { useRef, useState, useEffect } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Edit03 } from '@ui/media/icons/Edit03';
import { UseDisclosureReturn } from '@ui/utils';
import { File02 } from '@ui/media/icons/File02';
import { FormSelect } from '@ui/form/SyncSelect';
import { IconButton } from '@ui/form/IconButton';
import { Heading } from '@ui/typography/Heading';
import { Collapse } from '@ui/transitions/Collapse';
import { ChevronUp } from '@ui/media/icons/ChevronUp';
import { DatePicker } from '@ui/form/DatePicker/DatePicker';
import { Card, CardBody, CardFooter, CardHeader } from '@ui/presentation/Card';
import { billingFrequencyOptions } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { ContractDTO } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Contract.dto';
import { ServiceModal } from '@organization/src/components/Tabs/panels/AccountPanel/Services/ServiceModal';
import { ServicesList } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/ServicesList';
import { RenewalARRCard } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/RenewalARRCard';
import { ContractStatusSelect } from '@organization/src/components/Tabs/panels/AccountPanel/Contract/contractStatuses/ContractStatusSelect';

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
  const [isExpanded, setIsExpanded] = useState(false);

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
    defaultValues.renewalCycle,
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
      <CardHeader as={Flex} p='0' pb={2} w='full' flexDir='column'>
        <Flex justifyContent='space-between' w='full'>
          <Heading size='sm' color='gray.700' noOfLines={1}>
            {name} ContractCard
          </Heading>
          <Flex alignItems='center' gap={2}>
            <File02 color='gray.400' />
            <ContractStatusSelect />

            {isExpanded && (
              <IconButton
                size='xs'
                variant='ghost'
                aria-label='Show validationresults'
                icon={<ChevronUp />}
                onClick={() => setIsExpanded(false)}
              />
            )}
          </Flex>
        </Flex>

        {!isExpanded && (
          <Button
            bg='transparent'
            _hover={{
              bg: 'transparent',
              svg: { opacity: 1, transition: 'opacity 0.2s linear' },
            }}
            sx={{ svg: { opacity: 0, transition: 'opacity 0.2s linear' } }}
            size='xs'
            fontSize='sm'
            fontWeight='normal'
            color='gray.500'
            p={0}
            alignItems='center'
            justifyContent='flex-start'
            onClick={() => setIsExpanded(true)}
          >
            Service starts 1 Now 2023
            <Edit03 ml={1} color='gray.400' boxSize='3' />
          </Button>
        )}
      </CardHeader>
      {isExpanded && (
        <Collapse
          in={isExpanded}
          style={{ overflow: 'unset' }}
          delay={{
            exit: 2,
          }}
        >
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

              <FormSelect
                label='Contract renews'
                placeholder='Contract renews'
                isLabelVisible
                name='renewalCycle'
                formId={formId}
                options={billingFrequencyOptions}
              />
            </Flex>
            {/*<Divider  />*/}
          </CardBody>
        </Collapse>
      )}

      <CardFooter p='0' w='full' flexDir='column'>
        <RenewalARRCard />

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
