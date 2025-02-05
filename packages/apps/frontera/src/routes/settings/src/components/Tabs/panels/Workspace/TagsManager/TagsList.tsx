import { useState, createRef, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { TagDatum, TagStore } from '@store/Tags/Tag.store';

import { Input } from '@ui/form/Input';
import { EntityType } from '@graphql/types';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

import { TagColorPicker } from './TagColorPicker';

const entityTypes = {
  [EntityType.Organization]: { label: 'Company' },
  [EntityType.Contact]: { label: 'Contact' },
  [EntityType.LogEntry]: { label: 'Log entry' },
};

export const TagList = observer(
  ({
    tags,
    title,
    entityType,
  }: {
    title: string;
    tags: TagStore[];
    entityType?: EntityType;
  }) => {
    const store = useStore();
    const [editingTag, setEditingTag] = useState<{
      id: string;
      name: string;
    } | null>(null);
    const [newTag, setNewTag] = useState('');
    const [newTagColor, setNewTagColor] = useState('grayModern');
    const [deletingTag, setDeletingTag] = useState<TagDatum | null>(null);
    const inputRef = createRef<HTMLInputElement>();

    const [isCollapsed, setIsCollapsed] = useState(false);
    const [showNewTagInput, setShowNewTagInput] = useState<EntityType | null>(
      null,
    );
    const { open: isOpen, onOpen, onClose } = useDisclosure();

    useEffect(() => {
      if (editingTag && inputRef.current) {
        inputRef.current.focus();
        inputRef.current.select();
      }
    }, [editingTag]);

    const handleDeleteTag = (tagId: string) => {
      const tag = store.tags.value.get(tagId);

      if (tag) {
        store.organizations.toArray().forEach((organization) => {
          organization.deleteTag(tagId);
        });
        store.contacts.toArray().forEach((contact) => {
          contact.deletePersona(tagId);
        });
        store.tableViewDefs.toArray().forEach((tableViewDef) => {
          if (tableViewDef) {
            const personaFilter = tableViewDef.getFilter('CONTACTS_PERSONA');

            if (personaFilter) {
              tableViewDef.toggleFilter(personaFilter);
            }
            const tagsFilter = tableViewDef.getFilter('ORGANIZATIONS_TAGS');

            if (tagsFilter) {
              tableViewDef.toggleFilter(tagsFilter);
            }
          }
        });
        store.tags.deleteTag(tagId);
      }
      setDeletingTag(null);
    };

    const handleNewTagSubmit = (entityType: EntityType) => {
      if (newTag) {
        store.tags.create(
          { name: newTag, entityType, colorCode: newTagColor },
          {
            onSucces: (serverId) => {
              const newTagStore = store.tags.value.get(serverId);

              if (newTagStore) {
                store.tags.value.delete(serverId);

                const updatedTags = new Map([[serverId, newTagStore]]);

                store.tags.value.forEach((value, key) => {
                  if (key !== serverId) {
                    updatedTags.set(key, value);
                  }
                });
                store.tags.value = updatedTags;
              }
            },
          },
        );
        setShowNewTagInput(null);
        setNewTag('');
      }
    };

    const handleEditTag = (tagId: string, newName: string) => {
      const tag = store.tags.value.get(tagId);

      if (tag) {
        tag.update((value) => {
          value.name = newName;

          return value;
        });
      }
      setEditingTag(null);
      setNewTag('');
    };
    const deleteTagDescription = `Deleting this tag will remove it from all ${
      deletingTag?.entityType
        ? // eslint-disable-next-line @typescript-eslint/no-explicit-any
          (entityTypes as any)[deletingTag.entityType].label
        : entityTypes?.[EntityType.Organization]?.label
    } where it’s currently used.`;

    return (
      <div className='mb-6'>
        <div className='flex justify-between mb-2'>
          <div className='flex gap-2'>
            <p className='text-sm font-medium text-gray-700 '>{`${title} tags`}</p>
            {tags.length > 0 && (
              <IconButton
                size='xxs'
                variant='ghost'
                aria-label='collapse the list'
                onClick={() => setIsCollapsed(!isCollapsed)}
                icon={!isCollapsed ? <ChevronCollapse /> : <ChevronExpand />}
              />
            )}
          </div>
          {entityType && (
            <IconButton
              size='xxs'
              icon={<Plus />}
              aria-label='add new tag'
              onClick={() => {
                setShowNewTagInput(entityType);
                setNewTagColor(getRandomColor());
                setNewTag('');
              }}
            />
          )}
        </div>

        {showNewTagInput === entityType && (
          <div
            className='border border-gray-200 rounded-md mb-2 flex gap-2'
            onBlur={() => {
              if (newTag) {
                handleNewTagSubmit(entityType);
              } else {
                setShowNewTagInput(null);
              }
            }}
          >
            <TagColorPicker
              colorCode={newTagColor}
              onSelect={(color) => {
                setNewTagColor(color);
              }}
            />
            <Input
              autoFocus
              size='sm'
              value={newTag}
              variant='unstyled'
              placeholder='Enter new tag name...'
              className=' placeholder:text-sm text-sm bg-white'
              onBlur={(e) => {
                e.stopPropagation();
              }}
              onChange={(e) => {
                setNewTag(e.target.value);
                e.stopPropagation();
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleNewTagSubmit(entityType);
                }
                e.stopPropagation();
              }}
            />
          </div>
        )}

        {!isCollapsed && (
          <>
            {tags.length === 0 ? (
              <p className='text-sm text-gray-500'>No tags in sight</p>
            ) : (
              tags.map((tag) => (
                <div
                  key={tag.value.metadata.id}
                  className='py-1 max-h-[30px] mb-1 border rounded-md border-gray-200 flex justify-between items-center group bg-white'
                >
                  <TagColorPicker
                    colorCode={tag.value.colorCode}
                    onSelect={(color) => {
                      tag?.update((value) => {
                        value.colorCode = color;

                        return value;
                      });
                    }}
                  />
                  <div className='flex flex-grow '>
                    {editingTag?.id === tag.value.metadata.id ? (
                      <div className='ml-2 overflow-hidden'>
                        <Input
                          autoFocus
                          size='xxs'
                          ref={inputRef}
                          variant='unstyled'
                          className='mb-[3px] bg-white '
                          defaultValue={editingTag?.name}
                          onFocus={(e) => {
                            e.target.select();

                            if (!newTag) {
                              setNewTag(e.target.value);
                            }
                          }}
                          onBlur={(e) => {
                            const value = e.target.value.trim();

                            if (value && value !== editingTag?.name) {
                              handleEditTag(editingTag.id, value);
                            }
                            setEditingTag(null);
                            setNewTag('');
                          }}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              const value = e.currentTarget.value.trim();

                              if (value && value !== editingTag?.name) {
                                handleEditTag(editingTag.id, value);
                              }
                              setEditingTag(null);
                              setNewTag('');
                            }
                          }}
                        />
                      </div>
                    ) : (
                      <>
                        <span
                          className='cursor-pointer ml-2 text-sm break-all line-clamp-1 w-full '
                          onClick={() =>
                            setEditingTag({
                              id: tag.value.metadata.id,
                              name: tag.value.name,
                            })
                          }
                        >
                          {tag.value.name}
                        </span>
                      </>
                    )}
                  </div>
                  <div className='flex items-center opacity-0 transition-opacity duration-200 group-hover:opacity-100 pr-3'>
                    <IconButton
                      size='xxs'
                      variant='ghost'
                      aria-label='Delete tag'
                      icon={<Trash01 className='w-4 h-4' />}
                      onClick={() => {
                        setDeletingTag(tag.value);
                        onOpen();
                      }}
                    />
                  </div>
                </div>
              ))
            )}
          </>
        )}
        <ConfirmDeleteDialog
          isOpen={isOpen}
          hideCloseButton
          onClose={onClose}
          confirmButtonLabel='Delete tag'
          description={deleteTagDescription}
          label={`Delete '${deletingTag?.name}'?`}
          onConfirm={() => {
            if (deletingTag?.metadata.id) {
              handleDeleteTag(deletingTag.metadata.id);
            }
          }}
          body={
            <div className='flex flex-col gap-2'>
              <p className='text-sm'>This action cannot be undone.</p>
            </div>
          }
        />
      </div>
    );
  },
);

function getRandomColor(): string {
  const colors = [
    'grayModern',
    'error',
    'warning',
    'success',
    'grayWarm',
    'moss',
    'blueLight',
    'indigo',
    'violet',
    'pink',
  ];

  return colors[Math.floor(Math.random() * colors.length)];
}
