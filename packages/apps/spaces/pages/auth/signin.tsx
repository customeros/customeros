import type {
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next';
import { getProviders, signIn } from 'next-auth/react';
import { getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';
import { ChakraProvider, Flex } from '@chakra-ui/react';
import { theme } from '@ui/theme/theme';
import { Button } from '@ui/form/Button';
import { Center } from '@ui/layout/Center';
import { Icons } from '@ui/media/Icon';
import { Image } from '@ui/media/Image';
import { Card, CardBody } from '@ui/presentation/Card';
import CustomerOsLogo from '../../public/images/customeros-logo-dark-2.png';
import styles from '@spaces/layouts/page-content-layout/page-content-layout.module.scss';

export default function SignIn({
  providers,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <ChakraProvider theme={theme}>
      {Object.values(providers).map((provider) => (
        <Center height='100%' key={provider.name} bg='gray.100'>
          <Flex flexDirection={'column'} align={'center'}>
            <Image
              src={CustomerOsLogo}
              width={482}
              height={62}
              alt={'CustomerOS'}
            />
            <br />
            <br />
            <br />
            <br />
            <Card bg='gray.100'>
              <CardBody>
                <Button
                  leftIcon={<Icons.GOOGLE />}
                  onClick={() => signIn(provider.id)}
                >
                  Sign in with {provider.name}
                </Button>
              </CardBody>
            </Card>
          </Flex>
        </Center>
      ))}
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
