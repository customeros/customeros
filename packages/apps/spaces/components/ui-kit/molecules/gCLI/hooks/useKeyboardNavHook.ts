import {useState, useEffect, KeyboardEventHandler} from 'react';

function useKeyboardNavigation(navItems: any, initialIndex = 0) {
    const [activeIndex, setActiveIndex] = useState(initialIndex);

    useEffect(() => {
        function handleKeyDown(event:KeyboardEvent) {
            switch (event.key) {
                case 'ArrowUp':
                    event.preventDefault();
                    setActiveIndex(prevActiveIndex => (prevActiveIndex - 1 + navItems.length) % navItems.length);
                    break;
                case 'ArrowDown':
                    event.preventDefault();
                    setActiveIndex(prevActiveIndex => (prevActiveIndex + 1) % navItems.length);
                    break;
                case 'Enter':
                    navItems[activeIndex].onClick();
                    break;
                default:
                    break;
            }
        }

        document.addEventListener('keydown', handleKeyDown);

        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [navItems, activeIndex]);

    return activeIndex;
}