'use client';

import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { VStack } from '@ui/layout/Stack';
import { Input } from '@ui/form/Input';

const DetailGroup = ({
  label,
  value,
  children,
}: {
  label: string;
  value?: string;
  children?: React.ReactNode;
}) => {
  return (
    <Flex flexDir='column' w='full'>
      <Text color='gray.800' fontWeight='semibold'>
        {label}
      </Text>
      {!children ? (
        <>
          <Input
            ml='-5px'
            px='1'
            variant='unstyled'
            border='1px solid transparent'
            borderRadius='2px'
            _hover={{
              borderColor: 'gray.400',
            }}
            _focus={{
              borderColor: 'blue.400',
            }}
            _active={{
              borderColor: 'blue.400',
            }}
            placeholder={label}
            defaultValue={value}
            w='50%'
          />
        </>
      ) : (
        children
      )}
    </Flex>
  );
};

export const OrganizationDetails = () => {
  return (
    <Flex h='full' flexDir='column'>
      <Flex justify='space-between' flex='1' align='flex-end' mb='4'>
        <Heading fontSize='2xl'>ACME Org</Heading>
        <Text color='gray.600'>customeros.aa</Text>
      </Flex>

      <Flex minH='20' mb='6'>
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Possimus
        doloribus velit aperiam porro odio cupiditate ullam culpa molestiae
        illum doloremque eaque sint, natus ducimus animi iusto. Aliquam error
        rem blanditiis.
      </Flex>

      <VStack spacing='4' flex='1' align='flex-start' justify='flex-start'>
        <DetailGroup label='Industry' value='Internet Software & Services' />
        <DetailGroup label='Target Audience' value='SaaS' />
        <DetailGroup label='Business Type' value='B2B' />
        <DetailGroup label='Last funding' value='Pre-seed • $1,000,000' />
        <DetailGroup label='Number of employees' value='2-20' />
        <DetailGroup label='Relationship & stage' value='Customer • Live' />
        <DetailGroup label='Social'>
          <Text>twitter.com/customeros</Text>
          <Text>linkedin.com/in/customeros</Text>
        </DetailGroup>
      </VStack>
    </Flex>
  );
};
