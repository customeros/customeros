import { Copy02 } from '@ui/media/icons/Copy02';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { IconButton } from '@ui/form/IconButton/IconButton';
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
    <div className='flex items-center group/linkedin'>
      <Tooltip label={url ?? ''}>
        <p
          onClick={() => window.open(url, '_blank', 'noopener')}
          className='text-gray-700 truncate cursor-default hover:underline hover:cursor-pointer'
        >
          {displayLink}
        </p>
      </Tooltip>

      <IconButton
        size='xxs'
        variant='ghost'
        aria-label='social-link'
        icon={<Copy02 className='text-gray-500' />}
        onClick={() => copyToClipboard(url, 'LinkedIn profile copied')}
        className='ml-1 rounded-[5px] opacity-0 group-hover/linkedin:opacity-100'
      />
    </div>
  );
};
