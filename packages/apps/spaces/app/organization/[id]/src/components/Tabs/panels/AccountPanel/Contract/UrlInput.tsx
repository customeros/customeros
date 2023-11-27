import { useForm } from 'react-inverted-form';
import React, { FC, useState, ChangeEvent } from 'react';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { FormInput } from '@ui/form/Input';
import { Check } from '@ui/media/icons/Check';
import { Link01 } from '@ui/media/icons/Link01';
import { File02 } from '@ui/media/icons/File02';
import { IconButton } from '@ui/form/IconButton';
import { Trash01 } from '@ui/media/icons/Trash01';
import { Divider } from '@ui/presentation/Divider';
import { ContractUpdateInput } from '@graphql/types';
import { FileCheck02 } from '@ui/media/icons/FileCheck02';
import { getExternalUrl } from '@spaces/utils/getExternalLink';
import {
  Popover,
  PopoverArrow,
  PopoverContent,
  PopoverTrigger,
} from '@ui/overlay/Popover';
import {
  ContractDTO,
  TimeToRenewalForm,
} from '@organization/src/components/Tabs/panels/AccountPanel/Contract/Contract.dto';

interface UrlInputProps {
  formId: string;
  url?: string | null;
  contractId?: string | null;
  onSubmit: (props: { input: ContractUpdateInput }) => void;
}
export const UrlInput: FC<UrlInputProps> = ({
  url,
  onSubmit,
  contractId,
  formId,
}) => {
  const { onOpen, onClose, isOpen } = useDisclosure();
  const firstFieldRef = React.useRef(null);
  const [href, setHref] = useState<string>(url ?? '');
  const { state } = useForm<TimeToRenewalForm>({
    formId,
  });

  const handleUpdateContractUrl = (value: string) => {
    if (!contractId) return;

    onSubmit({
      input: {
        contractId,
        ...ContractDTO.toPayload({
          ...state.values,
          contractUrl: value,
        }),
      },
    });
    onClose();
  };

  return (
    <>
      <Popover
        isOpen={isOpen}
        initialFocusRef={firstFieldRef}
        onOpen={onOpen}
        onClose={onClose}
        placement='bottom-start'
        closeOnBlur
      >
        <PopoverTrigger>
          <IconButton
            borderRadius='sm'
            size='xs'
            variant='ghost'
            color='gray.400'
            aria-label='Click to add contract link'
            icon={url ? <FileCheck02 /> : <File02 />}
            onClick={onOpen}
          />
        </PopoverTrigger>
        <PopoverContent py={1} px={3} background='gray.700' borderRadius='8px'>
          <PopoverArrow borderColor='gray.700' bg='gray.700' />
          <Flex alignItems='center'>
            <IconButton
              size='xs'
              variant='ghost'
              aria-label='Go to url'
              disabled={!href}
              onClick={() => {
                window.open(
                  getExternalUrl(href),
                  '_blank',
                  'noopener noreferrer',
                );
              }}
              icon={<Link01 color='gray.25' />}
              mr={2}
              borderRadius='sm'
              _hover={{ background: 'gray.600', color: 'gray.25' }}
            />
            <FormInput
              label='Contract link'
              formId={formId}
              name='contractUrl'
              background='gray.700'
              fontSize='sm'
              color='gray.25'
              whiteSpace='nowrap'
              overflow='hidden'
              textOverflow='ellipsis'
              _placeholder={{
                color: 'gray.400',
              }}
              _focusVisible={{
                outline: 'none',
              }}
              tabIndex={1}
              border='none'
              placeholder='Paste or enter a contract link'
              onChange={(event: ChangeEvent<HTMLInputElement>) =>
                setHref(event.target.value)
              }
              value={href}
              onKeyDown={(event) => {
                const { key } = event;
                if (key === 'Enter') {
                  handleUpdateContractUrl(href);
                }

                if (key === 'Escape') {
                  handleUpdateContractUrl('');
                }
              }}
            />
            {href && (
              <>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Save'
                  onClick={() => handleUpdateContractUrl(href)}
                  color='gray.400'
                  icon={<Check color='inherit' />}
                  mr={2}
                  ml={2}
                  borderRadius='sm'
                  _hover={{ background: 'gray.600', color: 'gray.25' }}
                />

                <Divider
                  orientation='vertical'
                  borderLeft='1px solid'
                  borderLeftColor='gray.400 !important'
                  height='14px'
                />

                <IconButton
                  ml={2}
                  borderRadius='sm'
                  size='xs'
                  variant='ghost'
                  aria-label='Remove link'
                  onClick={() => handleUpdateContractUrl('')}
                  color='gray.400'
                  icon={<Trash01 color='inherit' />}
                  _hover={{ background: 'gray.600', color: 'gray.25' }}
                />
              </>
            )}
          </Flex>
        </PopoverContent>
      </Popover>
    </>
  );
};
