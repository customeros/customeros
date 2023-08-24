'use client';
import React, { FC, useCallback, useRef, useState } from 'react';
import { Button } from '@ui/form/Button';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
// import { FileUpload } from '@spaces/atoms/index';

import { Flex } from '@ui/layout/Flex';
import { ModeChangeButtons } from '@organization/components/Timeline/events/email/compose-email/EmailResponseModeChangeButtons';
import { Box } from '@ui/layout/Box';
import { ParticipantsSelectGroup } from '@organization/components/Timeline/events/email/compose-email/ParticipantsSelectGroup';
import {
  extraAttributes,
  SocialEditor,
} from '@ui/form/RichTextEditor/SocialEditor';
import { TableExtension } from '@remirror/extension-react-tables';
import {
  AnnotationExtension,
  BlockquoteExtension,
  BoldExtension,
  BulletListExtension,
  FontSizeExtension,
  HistoryExtension,
  ImageExtension,
  ItalicExtension,
  LinkExtension,
  MentionAtomExtension,
  OrderedListExtension,
  StrikeExtension,
  TextColorExtension,
  UnderlineExtension,
  wysiwygPreset,
} from 'remirror/extensions';
import { useRemirror } from '@remirror/react';
import { NoteEditorModes } from '@spaces/organization/editor/types';

interface ComposeEmail {
  onModeChange?: (status: 'reply' | 'reply-all' | 'forward') => void;
  onSubmit: () => void;
  formId: string;
  modal: boolean;
  isSending: boolean;
  to: Array<{ label: string; value: string }>;
  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
}

export const ComposeEmail: FC<ComposeEmail> = ({
  onModeChange,
  formId,
  modal,
  isSending,
  onSubmit,
  to,
  cc,
  bcc,
}) => {
  const myRef = useRef<HTMLDivElement>(null);
  const height =
    modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 100 + 24;
  const editorRef = useRef<any | null>(null);

  const remirrorExtentions = [
    new TableExtension(),
    new MentionAtomExtension({
      matchers: [
        { name: 'at', char: '@' },
        { name: 'tag', char: '#' },
      ],
    }),

    ...wysiwygPreset(),
    new BoldExtension(),
    new ItalicExtension(),
    new BlockquoteExtension(),
    new ImageExtension({
      // uploadHandler: (e) => console.log('upload handler', e),
    }),
    new LinkExtension({ autoLink: true }),
    new TextColorExtension(),
    new UnderlineExtension(),
    new FontSizeExtension(),
    new HistoryExtension(),
    new AnnotationExtension(),
    new BulletListExtension(),
    new OrderedListExtension(),
    new StrikeExtension(),
  ];
  const extensions = useCallback(() => [...remirrorExtentions], []);

  const { manager, state, setState, getContext } = useRemirror({
    extensions,
    extraAttributes,
    // state can created from a html string.
    stringHandler: 'html',
  });

  console.log('🏷️ ----- state: '
      , state);
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
      }}
    >
      {!!onModeChange && (
        <div style={{ position: 'relative' }}>
          <ModeChangeButtons handleModeChange={onModeChange} />
        </div>
      )}
      <Box ref={myRef}>
        <ParticipantsSelectGroup
          to={to}
          cc={cc}
          bcc={bcc}
          modal={modal}
          formId={formId}
        />
      </Box>

      <Flex direction='column' align='flex-start' mt={2} flex={1} maxW='100%'>
        {/*<FormAutoresizeTextarea*/}
        {/*  placeholder='Write something here...'*/}
        {/*  size='md'*/}
        {/*  formId={formId}*/}
        {/*  name='content'*/}
        {/*  mb={3}*/}
        {/*  transform={!modal ? 'translateY(-16px)' : undefined}*/}
        {/*  resize='none'*/}
        {/*  borderBottom='none'*/}
        {/*  outline='none'*/}
        {/*  borderBottomWidth={0}*/}
        {/*  minHeight={modal ? '100px' : '30px'}*/}
        {/*  maxHeight={modal ? `calc(50vh - ${height}px) !important` : 'auto'}*/}
        {/*  height={modal ? `calc(50vh - ${height}px) !important` : 'auto'}*/}
        {/*  position='initial'*/}
        {/*  overflowY='auto'*/}
        {/*  _focusVisible={{*/}
        {/*    boxShadow: 'none',*/}
        {/*  }}*/}
        {/*/>*/}

        <SocialEditor/>
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
          {/*<Button*/}
          {/*  variant='outline'*/}
          {/*  fontWeight={600}*/}
          {/*  borderRadius='lg'*/}
          {/*  pt={0}*/}
          {/*  pb={0}*/}
          {/*  pl={3}*/}
          {/*  pr={3}*/}
          {/*  size='sm'*/}
          {/*  fontSize='sm'*/}
          {/*  isDisabled={isSending}*/}
          {/*  isLoading={isSending}*/}
          {/*  loadingText='Sending'*/}
          {/*  onClick={onSubmit}*/}
          {/*>*/}
          {/*  Send*/}
          {/*</Button>*/}
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
