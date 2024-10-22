import { useState, useEffect, useCallback } from 'react';

/**
 * Detect if the user currently presses or releases a pointer/mouse button.
 */
export function usePointerInteractions() {
  const [isPointerDown, setIsPointerDown] = useState(false);
  const [isPointerReleased, setIsPointerReleased] = useState(true);

  // Use useCallback to maintain function reference
  const handlePointerUp = useCallback(() => {
    setIsPointerDown(false);
    setIsPointerReleased(true);
  }, []);

  const handlePointerDown = useCallback(() => {
    setIsPointerDown(true);
    setIsPointerReleased(false);
  }, []);

  useEffect(() => {
    // Use pointer events instead of mouse events
    window.addEventListener('pointerdown', handlePointerDown);
    window.addEventListener('pointerup', handlePointerUp);

    // Clean up both event listeners
    return () => {
      window.removeEventListener('pointerdown', handlePointerDown);
      window.removeEventListener('pointerup', handlePointerUp);
    };
  }, [handlePointerDown, handlePointerUp]);

  return { isPointerDown, isPointerReleased };
}
