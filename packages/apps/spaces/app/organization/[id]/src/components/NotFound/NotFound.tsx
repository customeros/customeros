'use client';
import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Heading } from '@ui/typography/Heading';
import { SearchSm } from '@ui/media/icons/SearchSm';
import HalfCirclePattern from '../../../../../src/assets/HalfCirclePattern';

export default function NotFound() {
  return (
    <Box
      p={0}
      flex={1}
      as={Flex}
      flexDirection='column'
      backgroundRepeat='no-repeat'
      backgroundSize='contain'
      w='100vw'
      position='relative'
      alignItems='center'
      justifyContent='center'
      backgroundColor='gray.25'
      border='1px solid'
      borderColor='gray.200'
      borderRadius='xl'
    >
      <Box
        position='absolute'
        height='50vh'
        maxH='768px'
        width='768px'
        top='50%'
        left='50%'
        style={{
          transform: 'translate(-50%, -90%) rotate(180deg)',
        }}
      >
        <HalfCirclePattern />
      </Box>
      <Flex
        position='relative'
        direction='column'
        alignItems='center'
        justifyContent='center'
        h='50vh'
      >
        <FeaturedIcon colorScheme='primary' size='lg'>
          <SearchSm boxSize='5' />
        </FeaturedIcon>
        <Heading fontWeight={600} fontSize='5xl' color='gray.900' py={6}>
          This organization cannot be found
        </Heading>
        <Text color='gray.600' fontSize='2xl' pb={12} px={8} textAlign='center'>
          It appears the organization does not exist or you do not have
          sufficient rights to preview it.
        </Text>
      </Flex>
    </Box>
  );
}
