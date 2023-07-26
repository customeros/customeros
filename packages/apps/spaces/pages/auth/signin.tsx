import type {
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next';
import { getProviders, signIn } from 'next-auth/react';
import { getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';
import {
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
import LoginBg from '@spaces/atoms/backgrounds/Login';
import GoogleLogo from '@spaces/atoms/icons/GoogleLogo';
import React from 'react';
import CustomerOsLogo from '@spaces/atoms/icons/CustomerOsLogo';

export default function SignIn({
  providers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <ChakraProvider theme={theme}>
      <Grid templateColumns={['1fr', '1fr', '1fr', '10fr 11fr']} gap={1} h='100vh'>
        <GridItem h='100vh'>
          <Center height='100%'>
            <Flex flexDirection={'column'} align={'center'} width={360}>
              <CustomerOsLogo height={64} />
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
          h='100vh'
          w='50vw'
          >
            <LoginBg />   
        </GridItem>
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
