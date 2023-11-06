import { FC } from 'react';

import { signIn } from 'next-auth/react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { FeaturedIcon } from '@ui/media/Icon';
import { Mail01 } from '@ui/media/icons/Mail01';
import { Google } from '@ui/media/logos/Google';
import { toastError } from '@ui/presentation/Toast';

export const MissingPermissionsPrompt: FC<{
  modal: boolean;
}> = ({ modal }) => {
  const signInWithScopes = async () => {
    const scopes = [
      'openid',
      'email',
      'profile',
      'https://www.googleapis.com/auth/gmail.readonly',
      'https://www.googleapis.com/auth/gmail.send',
      'https://www.googleapis.com/auth/calendar.readonly',
    ];

    try {
      await signIn(
        'google',
        { callbackUrl: window.location.href },
        {
          prompt: 'login',
          scope: scopes.join(' '),
        },
      );
    } catch (error) {
      toastError('Something went wrong!', `unable-to-sign-in-with-scopes`);
    }
  };

  return (
    <Box
      alignItems='center'
      mt='4'
      borderTop={modal ? '1px dashed var(--gray-200, #EAECF0)' : 'none'}
      background={modal ? '#F8F9FC' : 'white'}
      borderRadius={modal ? 0 : 'lg'}
      borderBottomRadius='2xl'
      as='form'
      p={6}
      overflow='visible'
      maxHeight={modal ? '50vh' : 'auto'}
    >
      <Flex
        direction='column'
        alignItems='center'
        bg={modal ? '#F8F9FC' : 'white'}
        p={6}
      >
        <FeaturedIcon size='md' minW='10' colorScheme='gray' mb={4}>
          <Mail01 color='gray.700' boxSize='6' />
        </FeaturedIcon>
        <Text color='gray.700' fontWeight={600} mb={1}>
          Allow CustomerOS to send emails
        </Text>

        <Text color='gray.500' mb={6} textAlign='center'>
          To send emails, you need to allow CustomerOS to connect to your gmail
          account
        </Text>
        <Button variant='outline' colorScheme='gray' onClick={signInWithScopes}>
          <Google mr={2} />
          Allow with google
        </Button>
      </Flex>
    </Box>
  );
};
