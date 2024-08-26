import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Plus } from '@ui/media/icons/Plus';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { Tumbleweed } from '@ui/media/icons/Tumbleweed';
import { InputGroup, LeftElement } from '@ui/form/InputGroup';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Tag } from '@shared/types/__generated__/graphql.types';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

export const TagsManager = observer(() => {
  const store = useStore();
  const [newTag, setNewTag] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [showNewTagInput, setShowNewTagInput] = useState(false);
  const [editingTag, setEditingTag] = useState<{
    id: string;
    name: string;
  } | null>(null);

  const [deletingTag, setDeletingTag] = useState<Tag | null>(null);
  const { open: isOpen, onOpen, onClose } = useDisclosure();

  const handleAddNewTag = () => {
    setShowNewTagInput(true);
  };

  const handleNewTagSubmit = () => {
    if (newTag) {
      store.tags.create(
        { name: newTag },
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
      setShowNewTagInput(false);
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
  };

  const handleDeleteTag = (tagId: string) => {
    const tag = store.tags.value.get(tagId);

    if (tag) {
      store.tags.deleteTag(tagId);
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

  const organizationTagCount =
    deletingTag &&
    deletingTag.id &&
    store.organizations.toComputedArray((arr) => {
      return arr.filter((org) =>
        org.value.tags?.some((tag) => tag.id === deletingTag.id),
      );
    }).length;

  const contactTagCount =
    deletingTag &&
    deletingTag.id &&
    store.contacts.toComputedArray((arr) => {
      return arr.filter((contact) =>
        contact.value.tags?.some((persona) => persona.id === deletingTag.id),
      );
    }).length;

  const totalTagCount = Number(organizationTagCount) + Number(contactTagCount);

  const deleteTagDescription =
    totalTagCount > 1
      ? `This action will remove this tag from ${totalTagCount} contacts or organizations.`
      : totalTagCount === 1
      ? `This action will remove this tag from ${totalTagCount} contact or organization.`
      : 'This tag action will remove  this tag ';

  return (
    <>
      <div className='px-6 pb-4 max-w-[500px] h-full overflow-y-auto  border-r border-gray-200'>
        <div className='flex flex-col '>
          <div className='flex justify-between items-center pt-[5px] sticky top-0 bg-gray-25 '>
            <p className='text-gray-700 font-semibold'>Tags</p>
            <Button size='xs' leftIcon={<Plus />} onClick={handleAddNewTag}>
              New Tag
            </Button>
          </div>
          <p className='mb-4 text-sm'>Manage your workspace tags</p>

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
                onChange={(e) => {
                  setSearchTerm(e.target.value.toLowerCase());
                }}
              />
            </InputGroup>
          </div>
          {showNewTagInput && (
            <div className='border border-gray-200 rounded-md mb-1'>
              <Input
                autoFocus
                size='sm'
                value={newTag}
                variant='unstyled'
                placeholder='Enter new tag name...'
                className='pl-6 placeholder:text-sm text-sm'
                onChange={(e) => {
                  setNewTag(e.target.value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    handleNewTagSubmit();
                  }
                }}
                onBlur={() => {
                  if (newTag) {
                    handleNewTagSubmit();
                  } else {
                    setShowNewTagInput(false);
                  }
                }}
              />
            </div>
          )}
          {filteredTags.length === 0 ? (
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
            filteredTags.map((tag) => (
              <div
                key={tag.value.id}
                className='py-1 max-h-[32px] mb-1 border rounded-md border-gray-200 flex justify-between items-center group'
              >
                <div className='flex-grow'>
                  {editingTag?.id === tag.value.id ? (
                    <Input
                      autoFocus
                      size='xs'
                      variant='unstyled'
                      className='ml-6 mt-[2px]'
                      defaultValue={tag.value.name}
                      onChange={(e) => setNewTag(e.target.value)}
                      onFocus={(e) => {
                        e.target.select();
                      }}
                      onBlur={() => {
                        handleEditTag(tag.value.id, newTag);
                        setEditingTag(null);
                      }}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                          handleEditTag(tag.value.id, e.currentTarget.value);
                        }
                      }}
                    />
                  ) : (
                    <span
                      className='cursor-pointer ml-6 text-sm line-clamp-1'
                      onClick={() =>
                        setEditingTag({
                          id: tag.value.id,
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
                    size='xs'
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
          if (deletingTag?.id) {
            handleDeleteTag(deletingTag.id);
          }
        }}
        body={
          <div className='flex flex-col gap-2'>
            <p>This action cannot be undone.</p>
          </div>
        }
      />
    </>
  );
});
