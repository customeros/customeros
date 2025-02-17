import { useMemo, useState } from 'react';

import Fuse from 'fuse.js';
import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';

import { Icon } from '@ui/media/Icon';
import { Input } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';

import { data, Section, SearchResult, getAllShortcuts } from './shortcuts';

export const ShortcutsPanel = observer(() => {
  const store = useStore();
  const [searchTerm, setSearchTerm] = useState('');

  const fuse = useMemo(() => {
    const allShortcuts = getAllShortcuts(data.sections);

    return new Fuse(allShortcuts, {
      keys: [{ name: 'label' }],
      threshold: 0.3,
      ignoreLocation: true,
    });
  }, []);

  const groupShortcutsBySection = (
    searchResults: SearchResult[],
  ): Section[] => {
    return searchResults.reduce<Section[]>((acc, { item }) => {
      const section = acc.find((s) => s.title === item.section);

      if (section) {
        section.shortcuts.push(item);
      } else {
        acc.push({
          title: item.section,
          shortcuts: [item],
        });
      }

      return acc;
    }, []);
  };

  const filteredSections = useMemo(() => {
    if (!searchTerm.trim()) {
      return data.sections;
    }

    const searchResults = fuse.search(searchTerm) as SearchResult[];

    return groupShortcutsBySection(searchResults);
  }, [searchTerm, fuse]);

  useKey('Escape', () => {
    store.ui.setShortcutsPanel(false);
  });

  return (
    <div className='mb-10'>
      <div className='sticky top-0 bg-white mb-2'>
        <div className='flex items-center justify-between pt-4 px-4'>
          <p className='text-base font-medium'>Keyboard shortcuts</p>
          <IconButton
            size='xs'
            variant='ghost'
            icon={<Icon name='x-close' />}
            aria-label='close-contact-overview'
            onClick={() => store.ui.setShortcutsPanel(false)}
          />
        </div>
        <div className='flex items-center gap-x-2 px-4'>
          <Icon name={'search-sm'} className='size-4 text-gray-500' />
          <Input
            autoFocus
            size='sm'
            type='text'
            variant='unstyled'
            value={searchTerm}
            placeholder='Search shortcuts'
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      <div className='w-full px-4'>
        {filteredSections.length > 0 ? (
          filteredSections.map((section, idx) => (
            <div key={idx} className='mb-5'>
              <h3 className='text-sm font-medium text-gray-700 mb-2'>
                {section.title}
              </h3>
              {section.shortcuts.map((shortcut, shortcutIdx) => (
                <div
                  key={shortcutIdx}
                  className='mb-2 flex items-center justify-between text-sm'
                >
                  <p className='text-gray-700'>{shortcut.label}</p>
                  <div className='flex items-center gap-1'>
                    {shortcut?.sequence ? (
                      <div className='flex items-center gap-1'>
                        <kbd className='px-2 py-1 bg-gray-100 rounded text-xs'>
                          {shortcut.sequence[0]}
                        </kbd>
                        <span className='text-xs text-gray-500'>then</span>
                        <kbd className='px-2 py-1 bg-gray-100 rounded text-xs'>
                          {shortcut.sequence[1]}
                        </kbd>
                      </div>
                    ) : (
                      <div className='flex items-center gap-1'>
                        {shortcut?.modifierIcon && (
                          <kbd className='px-2 py-1 bg-gray-100 rounded text-xs'>
                            <Icon
                              className='size-3'
                              name={shortcut.modifierIcon}
                            />
                          </kbd>
                        )}
                        {shortcut?.key && (
                          <kbd className='px-2 py-1 bg-gray-100 rounded text-xs'>
                            {shortcut.key}
                          </kbd>
                        )}

                        {shortcut?.keyIcon && (
                          <kbd className='px-2 py-1 bg-gray-100 rounded text-xs'>
                            <Icon className='size-3' name={shortcut.keyIcon} />
                          </kbd>
                        )}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ))
        ) : (
          <div className='text-center text-sm text-gray-500 mt-4'>
            No shortcut found by that name…
          </div>
        )}
      </div>
    </div>
  );
});
