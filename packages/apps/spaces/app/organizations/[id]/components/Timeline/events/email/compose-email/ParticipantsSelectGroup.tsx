import React, { useEffect, useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { InputProps } from '@ui/form/Input';

import { EmailParticipantSelect } from '@organization/components/Timeline/events/email/compose-email/EmailParticipantSelect';
import { useOutsideClick } from '@spaces/hooks/useOutsideClick';
import { Box } from '@chakra-ui/react';
import { EmailSubjectInput } from '@organization/components/Timeline/events/email/compose-email/EmailSubjectInput';
import { Button } from '@ui/form/Button';
import Image from 'next/image';

interface ParticipantSelectGroupGroupProps extends InputProps {
  to: Array<{ label: string; value: string }>;
  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;

  modal?: boolean;
  formId: string;
}

export const ParticipantsSelectGroup = ({
  to = [],
  cc = [],
  bcc = [],
  modal,
  formId,
}: ParticipantSelectGroupGroupProps) => {
  const [showCC, setShowCC] = useState(false);
  const [showBCC, setShowBCC] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [focusedItemIndex, setFocusedItemIndex] = useState<false | number>(
    false,
  );
  const ref = React.useRef(null);
  useOutsideClick({
    ref: ref,
    handler: () => {
      setIsFocused(false);
      setFocusedItemIndex(false);
      setShowCC(false);
      setShowBCC(false);
    },
  });

  const handleFocus = (index: number) => {
    setIsFocused(true);
    setFocusedItemIndex(index);
  };

  useEffect(() => {
    if (showCC && !isFocused) {
      handleFocus(1);
    }
  }, [showCC]);

  useEffect(() => {
    if (showBCC && !isFocused) {
      handleFocus(2);
    }
  }, [showBCC]);

  return (
    <Flex justifyContent='space-between' mt={3} ref={ref}>
      <Box width='100%'>
        {isFocused && (
          <>
            <EmailParticipantSelect
              formId={formId}
              fieldName='to'
              entryType='To'
              autofocus={focusedItemIndex === 0}
            />
            {(showCC || !!cc.length) && (
              <EmailParticipantSelect
                formId={formId}
                fieldName='cc'
                entryType='CC'
                autofocus={focusedItemIndex === 1}
              />
            )}
            {(showBCC || !!bcc.length) && (
              <EmailParticipantSelect
                formId={formId}
                fieldName='bcc'
                entryType='BCC'
                autofocus={focusedItemIndex === 2}
              />
            )}
          </>
        )}

        {!isFocused && (
          <Flex mt={1} flex={isFocused ? 1 : 'unset'}>
            <Flex
              onClick={() => handleFocus(0)}
              role='button'
              aria-label='Click to input participant data'
              flex={!to.length ? 1 : 'unset'}
            >
              <Text as={'span'} color='gray.700' fontWeight={600} mr={1}>
                To:
              </Text>
              <Text color='gray.500' noOfLines={1}>
                {!!to?.length && (
                  <>{to?.map((email) => email.value).join(', ')}</>
                )}
              </Text>
            </Flex>

            {!!cc.length && (
              <Flex
                onClick={() => handleFocus(1)}
                role='button'
                aria-label='Click to input participant data'
                flex={!bcc.length ? 1 : 'unset'}
              >
                <Text
                  as={'span'}
                  color='gray.700'
                  fontWeight={600}
                  ml={2}
                  mr={1}
                >
                  CC:
                </Text>
                <Text color='gray.500' noOfLines={1}>
                  {[...cc].map((email) => email.value).join(', ')}
                </Text>
              </Flex>
            )}
            {!!bcc.length && (
              <Flex
                onClick={() => handleFocus(2)}
                role='button'
                aria-label='Click to input participant data'
              >
                <Text
                  as={'span'}
                  color='gray.700'
                  fontWeight={600}
                  ml={2}
                  mr={1}
                >
                  BCC:
                </Text>
                <Text color='gray.500' noOfLines={1}>
                  {[...bcc].map((email) => email.value).join(', ')}
                </Text>
              </Flex>
            )}
          </Flex>
        )}
        <EmailSubjectInput formId={formId} fieldName='subject' />
      </Box>
      <Flex maxW='64px'>
        {!showCC && (
          <Button
            variant='ghost'
            fontWeight={600}
            color='gray.400'
            size='sm'
            px={1}
            onClick={() => {
              setShowCC(true);
              setFocusedItemIndex(1);
            }}
          >
            CC
          </Button>
        )}

        {/*{!showBCC && (*/}
        {/*  <Button*/}
        {/*    variant='ghost'*/}
        {/*    fontWeight={600}*/}
        {/*    size='sm'*/}
        {/*    px={1}*/}
        {/*    color='gray.400'*/}
        {/*    onClick={() => {*/}
        {/*      setShowBCC(true);*/}
        {/*      setFocusedItemIndex(2);*/}
        {/*    }}*/}
        {/*  >*/}
        {/*    BCC*/}
        {/*  </Button>*/}
        {/*)}*/}
      </Flex>

      {!modal && (
        <Box position='relative'>
          <Image
            src={'/backgrounds/organization/post-stamp.webp'}
            alt='Email'
            width={54}
            height={70}
            style={{
              filter: 'drop-shadow(0px 0.5px 1px #D8D8D8)',
              marginLeft: '8px',
            }}
          />
        </Box>
      )}
    </Flex>
  );
};
