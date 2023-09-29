'use client';
import React from 'react';
import { ClientSafeProvider, LiteralUnion, signIn } from 'next-auth/react';
import { GridItem } from '@ui/layout/Grid';
import { Flex } from '@ui/layout/Flex';
import { Box } from '@ui/layout/Box';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { Button } from '@ui/form/Button';
import { Center } from '@ui/layout/Center';
import { Image } from '@ui/media/Image';
import { Link } from '@ui/navigation/Link';
import BackgroundGridDot from '../../../public/backgrounds/grid/backgroundGridDot.png';
import { Google } from '@ui/media/logos/Google';

import CustomOsLogo from './CustomerOS-logo.png';
import Background from './login-bg.png';
import { BuiltInProviderType } from 'next-auth/providers';

export default function SignIn({
  providers,
}: {
  providers: Record<
    LiteralUnion<BuiltInProviderType, string>,
    ClientSafeProvider
  > | null;
}) {
  return (
    <>
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
            <Heading color='gray.900' size='lg' py={3} mt={-10}>
              Welcome back
            </Heading>
            <Text color='gray.500'>Sign in to your account</Text>
            {providers &&
              Object.values(providers).map((provider, i) => (
                <Button
                  mt={i === 0 ? 6 : 3}
                  key={provider.name}
                  size='md'
                  variant='outline'
                  leftIcon={<Google boxSize={6} />}
                  backgroundColor={'white'}
                  onClick={() => signIn(provider.id)}
                  width='100%'
                >
                  Sign in with {provider.name}
                </Button>
              ))}

            <Text color='gray.500' mt={2} textAlign='center' fontSize='xs'>
              By logging in you agree to CustomerOS&apos;s
              <Text color='gray.500'>
                <Link
                  color='primary.700'
                  href='https://customeros.ai/legal/terms-of-service'
                  mr={1}
                >
                  Terms of Service
                </Link>
                <Text as='span' mr={1}>
                  and
                </Text>
                <Link
                  color='primary.700'
                  href='https://www.customeros.ai/legal/privacy-policy'
                >
                  Privacy Policy
                </Link>
                .
              </Text>
            </Text>
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
    </>
  );
}
