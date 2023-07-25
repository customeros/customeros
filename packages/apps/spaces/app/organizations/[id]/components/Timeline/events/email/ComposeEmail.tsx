'use client';
import React, { FC, useState } from 'react';
import { Button, ButtonGroup, Flex } from '@chakra-ui/react';
import { CardFooter } from '@ui/layout/Card';
import { IconButton } from '@ui/form/IconButton';
import ReplyMany from '@spaces/atoms/icons/ReplyMany';
import Reply from '@spaces/atoms/icons/Reply';
import { EmailParticipantSelect } from './EmailParticipantSelect';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import Paperclip from '@spaces/atoms/icons/Paperclip';
import { FileUpload } from '@spaces/ui-kit/atoms';
import Forward from '@spaces/atoms/icons/Forward';
// import { FileTemplateUpload } from '@spaces/atoms/file-upload/FileTemplate';
import { useForm } from 'react-inverted-form';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/events/email/ComposeEmail.dto';
import { useOutsideClick } from '@chakra-ui/react-use-outside-click';
import { EmailSubjectInput } from '@organization/components/Timeline/events/email/EmailSubjectInput';
import { SendMailRequest } from '@spaces/molecules/conversation-timeline-item/types';
import axios from 'axios';
import { toast } from 'react-toastify';
import { useRecoilValue } from 'recoil';
import { userData } from '@spaces/globalState/userData';

interface ComposeEmail {
  subject: string;
}

