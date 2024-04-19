'use client';
import { observer } from 'mobx-react-lite';
import { Droppable, DragDropContext } from '@hello-pangea/dnd';

import { Button } from '@ui/form/Button/Button';
import { ColumnViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Columns02 } from '@ui/media/icons/Columns02';
import { Menu, MenuList, MenuGroup, MenuButton } from '@ui/overlay/Menu/Menu';

import { ColumnOption, ColumnOptionContent } from './ColumnOption';

type InvoicesColumnType =
  | ColumnViewType.InvoicesAmount
  | ColumnViewType.InvoicesBillingCycle
  | ColumnViewType.InvoicesContract
  | ColumnViewType.InvoicesDueDate
  | ColumnViewType.InvoicesIssueDatePast
  | ColumnViewType.InvoicesInvoicePreview
  | ColumnViewType.InvoicesIssueDate
  | ColumnViewType.InvoicesInvoiceStatus
  | ColumnViewType.InvoicesPaymentStatus;

const invoicesOptionsMap: Record<InvoicesColumnType, string> = {
  [ColumnViewType.InvoicesAmount]: 'Amount',
  [ColumnViewType.InvoicesBillingCycle]: 'Billing cycle',
  [ColumnViewType.InvoicesContract]: 'Contract',
  [ColumnViewType.InvoicesDueDate]: 'Due date',
  [ColumnViewType.InvoicesInvoicePreview]: 'Invoice preview',
  [ColumnViewType.InvoicesIssueDate]: 'Issue date',
  [ColumnViewType.InvoicesIssueDatePast]: 'Issue date',
  [ColumnViewType.InvoicesInvoiceStatus]: 'Invoice status',
  [ColumnViewType.InvoicesPaymentStatus]: 'Payment status',
};

export const EditColumns = observer(() => {
  const { tableViewDefsStore } = useStore();

  const tableViewDef = tableViewDefsStore.getById('5');

  const columns =
    tableViewDef?.value?.columns.map((c) => ({
      ...c,
      label: invoicesOptionsMap[c.columnType as InvoicesColumnType],
    })) ?? [];

  return (
    <>
      <Menu
        defaultOpen
        onOpenChange={(open) => {
          if (!open) {
            tableViewDef?.orderColumnsByVisibility();
          }
        }}
      >
        <MenuButton asChild>
          <Button size='xs' leftIcon={<Columns02 />}>
            Edit columns
          </Button>
        </MenuButton>
        <DragDropContext
          onDragEnd={(res) => {
            const sourceIndex = res.source.index;
            const destIndex = res?.destination?.index as number;

            if (sourceIndex === destIndex) return;

            tableViewDef?.reorderColumn(sourceIndex, destIndex);
          }}
        >
          <MenuList className='w-[300px]'>
            <Droppable
              key='active-columns'
              droppableId='active-columns'
              renderClone={(provided, snapshot, rubric) => {
                return (
                  <ColumnOptionContent
                    provided={provided}
                    snapshot={snapshot}
                    columnType={columns[rubric.source.index].columnType}
                    visible={columns[rubric.source.index].visible}
                    onCheck={() => {
                      tableViewDef?.update((value) => {
                        value.columns[rubric.source.index].visible =
                          !value.columns[rubric.source.index].visible;

                        return value;
                      });
                    }}
                    label={columns[rubric.source.index].label}
                  />
                );
              }}
            >
              {(provided, snapshot) => (
                <MenuGroup ref={provided.innerRef} {...provided.droppableProps}>
                  {columns.map((col, index) => (
                    <ColumnOption
                      index={index}
                      label={col.label}
                      visible={col?.visible}
                      key={col?.columnType}
                      onCheck={() => {
                        tableViewDef?.update((value) => {
                          value.columns[index].visible =
                            !value.columns[index].visible;

                          return value;
                        });
                      }}
                      columnType={col.columnType}
                    />
                  ))}
                  {provided.placeholder}
                </MenuGroup>
              )}
            </Droppable>
          </MenuList>
        </DragDropContext>
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
