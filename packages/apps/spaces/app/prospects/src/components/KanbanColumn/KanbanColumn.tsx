import React, { useState } from 'react';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Plus } from '@ui/media/icons/Plus';
import { Check } from '@ui/media/icons/Check';
import { Organization } from '@graphql/types';
import { ResizableInput } from '@ui/form/Input';
import { Skeleton } from '@ui/feedback/Skeleton';
import { IconButton } from '@ui/form/IconButton';
import { uuidv4 } from '@spaces/utils/generateUuid';

import { KanbanCard } from '../KanbanCard/KanbanCard';
import { useOrganizationsPageMethods } from '../../hooks';

interface CardColumnProps {
  title: string;
  cardCount: number;
  isLoading: boolean;
  cards: Organization[];
}

export const KanbanColumn = ({
  title,
  cardCount,
  cards,
  isLoading,
}: CardColumnProps) => {
  const [newData, setNewData] = useState<Array<{ id: string; name: string }>>(
    [],
  );
  const { createOrganization } = useOrganizationsPageMethods();
  const handleAddNew = () => {
    setNewData((prev) => [...prev, { id: uuidv4(), name: 'Unnamed' }]);
  };

  const handleUpdateNewData = (id: string, newName: string) => {
    setNewData((prev) =>
      prev.map((item) => (item.id === id ? { ...item, name: newName } : item)),
    );
  };

  const handleRemoveNewData = (id: string) => {
    setNewData((prev) => prev.filter((item) => item.id !== id));
  };

  const handleSaveNewData = (data: { id: string; name: string }) => {
    createOrganization.mutate(
      {
        input: {
          name: data.name,
        },
      },
      {
        onSuccess: () => {
          handleRemoveNewData(data.id);
        },
      },
    );
  };
  const getCardStyle = (title: string) => {
    switch (title) {
      case 'New prospects':
        return 'border border-l-4 border-cyan-300';
      case 'Contacted':
        return 'border border-l-4 border-success-300';
      case 'Opportunity':
        return 'border border-l-4 border-purple-300';
      case 'Abandoned':
        return 'border border-l-4 border-gray-400';
      default:
        return '';
    }
  };

  return (
    <div className='flex flex-col flex-shrink-0 w-72 '>
      <div className='flex items-center justify-between flex-shrink-0 h-10 px-2'>
        <div className='flex'>
          <span className='block text-sm font-semibold'>{title}</span>
          <span
            className={cn(
              'flex items-center justify-center w-5 h-5 ml-2 text-sm font-semibold  rounded text-gray-500 ',
              {
                'bg-cyan-50 text-cyan-500': title === 'New prospects',
                'bg-success-50 text-success-500': title === 'Contacted',
                'bg-purple-50 text-purple-500': title === 'Opportunity',
                'bg-gray-100 text-gray-500': title === 'Abandoned',
              },
            )}
          >
            {cardCount}
          </span>
        </div>

        {title === 'New prospects' && (
          <IconButton
            aria-label={'Add new prospect'}
            icon={<Plus />}
            variant='ghost'
            size='xs'
            onClick={handleAddNew}
          />
        )}
      </div>
      <div className='flex flex-col pb-2 overflow-auto pr-2'>
        {newData.map((data) => (
          <div
            key={data.id}
            className={cn(
              'relative flex flex-col items-start p-4 mt-3 bg-white rounded-lg cursor-pointer bg-opacity-90 group hover:bg-opacity-100',
              getCardStyle(title),
            )}
          >
            <ResizableInput
              value={data.name}
              className='text-sm font-medium shadow-none p-0 min-h-5'
              autoFocus
              onChange={(e) => handleUpdateNewData(data.id, e.target.value)}
            />

            <div className='flex justify-end w-full'>
              <IconButton
                variant='ghost'
                size='xs'
                aria-label='Cancel'
                className='p-1'
                icon={<X />}
                onClick={() => handleRemoveNewData(data.id)}
              />
              <IconButton
                variant='ghost'
                size='xs'
                aria-label='Save'
                className='p-1'
                isLoading={createOrganization.isPending}
                icon={<Check />}
                onClick={() => handleSaveNewData(data)}
              />
            </div>
          </div>
        ))}
        {isLoading && (
          <>
            <Skeleton className='h-[90px] min-h-[90px] rounded-lg mt-3' />
            <Skeleton className='h-[90px] min-h-[90px] rounded-lg mt-3' />
            <Skeleton className='h-[90px] min-h-[90px] rounded-lg mt-3' />
          </>
        )}
        {cards.map((card, index) => (
          <KanbanCard
            key={`card-${card.name}-${index}`}
            card={card}
            cardStyle={getCardStyle(title)}
          />
        ))}
      </div>
    </div>
  );
};
