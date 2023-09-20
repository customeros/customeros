import React from 'react';

import type {
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next';
import { getProviders, signIn } from 'next-auth/react';
import { getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';
import {
  Box,
  ChakraProvider,
  Flex,
  Grid,
  GridItem,
  Heading,
  Text,
} from '@chakra-ui/react';
import { theme } from '@ui/theme/theme';
import { Button } from '@ui/form/Button';
import { Center } from '@ui/layout/Center';
import { Image } from '@ui/media/Image';
import BackgroundGridDot from '../../public/backgrounds/grid/backgroundGridDot.png';
import GoogleLogo from '@spaces/atoms/icons/GoogleLogo';

import CustomOsLogo from './CustomerOS-logo.png';
import Background from './login-bg.png';

export default function SignIn({
  providers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <ChakraProvider theme={theme}>
      <Grid templateColumns={{ base: '1fr', md: '1fr 1fr' }} h='100vh'>
        <GridItem h='100vh'>
          <Box height='50%'>
            <Image
              alt=''
              src={BackgroundGridDot}
              width={480}
              top='-10%'
              margin='auto'
            />
          </Box>
          <Center height='100%' pos='relative' top='-50%'>
            <Flex flexDirection={'column'} align={'center'} width={360}>
              <Image
                src={CustomOsLogo}
                alt='CustomerOS'
                width={264}
                height={264}
              />
              <Heading color='gray.900' size='lg' py={3}>
                Welcome back
              </Heading>
              <Text color='gray.500'>Sign in to your account</Text>
              {Object.values(providers).map((provider, i) => (
                <Button
                  mt={i === 0 ? 6 : 3}
                  key={provider.name}
                  size='md'
                  variant='outline'
                  leftIcon={<GoogleLogo height={24} width={24} />}
                  backgroundColor={'white'}
                  onClick={() => signIn(provider.id)}
                  width='100%'
                >
                  Sign in with {provider.name}
                </Button>
              ))}
            </Flex>
          </Center>
        </GridItem>
        <GridItem
          borderTopLeftRadius='80px'
          borderBottomLeftRadius='80px'
          bg={`url(${Background.src})`}
          backgroundRepeat={'no-repeat'}
          backgroundSize={'cover'}
        ></GridItem>
      </Grid>
    </ChakraProvider>
  );
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getServerSession(context.req, context.res, authOptions);

  // If the user is already logged in, redirect.
  // Note: Make sure not to redirect to the same page
  // To avoid an infinite loop!
  if (session) {
    return { redirect: { destination: '/' } };
  }

  const providers = await getProviders();

  return {
    props: { providers: providers ?? [] },
  };
}
