'use client';
import React, { FC, PropsWithChildren, useRef } from 'react';

import { ModeChangeButtons } from '@organization/src/components/Timeline/events/email/compose-email/EmailResponseModeChangeButtons';
import { Box } from '@ui/layout/Box';
import { ParticipantsSelectGroup } from '@organization/src/components/Timeline/events/email/compose-email/ParticipantsSelectGroup';
import { RichTextEditor } from '@ui/form/RichTextEditor/RichTextEditor';
import { BasicEditorToolbar } from '@ui/form/RichTextEditor/menu/BasicEditorToolbar';
import {
  BasicEditorExtentions,
  RemirrorProps,
} from '@ui/form/RichTextEditor/types';
import { KeymapperCreate } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperCreate';

interface ComposeEmail extends PropsWithChildren {
  onModeChange?: (status: 'reply' | 'reply-all' | 'forward') => void;
  onSubmit: () => void;
  formId: string;
  modal: boolean;
  isSending: boolean;
  to: Array<{ label: string; value: string }>;
  cc: Array<{ label: string; value: string }>;
  bcc: Array<{ label: string; value: string }>;
  remirrorProps: RemirrorProps<BasicEditorExtentions>;
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
  remirrorProps,
  children,
}) => {
  const myRef = useRef<HTMLDivElement>(null);
  const height =
    modal && (myRef?.current?.getBoundingClientRect()?.height || 0) + 96;

  return (
    <Box
      borderTop={modal ? '1px dashed var(--gray-200, #EAECF0)' : 'none'}
      background={modal ? '#F8F9FC' : 'white'}
      borderRadius={modal ? 0 : 'lg'}
      borderBottomRadius='2xl'
      as='form'
      p={4}
      overflow='visible'
      maxHeight={modal ? '50vh' : 'auto'}
      pt={1}
      onSubmit={(e: any) => {
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

      <Box
        maxHeight={modal ? `calc(50vh - ${height}px) !important` : 'auto'}
        w='full'
      >
        <RichTextEditor
          {...remirrorProps}
          formId={formId}
          name='content'
          showToolbar
        >
          {children}
          <KeymapperCreate onCreate={onSubmit} />
          <BasicEditorToolbar isSending={isSending} onSubmit={onSubmit} />
        </RichTextEditor>
      </Box>

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
    </Box>
  );
};
