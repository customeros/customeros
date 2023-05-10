import React, {
  ButtonHTMLAttributes,
  useEffect,
  useRef,
  useState,
} from 'react';
import styles from './tags.module.scss';
import { AutoComplete } from 'primereact/autocomplete';
import { capitalizeFirstLetter } from '../../../../utils';

interface Tag {
  id: string;
  name: string;
}

// todo delete - click - animate 300ms- delete
// todo input should looks like tag pill
// todo delete button should appear only in edit mode
// todo add heading to contacts in org
// todo emails and phones should be actionable - show that there are more if there are
export interface TagProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  onNewTag: (tagName: string) => void;
  onTagChange: (tag: Tag) => void;
  onTagRemove: (id: string) => void;
  tags: Array<Tag>;
  options: Array<Tag>;
  onSetTags: (tags: Array<Tag>) => void;
  onTagSelect: (tag: Tag) => void;
  onTagDelete: (id: string) => void;
}
export const TagInput = ({
  onNewTag,
  onTagChange,
  onTagRemove,
  tags,
  options,
  onSetTags,
  onTagSelect,
}: TagProps) => {
  const [filteredOptions, setFilteredOptions] = useState(options);
  const inputRef = useRef(null);

  useEffect(() => {
    setFilteredOptions(options);
  }, [options.length]);
  const handleNewTag = (tag: Tag | string) => {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-expect-error
    if (onNewTag) onNewTag(tag);
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //@ts-expect-error
    if (onTagChange) onTagChange(tag);
  };

  const notDuplicate = (newTagName: string): boolean => {
    return !tags.map(({ name }) => name).includes(newTagName);
  };

  const addTag = (tag: Tag) => {
    if (notDuplicate(tag.name)) {
      onSetTags([...tags, tag]);
      handleNewTag(tag);
    }
  };
  const handleKeyDown = (e: any) => {
    const {
      key,
      target: { value },
    } = e;
    switch (key) {
      case 'Tab':
        if (value) e.preventDefault();
        break;
      case 'Enter':
      case ',':
        {
          const trimmedValue = value.trim();
          if (trimmedValue && notDuplicate(trimmedValue)) {
            addTag(trimmedValue);
          }
          if (inputRef?.current) {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            //@ts-expect-error
            inputRef.current.value = '';
          }
        }
        break;
      case 'Backspace':
        if (!value) {
          if (tags.length > 0) {
            onTagRemove(tags[tags.length - 1]?.id);
          }
        }
        break;
    }
  };

  const search = (event: any) => {
    const query = event.query;
    const filteredItems = (options || []).filter(
      (item) => item.name.toLowerCase().indexOf(query.toLowerCase()) !== -1,
    );

    setFilteredOptions(filteredItems || []);
  };

  return (
    <div className={`${styles.tagInputWrapper}`}>
      <AutoComplete
        field='name'
        inputRef={inputRef}
        multiple
        placeholder='Add a tag...'
        className={`${styles.autocomplete} tagInput`}
        value={tags}
        itemTemplate={(tag: Tag) => {
          return (
            <span className={styles.option} onClick={() => onTagSelect(tag)}>
              {capitalizeFirstLetter(tag.name)?.split('_')?.join(' ')}
            </span>
          );
        }}
        suggestions={filteredOptions}
        completeMethod={search}
        onChange={(e: { value: Array<Tag> }) => {
          if (tags.length < e.value.length) {
            onTagSelect(e.value[e.value.length - 1]);
          }
        }}
        onKeyDown={handleKeyDown}
      />
    </div>
  );
};
