'use client';
import { FC, ReactElement, MouseEventHandler } from 'react';

import { Flex } from '@ui/layout/Flex';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/presentation/Tooltip';
import { CornerUpLeft } from '@ui/media/icons/CornerUpLeft';
import { CornerUpLeft2 } from '@ui/media/icons/CornerUpLeft2';
import { CornerUpRight } from '@ui/media/icons/CornerUpRight';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

const TooltipButton: FC<{
  label: string;
  children: ReactElement;
  onClick: MouseEventHandler<HTMLButtonElement>;
}> = ({ label, children, onClick }) => (
  <Tooltip label={label}>
    <IconButton
      variant='ghost'
      aria-label={label}
      fontSize='14px'
      color='gray.400'
      borderRadius={0}
      marginInlineStart={0}
      size='xxs'
      pl={2}
      pr={2}
      onClick={onClick}
      icon={children}
    />
  </Tooltip>
);

interface ButtonsProps {
  handleModeChange: (mode: 'reply' | 'reply-all' | 'forward') => void;
}

export const ModeChangeButtons: FC<ButtonsProps> = ({ handleModeChange }) => (
  <Flex
    overflow='hidden'
    position='absolute'
    border='1px solid var(--gray-200, #EAECF0)'
    borderRadius={16}
    height='24px'
    gap={0}
    color='gray.25'
    background='gray.25'
    transform='translateY(-16px)'
  >
    <TooltipButton label='Reply' onClick={() => handleModeChange(REPLY_MODE)}>
      <CornerUpLeft height='16px' color='gray.400' />
    </TooltipButton>
    <TooltipButton
      label='Reply all'
      onClick={() => handleModeChange(REPLY_ALL_MODE)}
    >
      <CornerUpLeft2 height='14px' color='gray.400' />
    </TooltipButton>
    <TooltipButton
      label='Forward'
      onClick={() => handleModeChange(FORWARD_MODE)}
    >
      <CornerUpRight height='14px' color='gray.400' />
    </TooltipButton>
  </Flex>
);
