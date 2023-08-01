'use client';
import React, { FC, useCallback, useState } from 'react';
import { CardFooter } from '@ui/layout/Card';
import { Button } from '@ui/form/Button';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FileUpload } from '@spaces/atoms/index';
import { useForm } from 'react-inverted-form';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { useOutsideClick } from '@spaces/hooks/useOutsideClick';
import { EmailSubjectInput } from '@organization/components/Timeline/events/email/compose-email/EmailSubjectInput';
import { SendMailRequest } from '@spaces/molecules/conversation-timeline-item/types';
import axios from 'axios';
import { toast } from 'react-toastify';
import { useSession } from 'next-auth/react';
import { convert } from 'html-to-text';

import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { ModeChangeButtons } from '@organization/components/Timeline/events/email/compose-email/EmailResponseModeChangeButtons';
import { EmailParticipantSelect } from '@organization/components/Timeline/events/email/compose-email/EmailParticipantSelect';

interface ComposeEmail {
  subject: string;
  emailContent: string;
  to: Array<{ [x: string]: string; label: string }>;
  cc: Array<{ [x: string]: string; label: string }>;
  bcc: Array<{ [x: string]: string; label: string }>;
  from: Array<{ [x: string]: string; label: string }>;
}

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

export const ComposeEmail: FC<ComposeEmail> = ({
  emailContent,
  subject,
  cc,
  bcc,
  from,
}) => {
  const ref = React.useRef(null);
  const { data: session } = useSession();
  const text = convert(emailContent, {
    preserveNewlines: false,
    selectors: [
      {
        selector: 'a',
        options: { hideLinkHrefIfSameAsText: true, ignoreHref: true },
      },
    ],
  });
  useOutsideClick({
    ref: ref,
    handler: () => {
      setShowParticipantInputs(false);
      setShowBCC(false);
      setShowCC(false);
    },
  });

  const [mode, setMode] = useState(REPLY_MODE);
  const [isUploadAreaOpen, setUploadAreaOpen] = useState(false);
  const [showCC, setShowCC] = useState(false);
  const [showBCC, setShowBCC] = useState(false);
  const [showParticipantInputs, setShowParticipantInputs] = useState(
    !!from.length,
  );
  const [files, setFiles] = useState<any>([]);
  const [isSending, setIsSending] = useState(false);
  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: from,
    cc: [],
    bcc: [],
    subject: `Re: ${subject}`,
    content: '',
  });

  const SendMail = (
    textEmailContent: string,
    destination: Array<string> = [],
    replyTo: null | string,
    subject: null | string,
  ) => {
    if (!textEmailContent) return;
    const request: SendMailRequest = {
      channel: 'EMAIL',
      username: session?.user?.email || '',
      content: textEmailContent,
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
          reset();
          setIsSending(false);
        }
      })
      .catch((reason) => {
        setIsSending(false);
        toast.error('Something went wrong while sending the email');
      });
  };

  const { state, handleSubmit, setDefaultValues, reset } =
    useForm<ComposeEmailDtoI>({
      formId: 'compose-email-preview',
      defaultValues,

      stateReducer: (state, action, next) => {
        return next;
      },
      // @ts-expect-error fixme
      onSubmit: (values, metaProps) => {
        const destination = [...values.to, ...values.cc, ...values.bcc].map(
          ({ value }) => value,
        );
        setIsSending(true);

        return SendMail(values.content, destination, null, values.subject);
      },
    });

  const handleModeChange = useCallback(
    (newMode: string) => {
      let newDefaultValues = defaultValues;

      if (mode === newMode) {
        return;
      }
      if (newMode === REPLY_MODE) {
        newDefaultValues = new ComposeEmailDto({
          to: from,
          cc: [],
          bcc: [],
          subject: `Re: ${subject}`,
          content: mode === FORWARD_MODE ? '' : state.values.content,
        });
      }
      if (newMode === REPLY_ALL_MODE) {
        newDefaultValues = new ComposeEmailDto({
          to: from,
          cc,
          bcc,
          subject: `Re: ${subject}`,
          content: mode === FORWARD_MODE ? '' : state.values.content,
        });
      }
      if (newMode === FORWARD_MODE) {
        newDefaultValues = new ComposeEmailDto({
          to: [],
          cc: [],
          bcc: [],
          subject: `Re: ${subject}`,
          content: `${state.values.content}\n ${text}`,
        });
      }
      setMode(newMode);

      setDefaultValues(newDefaultValues);
    },
    [defaultValues, subject, state.values.content, from, cc, bcc],
  );
  return (
    <CardFooter
      borderTop='1px dashed var(--gray-200, #EAECF0)'
      background='#F8F9FC'
      borderBottomRadius='2xl'
      as='form'
      overflow='visible'
      maxHeight={'50vh'}
      pt={1}
      flexGrow={isUploadAreaOpen ? 2 : 1}
      onSubmit={(e) => {
        e.preventDefault();
        handleSubmit(e as any);
      }}
    >
      <div style={{ position: 'relative' }}>
        <ModeChangeButtons handleModeChange={handleModeChange} />
      </div>

      <Flex direction='column' align='flex-start' mt={2} flex={1} maxW='100%'>
        <Flex
          justifyContent='space-between'
          direction='row'
          flex={1}
          width='100%'
          ref={ref}
        >
          <Flex direction={'column'} flex={1} mt={2} maxWidth='90%'>
            {!showParticipantInputs && (
              <>
                <Flex
                  direction='row'
                  overflow='hidden'
                  alignItems='center'
                  alignContent='center'
                  flex={1}
                  maxWidth='100%'
                  overflowX='hidden'
                  overflowY='visible'
                  mt={1}
                  onClick={() => {
                    setShowParticipantInputs(true);
                    if (state.values.cc?.length > 0) {
                      setShowCC(true);
                    }
                    if (state.values.bcc?.length > 0) {
                      setShowBCC(true);
                    }
                  }}
                >
                  <Text
                    as={'span'}
                    color='gray.700'
                    fontWeight={600}
                    mr={1}
                    lineHeight={5}
                  >
                    To:
                  </Text>
                  <Text
                    color='gray.500'
                    overflow='hidden'
                    textOverflow='ellipsis'
                    whiteSpace='nowrap'
                  >
                    {!!state.values.to?.length && (
                      <>
                        {state.values.to
                          ?.map((email) => email.value)
                          .join(', ')}
                      </>
                    )}
                  </Text>

                  {!showCC && !!state.values.cc?.length && (
                    <>
                      <Text
                        as={'span'}
                        color='gray.700'
                        fontWeight={600}
                        ml={2}
                        mr={1}
                        lineHeight={5}
                      >
                        CC:
                      </Text>
                      <Text
                        color='gray.500'
                        overflow='hidden'
                        textOverflow='ellipsis'
                        whiteSpace='nowrap'
                      >
                        {[...state.values.cc]
                          .map((email) => email.value)
                          .join(', ')}
                      </Text>
                    </>
                  )}

                  {!showBCC && !!state.values.bcc?.length && (
                    <>
                      <Text
                        as={'span'}
                        color='gray.700'
                        fontWeight={600}
                        ml={2}
                        mr={1}
                        lineHeight={5}
                      >
                        BCC:
                      </Text>
                      <Text
                        color='gray.500'
                        overflow='hidden'
                        textOverflow='ellipsis'
                        whiteSpace='nowrap'
                      >
                        {[...state.values.bcc]
                          .map((email) => email.value)
                          .join(', ')}
                      </Text>
                    </>
                  )}
                </Flex>
              </>
            )}

            {showParticipantInputs && (
              <>
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
              </>
            )}
            <EmailSubjectInput
              formId='compose-email-preview'
              fieldName='subject'
            />
          </Flex>

          <Flex direction={'row'}>
            {!showCC && (
              <Button
                variant='ghost'
                fontWeight={600}
                color='gray.400'
                size='sm'
                px={1}
                onClick={() => {
                  setShowCC(true);
                  setShowParticipantInputs(true);
                }}
              >
                CC
              </Button>
            )}

            {!showBCC && (
              <Button
                variant='ghost'
                fontWeight={600}
                size='sm'
                px={1}
                color='gray.400'
                onClick={() => {
                  setShowBCC(true);
                  setShowParticipantInputs(true);
                }}
              >
                BCC
              </Button>
            )}
          </Flex>
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
          minHeight='100px'
          maxHeight={
            showBCC || showCC ? `calc(50vh - 16rem)` : `calc(50vh - 12rem)`
          }
          position='initial'
          overflowY='auto'
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
          {/*<IconButton*/}
          {/*  size='sm'*/}
          {/*  mr={2}*/}
          {/*  borderRadius='lg'*/}
          {/*  variant='ghost'*/}
          {/*  aria-label='Add attachement'*/}
          {/*  onClick={() => {*/}
          {/*    setUploadAreaOpen(!isUploadAreaOpen);*/}
          {/*  }}*/}
          {/*  isDisabled*/}
          {/*  icon={<Paperclip color='gray.400' height='20px' />}*/}
          {/*/>*/}
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
            background='white'
            type='submit'
            isDisabled={isSending}
            isLoading={isSending}
            loadingText='Sending'
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
