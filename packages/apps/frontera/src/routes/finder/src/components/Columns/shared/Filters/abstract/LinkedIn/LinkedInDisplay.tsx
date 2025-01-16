import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { getExternalUrl, getFormattedLink } from '@utils/getExternalLink';

interface LinkedInDisplayProps {
  link: string;
  type: string;
  alias?: string;
}

export const LinkedInDisplay = ({
  alias,
  link,
  type,
}: LinkedInDisplayProps) => {
  const [_, copyToClipboard] = useCopyToClipboard();

  const formattedLink = getFormattedLink(link).replace(
    /^linkedin\.com\/(?:in\/|company\/)?/,
    '/',
  );

  const displayLink = alias ? `/${alias}` : formattedLink;
  const url = formattedLink
    ? link.includes('linkedin')
      ? getExternalUrl(`https://linkedin.com/${type}${displayLink}`)
      : getExternalUrl(link)
    : '';

  return (
    <div className='flex items-center group'>
      <Tooltip label={url ?? ''}>
        <p
          className='text-gray-700 truncate cursor-default'
          onClick={() => copyToClipboard(url, 'LinkedIn profile copied')}
        >
          {displayLink}
        </p>
      </Tooltip>

      <IconButton
        size='xxs'
        variant='ghost'
        aria-label='contact website'
        icon={<LinkExternal02 className='text-gray-500' />}
        onClick={() => window.open(url, '_blank', 'noopener')}
        className='ml-1 rounded-[5px] opacity-0 group-hover:opacity-100'
      />
    </div>
  );
};