export const ComposeEmail: FC<ComposeEmail> = ({ subject }) => {
  const ref = React.useRef(null);
  const loggedInUserData = useRecoilValue(userData);

  useOutsideClick({
    ref: ref,
    handler: () => {
      setShowBCC(false);
      setShowCC(false);
    },
  });

  const [isUploadAreaOpen, setUploadAreaOpen] = useState(false);
  const [isTextAreaEditable, setIsTextAreaEditable] = useState(false);
  const [showCC, setShowCC] = useState(false);
  const [showBCC, setShowBCC] = useState(false);
  const [files, setFiles] = useState<any>([]);
  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: [],
    cc: [],
    bcc: [],
    subject: `Re: ${subject}`,
    content: '',
    files: [],
  });

  const SendMail = (
    text: string,
    onSuccess: () => void,
    destination: Array<string> = [],
    replyTo: null | string,
    subject: null | string,
  ) => {
    if (!text) return;
    const request: SendMailRequest = {
      channel: 'EMAIL',
      username: loggedInUserData.identity,
      content: text,
      direction: 'OUTBOUND',
      destination: destination,
    };
    if (replyTo) {
      request.replyTo = replyTo;
    }
    if (subject) {
      request.subject = subject;
    }
    return axios
      .post(`/comms-api/mail/send`, request, {
        headers: {
          'X-Openline-Mail-Api-Key': `${process.env.COMMS_MAIL_API_KEY}`,
        },
      })
      .then((res) => {
        if (res.data) {
          onSuccess();
        }
      })
      .catch((reason) => {
        console.log('🏷️ ----- : '
            , reason);
        toast.error('Something went wrong while sending request');
      });
  };

  const { state, handleSubmit } = useForm<ComposeEmailDtoI>({
    formId: 'compose-email-preview',
    defaultValues,

    stateReducer: (state, action, next) => {
      return next;
    },
    onSubmit: (values, metaProps) => {
      const destination = [...values.to, ...values.cc, ...values.bcc].map(
        ({ value }) => value,
      );

      return SendMail(
        values.content,
        () => {
          console.log('send');
        },
        destination,
        null,
        values.subject,
      );

    },
  });

  return (
    <CardFooter
      borderTop='1px dashed var(--gray-200, #EAECF0)'
      position='relative'
      background='#F8F9FC'
      borderBottomRadius='2xl'
      as='form'
      pt={1}
      flexGrow={isUploadAreaOpen ? 2 : 1}
      onBlur={() => setIsTextAreaEditable(false)}
      onFocus={() => setIsTextAreaEditable(true)}
      onSubmit={(e) => {
        e.preventDefault();
        handleSubmit(e);
      }}
    >
      <ButtonGroup
        overflow='hidden'
        position='absolute'
        border='1px solid var(--gray-200, #EAECF0)'
        borderRadius={16}
        height='24px'
        gap={0}
        color='#FCFCFD'
        background='#FCFCFD'
        top='-14px'
      >
        <IconButton
          variant='ghost'
          aria-label='Call Sage'
          fontSize='14px'
          color='gray.400'
          borderRadius={0}
          marginInlineStart={0}

          size='xxs'
          icon={<Reply height='16px' color='gray.400' />}
          pl={2}
          pr={1}
        />
        <IconButton
          variant='ghost'
          aria-label='Call Sage'
          fontSize='14px'
          color='gray.400'
          marginInlineStart={0}
          borderRadius={0}
          size='xxs'
          icon={<ReplyMany height='14px' color='gray.400' />}
          pl={1}
          pr={1}
        />
        <IconButton
          variant='ghost'
          aria-label='Call Sage'
          fontSize='14px'
          color='gray.400'
          marginInline={0}
          marginInlineStart={0}
          borderRadius={0}
          size='xxs'
          icon={<Forward height='14px' color='gray.400' />}
          pl={1}
          pr={2}
        />
      </ButtonGroup>

      <Flex direction='column' align='flex-start' mt={2} flex={1}>
        <Flex
          justifyContent='space-between'
          direction='row'
          flex={1}
          width='100%'
          ref={ref}
        >
          <Flex direction='column' flex={1}>
            <EmailParticipantSelect
              formId='compose-email-preview'
              fieldName='to'
              entryType='To'
            />

            {showCC && (
              <EmailParticipantSelect
                formId='compose-email-preview'
                fieldName='cc'
                entryType='CC'
              />
            )}
            {showBCC && (
              <EmailParticipantSelect
                formId='compose-email-preview'
                fieldName='Bcc'
                entryType='BCC'
              />
            )}
            <EmailSubjectInput
              formId='compose-email-preview'
              fieldName='subject'
            />
          </Flex>
          <div>
            {!showCC && (
              <Button
                variant='ghost'
                fontWeight={600}
                color='gray.400'
                size='sm'
                onClick={() => setShowCC(true)}
              >
                CC
              </Button>
            )}

            {!showBCC && (
              <Button
                variant='ghost'
                fontWeight={600}
                color='gray.400'
                size='sm'
                onClick={() => setShowBCC(true)}
              >
                BCC
              </Button>
            )}
          </div>
        </Flex>

        <FormAutoresizeTextarea
          placeholder='Write something here...'
          size='md'
          mt={1}
          formId='compose-email-preview'
          name='content'
          mb={3}
          resize='none'
          borderBottom='none'
          outline='none'
          borderBottomWidth={0}
          onFocus={() => setIsTextAreaEditable(true)}
          minHeight='30px'
          overflowY='auto'
          maxHeight={'20vh'}
          _focusVisible={{
            boxShadow: 'none',
          }}
        />
        {/*<Flex>*/}
        {/*  {data?.length > 0 &&*/}
        {/*    data.map((file: any, index: number) => {*/}
        {/*      return (*/}
        {/*        <FileTemplateUpload*/}
        {/*          key={`uploaded-file-${file?.name}-${file.extension}-${index}`}*/}
        {/*          file={file}*/}
        {/*          fileType={file.extension}*/}
        {/*          onFileRemove={() => console.log('REMOVE')}*/}
        {/*        />*/}
        {/*      );*/}
        {/*    })}*/}
        {/*</Flex>*/}

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
            onClick={() => {
              setUploadAreaOpen(!isUploadAreaOpen);
            }}
            isDisabled
            icon={<Paperclip color='gray.400' height='20px' />}
          />
          <Button
            pointerEvents={isTextAreaEditable ? 'all' : 'none'}
            opacity={isTextAreaEditable ? '1' : '0.5'}
            variant='outline'
            fontWeight={600}
            borderRadius='lg'
            pt={0}
            pb={0}
            pl={3}
            pr={3}
            size='sm'
            fontSize='sm'
            background='white'
            type='submit'
          >
            Send
          </Button>
        </Flex>
        {isUploadAreaOpen && (
          <FileUpload
            files={files}
            onBeginFileUpload={(fileKey: string) => {
              setFiles((prevFiles: any) => [
                ...prevFiles,
                {
                  key: fileKey,
                  uploaded: false,
                },
              ]);
            }}
            onFileUpload={(newFile: any) => {
              setFiles((prevFiles: any) => {
                return prevFiles.map((file: any) => {
                  if (file.key === newFile.key) {
                    file = {
                      id: newFile.id,
                      key: newFile.key,
                      name: newFile.name,
                      extension: newFile.extension,
                      uploaded: true,
                    };
                  }
                  return file;
                });
              });
            }}
            onFileUploadError={(fileKey: any) => {
              setFiles((prevFiles: any) => {
                // TODO do not remove the file from the list
                // show the error instead for that particular file
                return prevFiles.filter((file: any) => file.key !== fileKey);
              });
            }}
            onFileRemove={(fileId: any) => {
              setFiles((prevFiles: any) => {
                return prevFiles.filter((file: any) => file.id !== fileId);
              });
            }}
          />
        )}
      </Flex>
    </CardFooter>
  );
};
