'use client';

import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';

export const BankTransferCard = () => {
  const formId = 'bank-transfer-form';

  return (
    <>
      <Card
        py={2}
        px={4}
        borderRadius='lg'
        boxShadow='none'
        border='1px solid'
        borderColor='gray.200'
        _hover={{
          '& #help-button': {
            visibility: 'visible',
          },
        }}
      >
        <CardHeader p='0' pb={1}>
          <FormInput
            fontSize='md'
            fontWeight='semibold'
            autoComplete='off'
            label='Sort code'
            placeholder='Sort code'
            name='name'
            formId={formId}
            border='none'
            _hover={{
              border: 'none',
            }}
            _focus={{
              border: 'none',
            }}
            _focusVisible={{
              border: 'none',
            }}
          />
        </CardHeader>
        <CardBody p={0} gap={2}>
          <Flex pb={1}>
            <FormInput
              autoComplete='off'
              label='Sort code'
              placeholder='Sort code'
              isLabelVisible
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                fontWeight: 'semibold',
              }}
              name='sortCode'
              formId={formId}
            />
            <FormInput
              autoComplete='off'
              label='Account number'
              placeholder='Bank account #'
              isLabelVisible
              labelProps={{
                fontSize: 'sm',
                mb: 0,
                fontWeight: 'semibold',
              }}
              name='accountNumber'
              formId={formId}
            />
          </Flex>
          <FormInput
            autoComplete='off'
            label='BIC/Swift'
            placeholder='BIC/Swift'
            isLabelVisible
            labelProps={{
              fontSize: 'sm',
              mb: 0,
              fontWeight: 'semibold',
            }}
            name='bic'
            formId={formId}
          />
          <FormInput
            autoComplete='off'
            label='Other details'
            placeholder='Other details'
            name='comments'
            formId={formId}
          />
        </CardBody>
      </Card>
    </>
  );
};
