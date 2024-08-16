import { twMerge } from 'tailwind-merge';

import { cn } from '@ui/utils/cn.ts';
import { colors } from '@ui/theme/colors.ts';
import { HighlightColor } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/components/Services/components/highlighters/utils.ts';

interface HighlighterProps extends React.SVGAttributes<SVGElement> {
  color?: string;
  className?: string;
  children?: React.ReactNode;
  highlightVersion?: number | string;
  backgroundColor?: 'transparent' | 'grayWarm';
}

const HighlighterVariant1 = ({
  color,
  className,
  ...props
}: HighlighterProps) => {
  return (
    <svg
      height='21'
      fill='none'
      width='100%'
      viewBox='0 0 57 21'
      preserveAspectRatio={'none'}
      {...props}
      className={twMerge('inline-block', className)}
    >
      <path
        fill='currentColor'
        d='M53 0H4V2H3V4.5H1.5V6.5H4V9H3V11H1.5V13.5H0V16.5H2.5V20.5H56.5V16.5H54V13.5H56V11H51.5V9H53V6.5H54.5V4.5H52V2.5H53V0Z'
      />
    </svg>
  );
};

const HighlighterVariant2 = ({ color, ...props }: HighlighterProps) => {
  return (
    <svg
      height='21'
      fill='none'
      width='100%'
      viewBox='0 0 54 21'
      {...props}
      preserveAspectRatio={'none'}
    >
      <path
        fill='currentColor'
        d='M51.5 0.5H1.5V1.5H2.5V5H1V7H2.5V8.5H1.5V11.5H0V16H3.5V18H1V21H53V18H51.5V16H52.5V11.5H54V8.5H50.5V7H53V5H48.5V2H51.5V0.5Z'
      />
    </svg>
  );
};

const HighlighterVariant3 = ({ color, ...props }: HighlighterProps) => {
  return (
    <svg
      height='21'
      fill='none'
      width='100%'
      viewBox='0 0 55 21'
      {...props}
      preserveAspectRatio={'none'}
    >
      <path
        fill='currentColor'
        d='M3.5 20.5H53.5V17.5H50.5V16H53V14H51.5V12.5H52.5V9.5H55V6H53.5V4H54V0H0V4H2.5V6H1.5V9.5H1V12.5H2.5V14H2V16H4.5V17H3.5V20.5Z'
      />
    </svg>
  );
};

export const Highlighter = ({
  children,
  backgroundColor = 'transparent',
  highlightVersion,
}: HighlighterProps) => {
  const color =
    backgroundColor === HighlightColor.GrayWarm
      ? colors[backgroundColor as keyof typeof colors]?.['200']
      : colors[backgroundColor as keyof typeof colors]?.['100'];

  return (
    <div className={cn('relative max-h-4 flex items-center')}>
      {color && (
        <div className='flex items-center absolute top-0 bottom-0 -right-1 -left-1 overflow-visible'>
          {(!highlightVersion || `${highlightVersion}` === '1') && (
            <HighlighterVariant1 style={{ color }} />
          )}
          {`${highlightVersion}` === '2' && (
            <HighlighterVariant2 style={{ color }} />
          )}
          {`${highlightVersion}` === '3' && (
            <HighlighterVariant3 style={{ color }} />
          )}
        </div>
      )}

      <div className='flex relative z-1 items-baseline'>{children}</div>
    </div>
  );
};
