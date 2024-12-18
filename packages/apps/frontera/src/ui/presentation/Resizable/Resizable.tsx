import { useState, useEffect } from 'react';

interface ResizableProps {
  className?: string;
  isMovable?: boolean; // Enable or disable moving
  defaultWidth?: number;
  defaultHeight?: number;
  defaultPosition?: { x: number; y: number };
  resizeDirection?: 'both' | 'horizontal' | 'vertical'; // Resize direction
  children:
    | React.ReactNode
    | ((
        isDragging: boolean,
        startMove: (e: React.MouseEvent, ignoreKeys: boolean) => void,
      ) => React.ReactNode);
}

export const Resizable = ({
  children,
  className,
  defaultWidth = 200,
  defaultHeight = 200,
  defaultPosition = { x: 0, y: 0 },
  isMovable = true,
  resizeDirection = 'both', // Default: resizable in both directions
}: ResizableProps) => {
  const [dimensions, setDimensions] = useState({
    width: defaultWidth,
    height: defaultHeight,
  });
  const [position, setPosition] = useState(defaultPosition);
  const [isDragging, setIsDragging] = useState(false);
  const [isMoving, setIsMoving] = useState(false);
  const [draggingSide, setDraggingSide] = useState<string | null>(null);

  const disableTextSelection = () => {
    document.body.style.userSelect = 'none';
  };

  const enableTextSelection = () => {
    document.body.style.userSelect = 'auto';
  };

  const startResize = (side: string) => {
    setIsDragging(true);
    setDraggingSide(side);
    disableTextSelection();
  };

  const startMove = (e: React.MouseEvent, ignoreKeys: boolean = false) => {
    if (isMovable && (ignoreKeys || e.metaKey || e.ctrlKey)) {
      setIsMoving(true);
      setDraggingSide(null); // Ensure no resizing occurs during move
      disableTextSelection();
    }
    e.stopPropagation();
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (isDragging && draggingSide) {
      // Handle resizing
      setDimensions((prev) => {
        const { width, height } = prev;
        const deltaX = e.movementX;
        const deltaY = e.movementY;

        const newDimensions = { ...prev };

        if (resizeDirection !== 'vertical') {
          if (draggingSide.includes('right'))
            newDimensions.width = width + deltaX;
          if (draggingSide.includes('left'))
            newDimensions.width = Math.max(50, width - deltaX);
        }

        if (resizeDirection !== 'horizontal') {
          if (draggingSide.includes('bottom'))
            newDimensions.height = height + deltaY;
          if (draggingSide.includes('top'))
            newDimensions.height = Math.max(50, height - deltaY);
        }

        return newDimensions;
      });
    } else if (isMoving) {
      // Handle dragging
      setPosition((prev) => ({
        x: prev.x + e.movementX,
        y: prev.y + e.movementY,
      }));
    }
  };

  const stopInteractions = () => {
    setIsDragging(false);
    setIsMoving(false);
    setDraggingSide(null);
    enableTextSelection();
  };

  useEffect(() => {
    if (isDragging || isMoving) {
      window.addEventListener('mousemove', handleMouseMove);
      window.addEventListener('mouseup', stopInteractions);
    } else {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', stopInteractions);
    }

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', stopInteractions);
    };
  }, [isDragging, isMoving]);

  const containerStyle: React.CSSProperties = {
    position: isMovable ? 'absolute' : 'relative',
    top: isMovable ? position.y : undefined,
    left: isMovable ? position.x : undefined,
    width: dimensions.width,
    height: resizeDirection === 'horizontal' ? '100%' : dimensions.height,
    cursor: isMovable
      ? isMoving
        ? 'grabbing'
        : isDragging
        ? 'drag'
        : 'default'
      : isDragging
      ? 'default'
      : 'auto',
  };

  return (
    <div className={className} style={containerStyle} onMouseDown={startMove}>
      {typeof children === 'function'
        ? children(isDragging || isMoving, startMove)
        : children}

      {/* Resize Handles */}
      {resizeDirection !== 'horizontal' && (
        <>
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('top');
            }}
            style={{
              position: 'absolute',
              top: 0,
              left: '50%',
              width: '100%',
              height: '10px',
              cursor: 'ns-resize',
              transform: 'translate(-50%, -50%)',
            }}
          />
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('bottom');
            }}
            style={{
              position: 'absolute',
              bottom: 0,
              left: '50%',
              width: '100%',
              height: '10px',
              cursor: 'ns-resize',
              transform: 'translate(-50%, 50%)',
            }}
          />
        </>
      )}
      {resizeDirection !== 'vertical' && (
        <>
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('left');
            }}
            style={{
              position: 'absolute',
              top: '50%',
              left: 0,
              height: '100%',
              width: '10px',
              cursor: 'ew-resize',
              transform: 'translate(-50%, -50%)',
            }}
          />
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('right');
            }}
            style={{
              position: 'absolute',
              top: '50%',
              right: 0,
              height: '100%',
              width: '10px',
              cursor: 'ew-resize',
              transform: 'translate(50%, -50%)',
            }}
          />
        </>
      )}
      {resizeDirection === 'both' && (
        <>
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('top-left');
            }}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '10px',
              height: '10px',
              cursor: 'nwse-resize',
            }}
          />
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('top-right');
            }}
            style={{
              position: 'absolute',
              top: 0,
              right: 0,
              width: '10px',
              height: '10px',
              cursor: 'nesw-resize',
            }}
          />
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('bottom-left');
            }}
            style={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              width: '10px',
              height: '10px',
              cursor: 'nesw-resize',
            }}
          />
          <div
            onMouseDown={(e) => {
              e.stopPropagation();
              startResize('bottom-right');
            }}
            style={{
              position: 'absolute',
              bottom: 0,
              right: 0,
              width: '10px',
              height: '10px',
              cursor: 'nwse-resize',
            }}
          />
        </>
      )}
    </div>
  );
};
