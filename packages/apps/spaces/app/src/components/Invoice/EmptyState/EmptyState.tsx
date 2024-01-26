import { Flex } from '@ui/layout/Flex';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { File04 } from '@ui/media/icons/File04';

import HalfCirclePattern from '../../../assets/HalfCirclePattern';

export const EmptyState = ({
  maxW,
  withBorder,
}: {
  withBorder?: boolean;
  maxW?: string | number;
}) => {
  return (
    <Center h='100%' width={maxW || 500}>
      <Flex
        direction='column'
        h='100%'
        width={maxW || 500}
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
            top='22%'
            right='35%'
          >
            <File04 boxSize='5' />
          </FeaturedIcon>
          <HalfCirclePattern height={500} width={500} />
        </Flex>
        <Flex
          flexDir='column'
          textAlign='center'
          align='center'
          transform='translateY(-280px)'
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
