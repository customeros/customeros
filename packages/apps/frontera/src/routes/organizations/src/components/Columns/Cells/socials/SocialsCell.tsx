import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { Social } from '@shared/types/__generated__/graphql.types';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02.tsx';
import {
  getExternalUrl,
  getFormattedLink,
} from '@spaces/utils/getExternalLink';
import { isKnownUrl } from '@organization/components/Tabs/shared/FormSocialInput/util.ts';
import { SocialIcon } from '@organization/components/Tabs/shared/FormSocialInput/SocialIcons.tsx';

interface SocialsCellProps {
  socials?: Social[] | null;
}

export const SocialsCell = ({ socials }: SocialsCellProps) => {
  if (!socials?.length) return <p className='text-gray-400'>Unknown</p>;

  return (
    <div className='flex space-evenly items-center w-full h-full'>
      {socials?.map((social) =>
        isKnownUrl(social.url) ? (
          <Tooltip label={getFormattedLink(social.url)} key={social.id}>
            <IconButton
              className='ml-1 rounded-[5px]'
              variant='ghost'
              size='xs'
              onClick={() =>
                window.open(
                  getExternalUrl(social.url ?? '/'),
                  '_blank',
                  'noopener',
                )
              }
              aria-label={social.url}
              icon={<SocialIcon url={social.url} />}
            />
          </Tooltip>
        ) : (
          <div className='flex items-center' key={social.id}>
            <Tooltip label={getFormattedLink(social.url)}>
              <IconButton
                className='ml-1 rounded-[5px]'
                variant='ghost'
                size='xs'
                onClick={() =>
                  window.open(
                    getExternalUrl(social.url ?? '/'),
                    '_blank',
                    'noopener',
                  )
                }
                aria-label='organization website'
                icon={<LinkExternal02 className='text-gray-500' />}
              />
            </Tooltip>
          </div>
        ),
      )}
    </div>
  );
};
