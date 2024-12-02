import { useState, createRef, useEffect } from 'react';

import { observer } from 'mobx-react-lite';
import { TagDatum } from '@store/Tags/Tag.store';

import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { EntityType } from '@shared/types/__generated__/graphql.types';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

const entityTypes = {
  [EntityType.Organization]: { label: 'Organization' },
  [EntityType.Contact]: { label: 'Contact' },
  [EntityType.LogEntry]: { label: 'Log entry' },
};
export const TagsManager = observer(() => {
  const store = useStore();
  const [newTag, setNewTag] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [showNewTagInput, setShowNewTagInput] = useState<EntityType | null>(
    null,
  );
  const [editingTag, setEditingTag] = useState<{
    id: string;
    name: string;
  } | null>(null);

  const [deletingTag, setDeletingTag] = useState<TagDatum | null>(null);
  const inputRef = createRef<HTMLInputElement>();
  const { open: isOpen, onOpen, onClose } = useDisclosure();

  const handleNewTagSubmit = (entityType: EntityType) => {
    if (newTag) {
      store.tags.create(
        { name: newTag, entityType },
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

  const filteredTags = store.tags.toComputedArray((arr) => {
    arr = arr.filter((entity) => {
      const name = entity.value.name.toLowerCase().includes(searchTerm || '');

      return name;
    });

    return arr;
  });

  const filteredTagsByType = (entityType: EntityType) => {
    const allTags = store.tags.getByEntityType(entityType);

    if (searchTerm) {
      return allTags.filter((tag) =>
        tag.value.name.toLowerCase().includes(searchTerm),
      );
    }

    return allTags;
  };

  const deleteTagDescription = `Deleting this tag will remove it from all ${
    deletingTag?.entityType
      ? // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (entityTypes as any)[deletingTag.entityType].label
      : entityTypes[EntityType.Organization].label
  } where it’s currently used.`;

  const tags = store.tags.toArray().length;

  useEffect(() => {
    if (editingTag && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [editingTag]);

  const TagList = ({
    tags,
    title,
    entityType,
  }: {
    title: string;
    entityType?: EntityType;
    tags: typeof filteredTags;
  }) => {
    const [isCollapsed, setIsCollapsed] = useState(false);

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
                setNewTag('');
              }}
            />
          )}
        </div>

        {showNewTagInput === entityType && (
          <div className='border border-gray-200 rounded-md mb-2'>
            <Input
              autoFocus
              size='sm'
              value={newTag}
              variant='unstyled'
              placeholder='Enter new tag name...'
              className='pl-6 placeholder:text-sm text-sm bg-white'
              onChange={(e) => {
                setNewTag(e.target.value);
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleNewTagSubmit(entityType);
                }
                e.stopPropagation();
              }}
              onBlur={() => {
                if (newTag) {
                  handleNewTagSubmit(entityType);
                } else {
                  setShowNewTagInput(null);
                }
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
                  className='py-1 max-h-[32px] mb-1 border rounded-md border-gray-200 flex justify-between items-center group bg-white'
                >
                  <div className='flex-grow '>
                    {editingTag?.id === tag.value.metadata.id ? (
                      <div className='ml-6 overflow-hidden'>
                        <Input
                          autoFocus
                          size='xs'
                          ref={inputRef}
                          variant='unstyled'
                          className='mb-[1px] bg-white'
                          defaultValue={newTag || tag.value.name}
                          onBlur={() => {
                            handleEditTag(
                              tag.value.metadata.id,
                              newTag || tag.value.name,
                            );
                            setEditingTag(null);
                          }}
                          onChange={(e) => {
                            const trimmedValue = e.target.value.trim();

                            if (trimmedValue.length > 0) {
                              setNewTag(trimmedValue);
                            }
                            e.stopPropagation();
                          }}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              handleEditTag(
                                tag.value.metadata.id,
                                e.currentTarget.value,
                              );
                            }
                            e.stopPropagation();
                          }}
                        />
                      </div>
                    ) : (
                      <span
                        className='cursor-pointer ml-6 text-sm break-all line-clamp-1 '
                        onClick={() =>
                          setEditingTag({
                            id: tag.value.metadata.id,
                            name: tag.value.name,
                          })
                        }
                      >
                        {tag.value.name}
                      </span>
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
      </div>
    );
  };

  return (
    <>
      <div className='px-6 pb-4 max-w-[500px] h-full overflow-y-auto border-r border-gray-200'>
        <div className='flex flex-col'>
          <div className='flex justify-between items-center pt-2 sticky top-0 bg-gray-25'>
            <p className='text-gray-700 font-semibold'>Tags</p>
          </div>
          <p className='mb-4 text-sm'>Manage your workspace tags</p>

          {tags > 0 && (
            <div className='mb-4'>
              <InputGroup className='gap-2'>
                <LeftElement>
                  <SearchSm />
                </LeftElement>
                <Input
                  size='xs'
                  className='w-full'
                  variant='unstyled'
                  placeholder='Search tags...'
                  onKeyDown={(e) => {
                    e.stopPropagation();
                  }}
                  onChange={(e) => {
                    setSearchTerm(e.target.value.toLowerCase());
                    e.stopPropagation();
                  }}
                />
              </InputGroup>
            </div>
          )}

          {filteredTags.length === 0 && searchTerm ? (
            <div className='flex justify-center items-center h-full'>
              <div className='flex flex-col items-center mt-4 gap-2'>
                <Tumbleweed className='w-8 h-8 text-gray-400' />
                <span className='text-sm text-gray-500'>
                  Empty here in
                  <span className='font-semibold'> No Resultsville</span>
                </span>
              </div>
            </div>
          ) : (
            <>
              {Object.entries(entityTypes).map(([type, meta]) => (
                <TagList
                  key={type}
                  title={meta.label}
                  entityType={type as EntityType}
                  tags={filteredTagsByType(type as EntityType)}
                />
              ))}
            </>
          )}
        </div>
      </div>
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
    </>
  );
});
