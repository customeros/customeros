import type { SVGAttributes } from 'react';

import { twMerge } from 'tailwind-merge';

export type IconName =
  | 'atom-01'
  | 'alert-circle'
  | 'archive'
  | 'arrow-dropdown'
  | 'bubbles'
  | 'chevron-down'
  | 'chevron-up'
  | 'chevron-left'
  | 'chevron-right'
  | 'copy-03'
  | 'download-04'
  | 'dot-single'
  | 'dot-live-success'
  | 'dot-live-primary'
  | 'cloud-off'
  | 'invoice'
  | 'invoice-upcoming'
  | 'invoice-check'
  | 'building-07'
  | 'check-heart'
  | 'radio-dot'
  | 'users-01'
  | 'heart-hand'
  | 'radar'
  | 'signature'
  | 'target-05'
  | 'plus-circle'
  | 'coins-stacked-01'
  | 'shuffle-01'
  | 'play'
  | 'log-out-01'
  | 'layers-two-01'
  | 'settings-02'
  | 'x-close'
  | 'share-07'
  | 'tag-01'
  | 'activity-heart'
  | 'seeding'
  | 'message-x-circle'
  | 'broken-heart'
  | 'align-horizontal-centre-02'
  | 'check-verified-02'
  | 'target-04'
  | 'users-02'
  | 'key-01'
  | 'eye-off'
  | 'eye'
  | 'columns-03'
  | 'user-03'
  | 'user-01'
  | 'activity'
  | 'thumbs-down'
  | 'message-question-circle'
  | 'book-closed'
  | 'code-browser'
  | 'keyboard-02'
  | 'help-circle'
  | 'message-smile-square'
  | 'arrow-narrow-up-right'
  | 'search-sm'
  | 'command'
  | 'delete'
  | 'arrow-block-up'
  | 'arrow-block-down'
  | 'arrow-if-path'
  | 'clock'
  | 'user-plus-01'
  | 'check'
  | 'arrow-switch-horizontal-02'
  | 'x-circle'
  | 'code-square-02'
  | 'edit-03'
  | 'brush-01'
  | 'chart-breakout-circle'
  | 'trend-up-01'
  | 'bar-chart-05'
  | 'bar-chart-06'
  | 'bar-chart-circle-02'
  | 'bar-chart-circle-03'
  | 'chart-breakout-square'
  | 'trend-down-01'
  | 'presentation-chart-02'
  | 'presentation-chart-03'
  | 'bell-04'
  | 'bell-off-03'
  | 'thumbs-up'
  | 'announcement-03'
  | 'mail-01'
  | 'infinity'
  | 'refresh-ccw-02'
  | 'inbox-unread'
  | 'dots-vertical'
  | 'send-03'
  | 'mail-05'
  | 'message-chat-circle'
  | 'message-text-circle-01'
  | 'phone-call-01'
  | 'send-01'
  | 'code-02'
  | 'brackets'
  | 'cpu-chip-01'
  | 'data'
  | 'dataflow-03'
  | 'dataflow-04'
  | 'puzzle-piece-01'
  | 'variable'
  | 'award-01'
  | 'beaker-01'
  | 'beaker-02'
  | 'book-open-01'
  | 'briefcase-02'
  | 'certificate-01'
  | 'glasses-02'
  | 'graduation-hat-01'
  | 'stand'
  | 'telescope'
  | 'trophy-01'
  | 'clipboard-check'
  | 'sticker-circle'
  | 'paperclip'
  | 'alert-triangle'
  | 'building-05'
  | 'plus';

interface IconProps extends SVGAttributes<SVGElement> {
  name: IconName;
  className?: string;
}

export const Icon = ({
  name,
  fill,
  width,
  height,
  stroke,
  className,
  strokeWidth,
  ...props
}: IconProps) => (
  <svg
    viewBox='0 0 24 24'
    width={width ?? 24}
    fill={fill ?? 'none'}
    height={height ?? 24}
    strokeLinecap='round'
    strokeLinejoin='round'
    strokeWidth={strokeWidth ?? 2}
    stroke={stroke ?? 'currentColor'}
    {...props}
    className={twMerge('inline-block size-4', className)}
  >
    <use xlinkHref={`/icons.svg#${name}`} />
  </svg>
);
