import * as React from 'react';
import styles from './message.module.scss';
import { DialogContent } from './DialogContent';
import { AnalysisContent } from './AnalysisContent';
import { MessageIcon, Phone } from '../icons';
import classNames from 'classnames';

interface Props {
  message: string;

  date: any;
  index: number;
  mode: 'CHAT' | 'PHONE_CALL' | 'LIVE_CONVERSATION';
}

export const Message = ({ message }: Props) => {
  return <div style={{ width: '100%' }}>{message}</div>;
};
