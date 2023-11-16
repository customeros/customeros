import React from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { RenewalLikelihoodProbability } from '@graphql/types';
import { Card, CardBody, CardHeader } from '@ui/presentation/Card';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';

// todo uncomment when BE contract with services is available
// interface RenewalARRCardProps {
//    withMultipleServices: boolean;
// }
export const RenewalARRCard = () => {
  const renewalForecast = '10000'; // remove todo when BE contract is available
  const renewalProbability = 'HIGH'; // remove todo when BE contract is available

  return (
    <Card
      px='4'
      py='3'
      w='full'
      my={2}
      size='lg'
      variant='outline'
      cursor='default'
      border='1px solid'
      borderColor='gray.200'
      position='relative'
      // sx={
      //   withMultipleServices
      //     ? {
      //         '&:after': {
      //           content: "''",
      //           width: 2,
      //           height: '80%',
      //           left: -2,
      //           top: '7px',
      //           bg: 'white',
      //           position: 'absolute',
      //           borderTopLeftRadius: 'md',
      //           borderBottomLeftRadius: 'md',
      //           border: '1px solid',
      //           borderColor: 'gray.200',
      //         },
      //       }
      //     : {}
      // }
    >
      <CardHeader as={Flex} p='0' w='full' alignItems='center' gap={4}>
        <FeaturedIcon
          size='md'
          minW='10'
          colorScheme={
            renewalForecast
              ? getFeatureIconColor(
                  renewalProbability as
                    | RenewalLikelihoodProbability
                    | undefined,
                )
              : 'gray'
          }
        >
          <ClockFastForward />
        </FeaturedIcon>
        <Flex
          alignItems='center'
          justifyContent='space-between'
          w='full'
          mt={-4}
        >
          <Flex flex={1} alignItems='center'>
            <Heading size='sm' color='gray.700' noOfLines={1}>
              Renewal ARR
            </Heading>
            <Text color='gray.500' ml={1} fontSize='sm'>
              in 1 yr
            </Text>
          </Flex>
          {/*TODO swap with real data after integrating with BE*/}
          <Text fontWeight='semibold'>$24,000</Text>
        </Flex>
      </CardHeader>

      <CardBody
        as={Text}
        w='full'
        color='gray.500'
        fontSize='sm'
        p={0}
        pl={14}
        mt={-5}
      >
        Likelihood{' '}
        <Text as='span' color='success.500'>
          High
        </Text>
      </CardBody>
    </Card>
  );
};
