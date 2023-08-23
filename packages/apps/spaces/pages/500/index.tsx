import React from 'react';
import { Box, Flex, Heading, Text } from '@chakra-ui/react';
import HalfCirclePattern from '@spaces/atoms/icons/HalfCirclePattern';
import { FeaturedIcon } from '@ui/media/Icon';
import Search from '@spaces/atoms/icons/Search';
import { Button } from '@ui/form/Button';
import { useRouter } from 'next/router';

export const ServerError: React.FC = () => {
  const { push } = useRouter();

  return (
    <Box
      p={0}
      flex={1}
      as={Flex}
      flexDirection='column'
      bgImage='/backgrounds/organization/half-circle-pattern.svg'
      backgroundRepeat='no-repeat'
      backgroundSize='contain'
      h='100vh'
      w='100vw'
      position='relative'
      alignItems='center'
      justifyContent='center'
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
          <Box>
            <Search />
          </Box>
        </FeaturedIcon>
        <Heading fontWeight={600} fontSize='6xl' color='gray.900' py={6}>
          We are sorry, but something went wrong on our end.
        </Heading>
        <Text color='gray.600' fontSize='xl' pb={12} px={8}>
          Our team has been notified of the issue and is working to fix it as
          soon as possible.
        </Text>
        <Button
          colorScheme='primary'
          variant='outline'
          size='lg'
          onClick={() => push('/')}
        >
          Take me home
        </Button>
      </Flex>
    </Box>
  );
};

export default ServerError;
