import React, { useId, useState, useCallback } from 'react';

import { useCallbackRef } from '@shared/hooks/useCallbackRef';

export interface UseDisclosureProps {
  id?: string;
  open?: boolean;
  onOpen?(): void;
  onClose?(): void;
  defaultOpen?: boolean;
}

type HTMLProps = React.HTMLAttributes<HTMLElement>;

export function useDisclosure(props: UseDisclosureProps = {}) {
  const {
    onClose: onCloseProp,
    onOpen: onOpenProp,
    open: openProp,
    id: idProp,
  } = props;

  const handleOpen = useCallbackRef(onOpenProp);
  const handleClose = useCallbackRef(onCloseProp);

  const [openState, setopen] = useState(props.defaultOpen || false);

  const open = openProp !== undefined ? openProp : openState;

  const isControlled = openProp !== undefined;

  const uid = useId();
  const id = idProp ?? `disclosure-${uid}`;

  const onClose = useCallback(() => {
    if (!isControlled) {
      setopen(false);
    }
    handleClose?.();
  }, [isControlled, handleClose]);

  const onOpen = useCallback(() => {
    if (!isControlled) {
      setopen(true);
    }
    handleOpen?.();
  }, [isControlled, handleOpen]);

  const onToggle = useCallback(() => {
    if (open) {
      onClose();
    } else {
      onOpen();
    }
  }, [open, onOpen, onClose]);

  function getButtonProps(props: HTMLProps = {}): HTMLProps {
    return {
      ...props,
      'aria-expanded': open,
      'aria-controls': id,
      onClick(event) {
        props.onClick?.(event);
        onToggle();
      },
    };
  }

  function getDisclosureProps(props: HTMLProps = {}): HTMLProps {
    return {
      ...props,
      hidden: !open,
      id,
    };
  }

  return {
    open,
    onOpen,
    onClose,
    onToggle,
    isControlled,
    getButtonProps,
    getDisclosureProps,
  };
}

export type UseDisclosureReturn = ReturnType<typeof useDisclosure>;
