'use client';
import { Portal } from '@radix-ui/react-dropdown-menu';
import { observer, useLocalObservable } from 'mobx-react-lite';
import { Droppable, Draggable, DragDropContext } from '@hello-pangea/dnd';

import { ColumnDef } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Columns02 } from '@ui/media/icons/Columns02';
import { Checkbox } from '@ui/form/Checkbox/Checkbox2';
import { HandleDrag } from '@ui/media/icons/HandleDrag';
import { mockedTableDefs } from '@shared/util/tableDefs.mock';
import {
  Menu,
  MenuList,
  MenuItem,
  MenuGroup,
  MenuLabel,
  MenuButton,
} from '@ui/overlay/Menu/Menu';

import { ColumnOption } from './ColumnOption';

const allColumns = mockedTableDefs[4].columns;

export const EditColumns = observer(() => {
  const { tableViewDefsStore } = useStore();

  const tableViewDef = tableViewDefsStore.getById('5');
  const activeColumns = tableViewDef?.value?.columns ?? [];

  const state = useLocalObservable(() => ({
    unusedColumns: [] as ColumnDef[],
    // setAllColumns(columns: ColumnDef[]) {
    //   this.allColumns = columns;
    // },
  }));

  return (
    <>
      <Menu defaultOpen>
        <MenuButton asChild>
          <Button size='xs' leftIcon={<Columns02 />}>
            Edit columns
          </Button>
        </MenuButton>
        <MenuList>
          <DragDropContext
            onDragEnd={(res) => {
              const sourceId = res.source.droppableId;
              const destId = res.destination?.droppableId;
              const sourceIndex = res.source.index;
              const destIndex = res?.destination?.index as number;

              if (sourceId === destId) {
                if (sourceIndex === destIndex) return;
                if (!destId) return;

                tableViewDef?.reorderColumn(sourceIndex, destIndex);
              }

              // console.log(res);
            }}
          >
            <Droppable droppableId='active-columns'>
              {(provided, snapshot) => (
                <div ref={provided.innerRef} {...provided.droppableProps}>
                  {activeColumns.map((col, index) => (
                    <Draggable
                      key={col?.id}
                      index={index}
                      draggableId={col?.id as string}
                    >
                      {(provided) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          {...provided.dragHandleProps}
                        >
                          {col?.columnType?.name}
                        </div>
                      )}
                    </Draggable>
                  ))}
                </div>
              )}
            </Droppable>
          </DragDropContext>
        </MenuList>
      </Menu>
    </>
  );
});

// const EmptyContent = ({ dropId }: { dropId: string }) => {
//   const { setNodeRef } = useDroppable({
//     id: dropId,
//   });

//   return (
//     <div ref={setNodeRef} className='w-full h-10'>
//       Add colums here
//     </div>
//   );
// };
