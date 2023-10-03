import { useRef, useEffect, useState } from 'react';
import { useForm } from 'react-inverted-form';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Divider } from '@ui/presentation/Divider';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import CurrencyDollar from '@spaces/atoms/icons/CurrencyDollar';
import { FormSelect } from '@ui/form/SyncSelect';
import CoinsSwap from '@spaces/atoms/icons/CoinsSwap';
import { Card, CardBody, CardFooter } from '@ui/layout/Card';
import { FormCurrencyInput } from '@ui/form/CurrencyInput/FormCurrencyInput';
import { BillingDetailsForm, BillingDetailsDTO } from './BillingDetails.dto';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { invalidateAccountDetailsQuery } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { useQueryClient } from '@tanstack/react-query';
import { useUpdateBillingDetailsMutation } from '@organization/src/graphql/updateBillingDetails.generated';
import { BillingDetails as BillingDetails } from '@graphql/types';

import { frequencyOptions } from '../utils';

interface BillingDetailsCardBProps {
  id: string;
  data?: BillingDetails | null;
}
export const BillingDetailsCard = ({ id, data }: BillingDetailsCardBProps) => {
  const [isFocused, setIsFocused] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();
  const defaultValues = BillingDetailsDTO.toForm(data);
  const client = getGraphQLClient();
  const updateBillingDetails = useUpdateBillingDetailsMutation(client, {
    onSuccess: () => {
      timeoutRef.current = setTimeout(
        () => invalidateAccountDetailsQuery(queryClient, id),
        500,
      );
    },
  });
  const formId = 'billing-details-form';

  const handleUpdateBillingDetails = (variables: BillingDetailsForm) => {
    const payload = BillingDetailsDTO.toPayload({
      organizationId: id,
      ...data,
      ...variables,
    });

    updateBillingDetails.mutate({
      input: {
        ...payload,
      },
    });
  };

  useForm<BillingDetailsForm>({
    formId,
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_BLUR') {
        setIsFocused(false);
      }
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'frequency': {
            handleUpdateBillingDetails({
              ...state.values,
              frequency: action.payload?.value,
            });

            return next;
          }
          default:
            return next;
        }
      }

      if (action.type === 'FIELD_BLUR' && action.payload.name === 'amount') {
        handleUpdateBillingDetails({
          ...state.values,
          amount: action.payload.value,
        });
      }

      return next;
    },
  });

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <Card
      p='4'
      w='full'
      size='lg'
      variant='outline'
      cursor='default'
      boxShadow={isFocused ? 'md' : 'xs'}
      _hover={{
        boxShadow: 'md',
      }}
      transition='all 0.2s ease-out'
    >
      <CardBody as={Flex} p='0' w='full' align='center'>
        <FeaturedIcon>
          <Icons.Coin1 />
        </FeaturedIcon>
        <Heading ml='5' size='sm' color='gray.700'>
          Billing details
        </Heading>
      </CardBody>

      <CardFooter as={Flex} flexDir='column' padding={0}>
        <Divider color='gray.200' my='4' />
        <Flex w='full' flexDir='column'>
          <Flex justifyItems='space-between' w='full' gap='4'>
            <FormCurrencyInput
              label='Billing amount'
              color='gray.700'
              isLabelVisible
              formId={formId}
              name='amount'
              min={0}
              onFocus={() => setIsFocused(true)}
              placeholder='Amount'
              leftElement={
                <Box color='gray.500'>
                  <CurrencyDollar height='16px' />
                </Box>
              }
            />

            <FormSelect
              isClearable
              label='Billing frequency'
              isLabelVisible
              name='frequency'
              placeholder='Monthly'
              options={frequencyOptions}
              formId={formId}
              onFocus={() => setIsFocused(true)}
              leftElement={
                <Box mr={3} color='gray.500'>
                  <CoinsSwap height={16} />
                </Box>
              }
            />
          </Flex>
        </Flex>
      </CardFooter>
    </Card>
  );
};
