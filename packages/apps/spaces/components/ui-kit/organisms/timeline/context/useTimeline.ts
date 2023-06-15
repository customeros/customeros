import { useContext } from 'react';
import { TimelineContext } from './timelineContext';

export const useTimeline = () => useContext(TimelineContext);
