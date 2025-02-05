import markdownToTxt from 'markdown-to-txt';

import { DateTimeUtils } from '@utils/date';
import { XClose } from '@ui/media/icons/XClose';
import { Link01 } from '@ui/media/icons/Link01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

interface TimelineEventPreviewHeaderProps {
  name: string;
  date?: string;
  copyLabel: string;
  onClose: () => void;
  children?: React.ReactNode;
}

export const TimelineEventPreviewHeader = ({
  date,
  name,
  onClose,
  copyLabel,
  children,
}: TimelineEventPreviewHeaderProps) => {
  const [_, copy] = useCopyToClipboard();
  const parsedName = markdownToTxt(name);

  return (
    <div
      onClick={(e) => e.stopPropagation()}
      className='sticky py-4 px-6 pb-1 top-0 rounded-xl'
    >
      <div>
        <div className='flex justify-between '>
          <span className='text-lg font-semibold text-gray-700'>
            {parsedName}
          </span>

          <div className='flex justify-end items-baseline'>
            {children}
            <Tooltip side='bottom' asChild={false} label={copyLabel}>
              <div>
                <IconButton
                  size='xs'
                  variant='ghost'
                  className='mr-1'
                  colorScheme='gray'
                  aria-label={copyLabel}
                  icon={<Link01 className='text-gray-500' />}
                  onClick={() => copy(window.location.href, 'Link copied')}
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
                  icon={<XClose className='text-gray-500' />}
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
