import { FC } from 'react';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Inbox01 } from '@ui/media/icons/Inbox01';
import { FeaturedIcon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button';
import { signIn } from 'next-auth/react';
import * as Sentry from '@sentry/nextjs';
import { toastError } from '@ui/presentation/Toast';
import { Box } from '@chakra-ui/react';
import { Google } from '@ui/media/logos/Google';

export const EmptyIssueMessage: FC<{
  onAllowSendingEmail: () => void;
  modal: boolean;
}> = ({ onAllowSendingEmail, modal }) => {
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
        { callbackUrl: '/settings?tab=oauth' },
        {
          prompt: 'login',
          scope: scopes.join(' '),
        },
      );
      onAllowSendingEmail();
    } catch (error) {
      Sentry.captureException(`Unable to sign in with scopes: ${error}`);
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
      <Flex direction='column' alignItems='center' bg='white' p={6}>
        <FeaturedIcon size='md' minW='10' colorScheme='gray' mb={2}>
          <Inbox01 color='gray.700' boxSize='6' />
        </FeaturedIcon>
        <Text color='gray.700' fontWeight={600} mb={1}>
          Allow CustomerOS to send emails
        </Text>

        <Text color='gray.500' mt={1} mb={6} textAlign='center'>
          To send emails, you need to allow CustomerOS to connect to your gmail
          account
        </Text>
        <Button
          onClick={signInWithScopes}
          variant='outlined'
          colorScheme='gray'
        >
          <Google />
          Allow with google
        </Button>
      </Flex>
    </Box>
  );
};
