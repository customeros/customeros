'use client';
import React, { FC, useCallback, useState } from 'react';
import { Button } from '@ui/form/Button';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
// import { FileUpload } from '@spaces/atoms/index';
import { useForm } from 'react-inverted-form';
import {
  ComposeEmailDto,
  ComposeEmailDtoI,
} from '@organization/components/Timeline/events/email/compose-email/ComposeEmail.dto';
import { SendMailRequest } from '@spaces/molecules/conversation-timeline-item/types';
import axios from 'axios';
import { useSession } from 'next-auth/react';
import { convert } from 'html-to-text';
import { Flex } from '@ui/layout/Flex';
import { ModeChangeButtons } from '@organization/components/Timeline/events/email/compose-email/EmailResponseModeChangeButtons';
import { useSearchParams } from 'next/navigation';
import { Box } from '@ui/layout/Box';
import { ParticipantsSelectGroup } from '@organization/components/Timeline/events/email/compose-email/ParticipantsSelectGroup';
import { toastSuccess } from '@ui/presentation/Toast';

interface ComposeEmail {
  subject: string;
  emailContent: string;
  to: Array<{ [x: string]: string; label: string }>;
  cc: Array<{ [x: string]: string; label: string }>;
  bcc: Array<{ [x: string]: string; label: string }>;
  from: Array<{ [x: string]: string; label: string }>;
  modal: boolean;
}

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

export const ComposeEmail: FC<ComposeEmail> = ({
  emailContent,
  subject,
  to,
  cc,
  bcc,
  from,
  modal,
}) => {
  const searchParams = useSearchParams();

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

  const [mode, setMode] = useState(REPLY_MODE);
  // const [isUploadAreaOpen, setUploadAreaOpen] = useState(false);
  // const [files, setFiles] = useState<any>([]);
  const [isSending, setIsSending] = useState(false);
  const defaultValues: ComposeEmailDtoI = new ComposeEmailDto({
    to: from,
    cc: [],
    bcc: [],
    subject: modal ? `Re: ${subject}` : '',
    content: '',
  });

  const SendMail = (
    textEmailContent: string,
    destination: Array<string> = [],
    replyTo: null | string,
    subject: null | string,
  ) => {
    const request: SendMailRequest = {
      channel: 'EMAIL',
      username: session?.user?.email || '',
      content: textEmailContent || '',
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
      .post(`/comms-api/mail/send/`, request, {
        headers: {
          'Content-Type': 'application/json',
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
        toastSuccess(
          'Something went wrong while sending the email',
          `send-email-error`,
        );
      });
  };

  const { state, handleSubmit, setDefaultValues, reset } =
    useForm<ComposeEmailDtoI>({
      formId: 'compose-email-preview',
      defaultValues,

      stateReducer: (state, action, next) => {
        return next;
      },
      onSubmit: (values, metaProps) => {
        const destination = [...values.to, ...values.cc, ...values.bcc].map(
          ({ value }) => value,
        );
        const params = new URLSearchParams(searchParams ?? '');

        setIsSending(true);
        const id = params.get('events');
        return SendMail(values.content, destination, id, values.subject);
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
          to: [...from, ...to],
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
    <Box
      borderTop={modal ? '1px dashed var(--gray-200, #EAECF0)' : 'none'}
      background={modal ? '#F8F9FC' : 'white'}
      borderRadius={modal ? 0 : 'lg'}
      borderBottomRadius='2xl'
      as='form'
      p={5}
      overflow='visible'
      maxHeight={modal ? '50vh' : 'auto'}
      pt={1}
      onSubmit={(e) => {
        e.preventDefault();
        handleSubmit(e as any);
      }}
    >
      {modal && (
        <div style={{ position: 'relative' }}>
          <ModeChangeButtons handleModeChange={handleModeChange} />
        </div>
      )}
      <ParticipantsSelectGroup
        to={state.values.to}
        cc={state.values.cc}
        bcc={state.values.bcc}
        modal={modal}
      />

      <Flex direction='column' align='flex-start' mt={2} flex={1} maxW='100%'>
        <FormAutoresizeTextarea
          placeholder='Write something here...'
          size='md'
          formId='compose-email-preview'
          name='content'
          mb={3}
          resize='none'
          borderBottom='none'
          outline='none'
          borderBottomWidth={0}
          minHeight={modal ? '100px' : '30px'}
          maxHeight={modal ? '45vh' : 'auto'}
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
            // background='white'
            type='submit'
            isDisabled={isSending}
            isLoading={isSending}
            loadingText='Sending'
          >
            Send
          </Button>
        </Flex>
        {/*{isUploadAreaOpen && (*/}
        {/*  <FileUpload*/}
        {/*    files={files}*/}
        {/*    onBeginFileUpload={(fileKey: string) => {*/}
        {/*      setFiles((prevFiles: any) => [*/}
        {/*        ...prevFiles,*/}
        {/*        {*/}
        {/*          key: fileKey,*/}
        {/*          uploaded: false,*/}
        {/*        },*/}
        {/*      ]);*/}
        {/*    }}*/}
        {/*    onFileUpload={(newFile: any) => {*/}
        {/*      setFiles((prevFiles: any) => {*/}
        {/*        return prevFiles.map((file: any) => {*/}
        {/*          if (file.key === newFile.key) {*/}
        {/*            file = {*/}
        {/*              id: newFile.id,*/}
        {/*              key: newFile.key,*/}
        {/*              name: newFile.name,*/}
        {/*              extension: newFile.extension,*/}
        {/*              uploaded: true,*/}
        {/*            };*/}
        {/*          }*/}
        {/*          return file;*/}
        {/*        });*/}
        {/*      });*/}
        {/*    }}*/}
        {/*    onFileUploadError={(fileKey: any) => {*/}
        {/*      setFiles((prevFiles: any) => {*/}
        {/*        // TODO do not remove the file from the list*/}
        {/*        // show the error instead for that particular file*/}
        {/*        return prevFiles.filter((file: any) => file.key !== fileKey);*/}
        {/*      });*/}
        {/*    }}*/}
        {/*    onFileRemove={(fileId: any) => {*/}
        {/*      setFiles((prevFiles: any) => {*/}
        {/*        return prevFiles.filter((file: any) => file.id !== fileId);*/}
        {/*      });*/}
        {/*    }}*/}
        {/*  />*/}
        {/*)}*/}
      </Flex>
    </Box>
  );
};
