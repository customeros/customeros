'use client';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { Users03 } from '@ui/media/icons/Users03';
import { Divider } from '@ui/presentation/Divider';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { Card, CardBody, CardHeader, CardFooter } from '@ui/layout/Card';

import { useGetBillableInfoQuery } from '../../../../graphql/getTenantBillableInfo.generated';

export const BillingPanel = () => {
  const client = getGraphQLClient();
  const { data } = useGetBillableInfoQuery(client);

  return (
    <Card
      flex='1'
      w='full'
      h='100vh'
      bg='#FCFCFC'
      borderRadius='2xl'
      flexDirection='column'
      boxShadow='none'
      background='gray.25'
    >
      <CardHeader px='6' pb='0' pt='4'>
        <Heading as='h1' fontSize='lg' color='gray.700'>
          <b>Billing</b>
        </Heading>
      </CardHeader>
      <CardBody px='6' w='full'>
        <Card
          p='4'
          w='full'
          maxW='23.5rem'
          size='lg'
          variant='outline'
          cursor='default'
          boxShadow='xs'
          _hover={{
            boxShadow: 'md',
          }}
          transition='all 0.2s ease-out'
        >
          <CardBody as={Flex} p='0' align='center'>
            <FeaturedIcon size='md' minW='10' colorScheme='gray'>
              <Users03 />
            </FeaturedIcon>
            <Flex
              ml='5'
              w='full'
              align='center'
              columnGap={4}
              justify='space-between'
            >
              <Heading
                size='sm'
                whiteSpace='nowrap'
                fontWeight='semibold'
                color='gray.700'
                mr={2}
              >
                Contacts
              </Heading>
            </Flex>
          </CardBody>

          <CardFooter p='0' as={Flex} flexDir='column'>
            <Divider mt='4' mb='2' />
            <Flex justify='space-between' align='center'>
              <Text color='gray.700'>Synced (no charge)</Text>
              <Text color='gray.700'>
                {data?.billableInfo.greylistedContacts ?? 0}
              </Text>
            </Flex>
            <Flex justify='space-between' align='center'>
              <Text color='gray.700'>Active (billed)</Text>
              <Text color='gray.700' fontWeight='semibold'>
                {data?.billableInfo.whitelistedContacts ?? 0}
              </Text>
            </Flex>
          </CardFooter>
        </Card>
      </CardBody>
    </Card>
  );
};
