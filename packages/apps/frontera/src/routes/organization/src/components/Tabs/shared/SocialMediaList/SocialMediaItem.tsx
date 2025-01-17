import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { RemoveOrgSocialMediaItemUsecase } from '@domain/usecases/organization-about-panel/remove-org-social-media-item.usecase.ts';

import { cn } from '@ui/utils/cn';
import { XCircle } from '@ui/media/icons/XCircle';
import { Copy01 } from '@ui/media/icons/Copy01.tsx';
import { Share03 } from '@ui/media/icons/Share03.tsx';
import { formatSocialUrl } from '@ui/form/UrlInput/util';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

import { SocialIcon } from './SocialIcons';

interface SocialMediaItemProps {
  id: string;
  value: string;
  index?: number;
  dataTest?: string;
  isReadOnly?: boolean;
  leftElement?: React.ReactNode;
}

const removeOrgSocialMediaItemUsecase = new RemoveOrgSocialMediaItemUsecase();
export const SocialMediaItem = observer(
  ({ value, dataTest, leftElement, id }: SocialMediaItemProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const orgId = useParams()?.id as string;
    const [_, copyToClipboard] = useCopyToClipboard();

    const href = value?.startsWith('http') ? value : `https://${value}`;
    const formattedUrl = formatSocialUrl(value);

    useEffect(() => {
      removeOrgSocialMediaItemUsecase.setId(orgId);
    }, [orgId]);

    return (
      <>
        <div className='w-full  group'>
          <div className='h-full relative w-full flex items-center'>
            <div className='h-full flex items-center '>
              <div
                tabIndex={0}
                role={'button'}
                data-test={dataTest}
                rel='noopener noreferrer'
                onClick={() => copyToClipboard(value, 'Link copied')}
                className='text-sm truncate cursor-default overflow-hidden overflow-ellipsis mr-2'
              >
                <SocialIcon url={value}>{leftElement}</SocialIcon>
                <span className='ml-3'>{formattedUrl}</span>
              </div>

              <div className='flex items-center gap-1'>
                <div>
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    aria-label={'Open in the new tab'}
                    icon={<Share03 className={'size-3 m-0'} />}
                    onClick={() =>
                      window.open(href, '_blank', 'noopener noreferrer')
                    }
                    className={cn('opacity-0 group-hover:opacity-100', {
                      '!opacity-100': isOpen,
                    })}
                  />
                </div>

                <Menu onOpenChange={(isOpen) => setIsOpen(isOpen)}>
                  <MenuButton asChild>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      aria-label={'Collapse'}
                      icon={<DotsVertical className={'size-3 m-0'} />}
                      className={cn('opacity-0 group-hover:opacity-100', {
                        '!opacity-100': isOpen,
                      })}
                    />
                  </MenuButton>

                  <MenuList className='min-w-[100px]'>
                    <MenuItem
                      onClick={() => copyToClipboard(value, 'Link copied')}
                    >
                      <Copy01 />
                      Copy link
                    </MenuItem>
                    <MenuItem
                      onClick={() => {
                        removeOrgSocialMediaItemUsecase.remove(id);
                      }}
                    >
                      <XCircle />
                      Remove social link
                    </MenuItem>
                  </MenuList>
                </Menu>
              </div>
            </div>
          </div>
        </div>
      </>
    );
  },
);
