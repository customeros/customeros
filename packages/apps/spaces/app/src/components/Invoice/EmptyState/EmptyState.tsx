import { Flex } from '@ui/layout/Flex';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { File04 } from '@ui/media/icons/File04';

import HalfCirclePattern from '../../../assets/HalfCirclePattern';

export const EmptyState = ({
  maxW = 500,
  withBorder,
}: {
  withBorder?: boolean;
  maxW?: string | number;
}) => {
  return (
    <Center h='100%' width={maxW}>
      <Flex
        direction='column'
        h='100%'
        width={maxW}
        borderRight={withBorder ? '1px solid' : 'none'}
        borderColor='gray.200'
      >
        <Flex position='relative'>
          <FeaturedIcon
            colorScheme='primary'
            size='lg'
            width='152px'
            height='120'
            position='absolute'
            top={withBorder ? '22%' : '20%'}
            right={withBorder ? '35%' : '33%'}
          >
            <File04 boxSize='5' />
          </FeaturedIcon>
          <HalfCirclePattern height={maxW} width={maxW} />
        </Flex>
        <Flex
          flexDir='column'
          textAlign='center'
          align='center'
          transform={withBorder ? 'translateY(-280px)' : 'translateY(-250px)'}
        >
          <Text color='gray.900' fontSize='md' fontWeight='semibold'>
            Awaiting your invoices
          </Text>
          <Text maxW='350px' fontSize='sm' color='gray.600' my={1}>
            Create your first contract with services, and your invoices will
            appear here.
          </Text>
        </Flex>
      </Flex>
    </Center>
  );
};
