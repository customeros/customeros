import { useForm } from 'react-inverted-form';
import { useRef, useState, useEffect } from 'react';

import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { FormSelect } from '@ui/form/SyncSelect';
import { Coins01 } from '@ui/media/icons/Coins01';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/layout/Card';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { BillingDetails as BillingDetails } from '@graphql/types';
import { FormCurrencyInput } from '@ui/form/CurrencyInput/FormCurrencyInput';
import { useUpdateBillingDetailsMutation } from '@organization/src/graphql/updateBillingDetails.generated';
import { invalidateAccountDetailsQuery } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

import { frequencyOptions } from '../utils';
import { BillingDetailsDTO, BillingDetailsForm } from './BillingDetails.dto';

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

  const { setDefaultValues } = useForm<BillingDetailsForm>({
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
    setDefaultValues(defaultValues);
  }, [defaultValues?.amount, defaultValues?.frequency?.value]);

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
              leftElement={<CurrencyDollar color='gray.500' />}
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
              leftElement={<Coins01 mr='3' color='gray.500' />}
            />
          </Flex>
        </Flex>
      </CardFooter>
    </Card>
  );
};
