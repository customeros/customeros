import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { RemoveOrgSocialMediaItemUsecase } from '@domain/usecases/organization-about-panel/remove-org-social-media-item.usecase.ts';

import { cn } from '@ui/utils/cn';
import { XCircle } from '@ui/media/icons/XCircle';
import { formatSocialUrl } from '@ui/form/UrlInput/util';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
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
              <a
                href={href}
                target={'_blank'}
                data-test={dataTest}
                rel='noopener noreferrer'
                className='text-sm cursor-pointer overflow-hidden overflow-ellipsis mr-2'
              >
                <SocialIcon url={value}>{leftElement}</SocialIcon>
                <span className='ml-3'>{formattedUrl}</span>
              </a>
              <Menu onOpenChange={(isOpen) => setIsOpen(isOpen)}>
                <MenuButton>
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
      </>
    );
  },
);
