import React, { useRef, useState, useEffect, useCallback } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';

interface PageLayoutProps {
  unstyled?: boolean;
  className?: string;
  isResizable?: boolean;
  children: React.ReactNode;
}

export const PageLayout = ({
  unstyled,
  className,
  children,
  isResizable,
}: PageLayoutProps) => {
  const [storedSidebarWidth, setStoredSidebarWidth] = useLocalStorage(
    'cos_sidebar_width',
    200,
  );
  const [sidebarWidth, setSidebarWidth] = useState(storedSidebarWidth);
  const [isDragging, setIsDragging] = useState(false);
  const dragHandleRef = useRef<HTMLDivElement>(null);

  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (!isResizable) return;
      e.preventDefault();
      setIsDragging(true);
      document.body.style.cursor = 'col-resize';
      document.body.style.userSelect = 'none';
    },
    [isResizable],
  );

  const handleMouseMove = useCallback(
    (e: MouseEvent) => {
      if (isDragging && !!isResizable) {
        const newWidth = e.clientX;

        if (newWidth >= 180 && newWidth <= 280) {
          requestAnimationFrame(() => {
            setSidebarWidth(newWidth);
          });
        }
      }
    },
    [isDragging, isResizable],
  );

  const handleMouseUp = useCallback(() => {
    if (isDragging) {
      setStoredSidebarWidth(sidebarWidth);
    }
    setIsDragging(false);
    document.body.style.removeProperty('cursor');
    document.body.style.removeProperty('user-select');
  }, [isDragging, sidebarWidth, setStoredSidebarWidth]);

  useEffect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    } else {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging, handleMouseMove, handleMouseUp]);

  if (unstyled) return <div className={className}>{children}</div>;

  return (
    <div
      className='h-screen grid bg-gray-25 relative'
      style={{
        gridTemplateAreas: `"sidebar content"`,
        gridTemplateColumns: `${
          isResizable ? sidebarWidth : '200'
        }px minmax(100px, 1fr)`,
      }}
    >
      {children}
      <div
        ref={dragHandleRef}
        onMouseDown={handleMouseDown}
        style={{
          left: `${isResizable ? sidebarWidth : 200}px`,
          width: isDragging ? '14rem' : '3px',
        }}
        className={cn(
          'absolute top-0 left-0 h-full z-10 hover:bg-transparent transition-colors hover:w-[1px]',
          {
            'border-l': isDragging,
            'hover:border-l cursor-col-resize': isResizable,
          },
        )}
      />
    </div>
  );
};
