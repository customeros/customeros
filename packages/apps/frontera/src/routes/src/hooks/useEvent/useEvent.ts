import { useRef, useEffect } from 'react';

/**
 * Custom React hook to send a JavaScript event with a payload, dispatch actions, and clean up listeners.
 * @param {string} eventName - The name of the event to dispatch.
 */

function useEvent<T>(eventName: string, listener?: (data: T) => void) {
  const eventListenerRef = useRef<((event: CustomEvent<T>) => void) | null>(
    null,
  );

  /**
   * Dispatch an event with a custom payload.
   * @param {T} payload - Data to send with the event.
   */
  const dispatchEvent = (payload: T) => {
    const event = new CustomEvent<T>(eventName, { detail: payload });

    window.dispatchEvent(event);
  };

  useEffect(() => {
    if (!listener) return;

    const _listener = (event: CustomEvent<T>) => {
      listener(event.detail);
    };

    eventListenerRef.current = _listener;
    window.addEventListener(eventName, _listener as EventListener);

    return () => {
      if (eventListenerRef.current) {
        window.removeEventListener(
          eventName,
          eventListenerRef.current as EventListener,
        );
        eventListenerRef.current = null;
      }
    };
  }, []);

  return { dispatchEvent };
}

export { useEvent };
