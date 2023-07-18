'use client';
import React, { FC, useState } from 'react';
import { Text } from '@ui/typography/Text';
import { Button, ButtonGroup, Flex } from '@chakra-ui/react';
import { InteractionEventParticipant } from '@graphql/types';
import { CardFooter } from '@ui/layout/Card';
import { IconButton } from '@ui/form/IconButton';
import ReplyLeft from '@spaces/atoms/icons/ReplyLeft';
import ReplyMany from '@spaces/atoms/icons/ReplyMany';
import Reply from '@spaces/atoms/icons/Reply';
import { EmailMetaDataEntry } from './EmailMetaDataEntry';
import { Textarea } from '@ui/form/Textarea';
import Paperclip from '@spaces/atoms/icons/Paperclip';
import { FileUpload } from '@spaces/ui-kit/atoms';

interface ComposeEmail {
  subject: string;
}

export const ComposeEmail: FC<ComposeEmail> = ({ subject }) => {
  const [isUploadAreaOpen, setUploadAreaOpen] = useState(false);
  return (
    <CardFooter
      borderTop='1px dashed var(--gray-200, #EAECF0)'
      position='relative'
      background='#F8F9FC'
      borderBottomRadius='2xl'
      flex={isUploadAreaOpen ? 2 : 1}
    >
      <ButtonGroup
        overflow='hidden'
        position='absolute'
        border='1px solid var(--gray-200, #EAECF0)'
        borderRadius={16}
        gap={0}
        color='gray.300'
        background='gray.50'
        top='-4px'
      >
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          size='xs'
          borderRadius={0}
          icon={<ReplyLeft height='10px' color='#98A2B3' />}
        />
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          size='xs'
          marginInlineStart={0}
          borderRadius={0}
          icon={<ReplyMany height='10px' color='#98A2B3' />}
        />
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          marginInlineStart={0}
          borderRadius={0}
          size='xs'
          icon={<Reply height='10px' color='#98A2B3' />}
        />
      </ButtonGroup>

      <Flex direction='column' align='flex-start' mt={2} flex={1}>
        <EmailMetaDataEntry entryType='To' content='test@test.com' />
        <EmailMetaDataEntry entryType='CC' content='test@test.com' />
        <EmailMetaDataEntry entryType='Subject' content={subject} />

        <Textarea
          placeholder='Write something here...'
          size='md'
          mt={1}
          mb={3}
          resize='none'
          borderBottom='none'
          outline='none'
          borderBottomWidth={0}
          minHeight='30px'
          _focusVisible={{ boxShadow: 'none', flex: 2 }}
        />
        <Flex
          justifyContent='flex-end'
          direction='row'
          flex={1}
          mt='lg'
          width='100%'
        >
          <IconButton
            size='sm'
            mr={2}
            borderRadius='lg'
            variant='ghost'
            aria-label='Add attachement'
            onClick={() => setUploadAreaOpen(true)}
            icon={<Paperclip color='#98A2B3' height='20px' />}
          />
          <Button
            variant='outline'
            fontWeight={600}
            borderRadius='lg'
            pt={0}
            pb={0}
            pl={3}
            pr={3}
            size='sm'
            fontSize='sm'
            background='#fff'
          >
            Send
          </Button>
        </Flex>
        {isUploadAreaOpen && <FileUpload />}
      </Flex>
    </CardFooter>
  );
};
