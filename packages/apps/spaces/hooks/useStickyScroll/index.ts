/*
  Idea and code comes from library https://github.com/vytenisu/use-chat-scroll/blob/master/lib/sticky-scroll.ts
  library was not used directly due to imports problem in next js
 */

import { useState, useRef, useEffect, MutableRefObject } from 'react';
import { useScroll } from './useScroll';

/**
 * React hook for keeping HTML element scroll at the bottom when content updates (if it is already at the bottom).
 * @param targetRef Reference of scrollable HTML element.
 * @param data Array of some data items displayed in a scrollable HTML element. It should normally come from a state.
 */
export const useStickyScroll = (
  targetRef: MutableRefObject<Element>,
  data: any[],
  options?: IUseStickyScrollOptions,
): IUseStickyScrollResponse => {
  const [enabled, setEnabled] = useState<boolean>(options?.enabled ?? true);
  const [sticky, setSticky] = useState<boolean>(true);
  const stickyRef = useRef(sticky);

  const moveScrollToBottom = () => {
    targetRef.current.scrollTop = targetRef.current.scrollHeight;
  };

  useEffect(() => {
    stickyRef.current = sticky;

    if (sticky) {
      moveScrollToBottom();
    }
  }, [data.length, targetRef, sticky]);

  const updateStuckToBottom = () => {
    const { scrollHeight, clientHeight, scrollTop } = targetRef.current;
    const currentlyAtBottom = scrollHeight === scrollTop + clientHeight;

    if (stickyRef.current && !currentlyAtBottom) {
      setSticky(false);
    } else if (!stickyRef.current && currentlyAtBottom) {
      setSticky(true);
    }
  };

  const handleScroll = () => {
    updateStuckToBottom();
  };

  const { setScrollEventHandler } = useScroll(targetRef);

  useEffect(() => {
    if (enabled) {
      setScrollEventHandler(handleScroll);
    } else {
      setScrollEventHandler(() => {
        return;
      });
    }
  }, [enabled]);

  /**
   * Scrolls to bottom.
   */
  const scrollToBottom = () => {
    moveScrollToBottom();
    setSticky(true);
  };

  /**
   * Enables sticky scroll behavior.
   */
  const enable = () => setEnabled(true);

  /**
   * Disables sticky scroll behavior.
   */
  const disable = () => setEnabled(false);

  return {
    enabled,
    sticky,
    scrollToBottom,
    enable,
    disable,
  };
};

/**
 * Accepted options for customizing useStickyScroll hook.
 */
export interface IUseStickyScrollOptions {
  /**
   * Defines whether sticky scroll behavior is enabled initially.
   */
  enabled?: boolean;
}

/**
 * Flags and methods provided by useStickyScroll hook.
 */
export interface IUseStickyScrollResponse {
  /**
   * True when sticky scroll behavior is enabled.
   */
  enabled: boolean;

  /**
   * True when scroll is stuck to the bottom of target element.
   */
  sticky: boolean;

  /**
   * Scrolls to bottom of the target element.
   */
  scrollToBottom: () => void;

  /**
   * Enables sticky scroll behavior.
   */
  enable: () => void;

  /**
   * Disables sticky scroll behavior.
   */
  disable: () => void;
}
