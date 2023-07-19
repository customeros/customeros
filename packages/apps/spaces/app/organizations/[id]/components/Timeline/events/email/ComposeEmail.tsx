'use client';
import React, { FC, useState } from 'react';
import { Button, ButtonGroup, Flex } from '@chakra-ui/react';
import { CardFooter } from '@ui/layout/Card';
import { IconButton } from '@ui/form/IconButton';
import ReplyMany from '@spaces/atoms/icons/ReplyMany';
import Reply from '@spaces/atoms/icons/Reply';
import { EmailMetaDataEntry } from './EmailMetaDataEntry';
import { Textarea } from '@ui/form/Textarea';
import Paperclip from '@spaces/atoms/icons/Paperclip';
import { FileUpload } from '@spaces/ui-kit/atoms';
import Forward from '@spaces/atoms/icons/Forward';
import { FileTemplateUpload } from '@spaces/atoms/file-upload/FileTemplate';

interface ComposeEmail {
  subject: string;
}
const data = [
  {
    id: '1',
    key: 'key1',
    name: 'File1',
    extension: '.txt',
    uploaded: true,
  },
  {
    id: '2',
    key: 'key2',
    name: 'File2',
    extension: '.doc',
    uploaded: true,
  },
  {
    id: '3',
    key: 'key3',
    name: 'File3',
    extension: '.pdf',
    uploaded: false,
  },
];
export const ComposeEmail: FC<ComposeEmail> = ({ subject }) => {
  const [isUploadAreaOpen, setUploadAreaOpen] = useState(false);
  const [isTextAreaEditable, setIsTextAreaEditable] = useState(false);
  const [files, setFiles] = useState<any>([]);

  return (
    <CardFooter
      borderTop='1px dashed var(--gray-200, #EAECF0)'
      position='relative'
      background='#F8F9FC'
      borderBottomRadius='2xl'
      flexGrow={isUploadAreaOpen ? 2 : 1}
      onBlur={() => setIsTextAreaEditable(false)}
      onFocus={() => setIsTextAreaEditable(true)}
    >
      <ButtonGroup
        overflow='hidden'
        position='absolute'
        border='1px solid var(--gray-200, #EAECF0)'
        borderRadius={16}
        height='24px'
        gap={0}
        color='gray.300'
        background='gray.50'
        top='-4px'
        marginInlineStart={0}
      >
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          borderRadius={0}
          size='xxs'
          icon={<Reply height='16px' color='#98A2B3' />}
          pl={2}
          pr={1}
        />
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          marginInlineStart={0}
          borderRadius={0}
          size='xxs'
          icon={<ReplyMany height='14px' color='#98A2B3' />}
          pl={1}
          pr={1}
        />
        <IconButton
          variant='ghost'
          color='gray.300'
          aria-label='Call Sage'
          fontSize='14px'
          marginInline={0}
          marginInlineStart={0}
          borderRadius={0}
          size='xxs'
          icon={<Forward height='14px' color='#98A2B3' />}
          pl={1}
          pr={2}
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
          onFocus={() => setIsTextAreaEditable(true)}
          minHeight='30px'
          height={
            isTextAreaEditable
              ? isUploadAreaOpen
                ? 'calc(30vh - 130px)'
                : '30vh'
              : '30px'
          }
          _focusVisible={{
            boxShadow: 'none',
          }}
        />
        <Flex>
          {data?.length > 0 &&
            data.map((file: any, index: number) => {
              return (
                <FileTemplateUpload
                  key={`uploaded-file-${file?.name}-${file.extension}-${index}`}
                  file={file}
                  fileType={file.extension}
                  onFileRemove={() => console.log('REMOVE')}
                />
              );
            })}
        </Flex>

        <Flex
          justifyContent='flex-end'
          direction='row'
          flex={1}
          mt='lg'
          width='100%'
          pointerEvents={isTextAreaEditable ? 'all' : 'none'}
          opacity={isTextAreaEditable ? '1' : '0.5'}
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
