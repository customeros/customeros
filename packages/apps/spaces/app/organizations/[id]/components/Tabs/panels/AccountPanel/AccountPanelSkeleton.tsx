import React from 'react';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { Card, CardBody } from '@ui/presentation/Card';
import { Flex } from '@ui/layout/Flex';
import { FeaturedIcon, Icons } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Divider } from '@ui/presentation/Divider';
import { FormLabel, Skeleton } from '@chakra-ui/react';
import { CardHeader } from '@ui/layout/Card';
import BillingDetails from '@spaces/atoms/icons/BillingDetails';
import { Box } from '@ui/layout/Box';
import { VStack } from '@ui/layout/Stack';
import CurrencyDollar from '@spaces/atoms/icons/CurrencyDollar';
import CoinsSwap from '@spaces/atoms/icons/CoinsSwap';
import ClockCheck from '@spaces/atoms/icons/ClockCheck';
import Calendar from '@spaces/atoms/icons/Calendar';

export const AccountPanelSkeleton: React.FC = () => {
  return (
    <OrganizationPanel title='Account'>
      <Card p='4' w='full' size='lg' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' colorScheme='gray'>
            <Icons.Building7 />
          </FeaturedIcon>
          <Flex ml='5' align='center' justify='space-between' w='full'>
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading size='sm' fontWeight='semibold' color='gray.700'>
                  Renewal likelihood
                </Heading>
              </Flex>
              <Text fontSize='xs' color='gray.500'>
                <Skeleton height='10px' width='90px' borderRadius='md' mt={1} />
              </Text>
            </Flex>

            <Heading fontSize='2xl' color='gray'>
              <Skeleton height='40px' width='73px' borderRadius='md' />
            </Heading>
          </Flex>
        </CardBody>
      </Card>

      <Card p='4' w='full' size='lg' boxShadow='xs' variant='outline'>
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' colorScheme='gray'>
            <Icons.Building7 />
          </FeaturedIcon>
          <Flex ml='5' align='center' justify='space-between' w='full'>
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading size='sm' fontWeight='semibold' color='gray.700'>
                  Renewal Forecast
                </Heading>
              </Flex>
              <Text fontSize='xs' color='gray.500'>
                <Skeleton height='10px' width='90px' borderRadius='md' mt={1} />
              </Text>
            </Flex>

            <Heading fontSize='2xl' color='gray'>
              <Skeleton height='40px' width='96px' borderRadius='md' />
            </Heading>
          </Flex>
        </CardBody>
      </Card>

      <Card
        size='sm'
        width='full'
        borderRadius='xl'
        border='1px solid'
        borderColor='gray.200'
        boxShadow='xs'
      >
        <CardHeader display='flex' alignItems='center'>
          <BillingDetails />
          <Heading ml={5} size='sm' color='gray.700'>
            Billing details
          </Heading>
        </CardHeader>
        <Box px={4}>
          <Divider color='gray.200' />
        </Box>

        <CardBody padding={4}>
          <VStack spacing='4' w='full'>
            <Flex justifyItems='space-between' w='full'>
              <Flex direction='column'>
                <FormLabel
                  mb={2}
                  fontWeight={600}
                  color='gray.700'
                  fontSize='sm'
                >
                  Billing amounts
                  <Flex alignItems='center'>
                    <Box color='gray.500' mr={4}>
                      <CurrencyDollar height='16px' />
                    </Box>
                    <Skeleton height='12px' width='140px' borderRadius='md' mt={1}/>
                  </Flex>
                </FormLabel>
              </Flex>
              <Flex direction='column'>
                <FormLabel
                  mb={2}
                  fontWeight={600}
                  color='gray.700'
                  fontSize='sm'
                >
                  Billing frequency
                  <Flex alignItems='center'>
                    <Box color='gray.500' mr={4}>
                      <CoinsSwap height={16} />
                    </Box>
                    <Skeleton height='12px' width='140px' borderRadius='md' mt={1}/>
                  </Flex>
                </FormLabel>
              </Flex>
            </Flex>
            <Flex justifyItems='space-between' w='full'>
              <Flex direction='column'>
                <FormLabel
                  mb={2}
                  fontWeight={600}
                  color='gray.700'
                  fontSize='sm'
                >
                  Renewal cycle
                  <Flex alignItems='center'>
                    <Box mr={3} color='gray.500'>
                      <ClockCheck height={16} />
                    </Box>
                    <Skeleton height='12px' width='140px' borderRadius='md' mt={1}/>
                  </Flex>
                </FormLabel>
              </Flex>{' '}
              <Flex direction='column'>
                <FormLabel
                  mb={2}
                  fontWeight={600}
                  color='gray.700'
                  fontSize='sm'
                >
                  Renewal cycle start
                  <Flex alignItems='center'>
                    <Box mr={4} color='gray.500'>
                      <Calendar height={16} />
                    </Box>
                    <Skeleton height='12px' width='140px' borderRadius='md' mt={1}/>
                  </Flex>
                </FormLabel>
              </Flex>
            </Flex>
          </VStack>
        </CardBody>
      </Card>
    </OrganizationPanel>
  );
};
