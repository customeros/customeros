'use client';
import React, { FC, MouseEventHandler, ReactElement } from 'react';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import ReplyMany from '@spaces/atoms/icons/ReplyMany';
import Reply from '@spaces/atoms/icons/Reply';
import Forward from '@spaces/atoms/icons/Forward';
import { Tooltip } from '@ui/presentation/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { Box, Flex } from '@chakra-ui/react';

const REPLY_MODE = 'reply';
const REPLY_ALL_MODE = 'reply-all';
const FORWARD_MODE = 'forward';

const TooltipButton: FC<{
  label: string;
  onClick: MouseEventHandler<HTMLButtonElement>;
  children: ReactElement;
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
      pl={1}
      pr={1}
      onClick={onClick}
      icon={children}
    />
  </Tooltip>
);

interface ButtonsProps {
  handleModeChange: (mode: string) => void;
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
    translateY='-16px'
  >
    <TooltipButton label='Reply' onClick={() => handleModeChange(REPLY_MODE)}>
      <Reply height='16px' color='gray.400' />
    </TooltipButton>
    <TooltipButton
      mx={1}
      label='Reply all'
      onClick={() => handleModeChange(REPLY_ALL_MODE)}
    >
      <ReplyMany height='14px' color='gray.400' />
    </TooltipButton>
    <TooltipButton
      label='Forward'
      onClick={() => handleModeChange(FORWARD_MODE)}
    >
      <Forward height='14px' color='gray.400' />
    </TooltipButton>
  </Flex>
);
