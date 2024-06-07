import { useState, useEffect } from 'react';

import { Tooltip, TooltipProps } from '@ui/overlay/Tooltip/Tooltip';

export const TableCellTooltip = (
  props: TooltipProps & { targetRef: React.RefObject<HTMLElement> },
) => {
  const [isOverflowing, setIsOverflowing] = useState(false);

  useEffect(() => {
    const element = props.targetRef.current;
    if (element) {
      const originalDisplay = element.style.display;
      element.style.display = 'block';

      const originalClientWidth = element.clientWidth;

      const clone = element.cloneNode(true) as HTMLElement;
      clone.style.display = 'inline-block';
      clone.style.position = 'absolute';
      clone.style.visibility = 'hidden';
      clone.style.whiteSpace = 'nowrap';

      document.body.appendChild(clone);
      const cloneScrollWidth = clone.clientWidth;
      document.body.removeChild(clone);

      const isOverflow = cloneScrollWidth > originalClientWidth;
      element.style.display = originalDisplay;

      setIsOverflowing(isOverflow);
    }
  }, []);

  if (!isOverflowing) return <>{props.children}</>;

  return <Tooltip {...props} />;
};
