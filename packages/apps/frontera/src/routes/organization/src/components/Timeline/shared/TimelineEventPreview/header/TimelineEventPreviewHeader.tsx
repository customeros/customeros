import { escapeForSlackWithMarkdown } from 'slack-to-html';

import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
import { DateTimeUtils } from '@spaces/utils/date';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

interface TimelineEventPreviewHeaderProps {
  name: string;
  date?: string;
  parse?: 'slack';
  copyLabel: string;
  onClose: () => void;
  children?: React.ReactNode;
}

export const TimelineEventPreviewHeader: React.FC<
  TimelineEventPreviewHeaderProps
> = ({ date, name, onClose, copyLabel, children, parse }) => {
  const [_, copy] = useCopyToClipboard();

  const parsedName =
    parse === 'slack' ? escapeForSlackWithMarkdown(name) : name;

  return (
    <div
      className='sticky py-4 px-6 pb-1 top-0 rounded-xl'
      onClick={(e) => e.stopPropagation()}
    >
      <div>
        <div className='flex justify-between '>
          <span
            className='text-lg font-semibold text-gray-700'
            dangerouslySetInnerHTML={
              parse === 'slack' ? { __html: parsedName } : undefined
            }
          >
            {parse !== 'slack' ? name : null}
          </span>

          <div className='flex justify-end items-center'>
            {children}
            <Tooltip label={copyLabel} side='bottom' asChild={false}>
              <div>
                <IconButton
                  className='mr-1'
                  variant='ghost'
                  aria-label={copyLabel}
                  colorScheme='gray'
                  size='xs'
                  icon={<Link03 color='gray.500' height='18px' />}
                  onClick={() => copy(window.location.href)}
                />
              </div>
            </Tooltip>
            <Tooltip label='Close' aria-label='close' side='bottom'>
              <div>
                <IconButton
                  variant='ghost'
                  aria-label='Close preview'
                  colorScheme='gray'
                  size='xs'
                  icon={<XClose color='gray.500' height='24px' />}
                  onClick={onClose}
                />
              </div>
            </Tooltip>
          </div>
        </div>
        {date && (
          <span className='text-[12px] text-gray-500'>
            {DateTimeUtils.format(date, DateTimeUtils.dateWithHour)}
          </span>
        )}
      </div>
    </div>
  );
};
