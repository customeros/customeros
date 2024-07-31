import { escapeForSlackWithMarkdown } from 'slack-to-html';

import { DateTimeUtils } from '@utils/date';
import { Link03 } from '@ui/media/icons/Link03';
import { XClose } from '@ui/media/icons/XClose';
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
      onClick={(e) => e.stopPropagation()}
      className='sticky py-4 px-6 pb-1 top-0 rounded-xl'
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
            <Tooltip side='bottom' asChild={false} label={copyLabel}>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  className='mr-1'
                  colorScheme='gray'
                  aria-label={copyLabel}
                  onClick={() => copy(window.location.href)}
                  icon={<Link03 height='18px' color='gray.500' />}
                />
              </div>
            </Tooltip>
            <Tooltip label='Close' side='bottom' aria-label='close'>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  onClick={onClose}
                  colorScheme='gray'
                  aria-label='Close preview'
                  icon={<XClose height='24px' color='gray.500' />}
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
