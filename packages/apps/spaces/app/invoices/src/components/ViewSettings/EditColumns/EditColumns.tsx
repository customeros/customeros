'use client';
import { observer } from 'mobx-react-lite';
import {
  Droppable,
  DragDropContext,
  OnDragEndResponder,
} from '@hello-pangea/dnd';

import { Button } from '@ui/form/Button/Button';
import { ColumnViewType } from '@graphql/types';
import { useStore } from '@shared/hooks/useStore';
import { Columns02 } from '@ui/media/icons/Columns02';
import { Menu, MenuList, MenuGroup, MenuButton } from '@ui/overlay/Menu/Menu';

import { ColumnItem, DraggableColumnItem } from './ColumnItem';

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

const invoicesHelperTextMap: Record<InvoicesColumnType, string> = {
  [ColumnViewType.InvoicesAmount]: 'E.g. $6,450',
  [ColumnViewType.InvoicesBillingCycle]: 'E.g. Monthly',
  [ColumnViewType.InvoicesContract]: 'E.g. Pile Contract',
  [ColumnViewType.InvoicesDueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesInvoicePreview]: 'E.g. RKD-04025',
  [ColumnViewType.InvoicesIssueDate]: 'E.g. 15 Aug 2019',
  [ColumnViewType.InvoicesIssueDatePast]: 'E.g. 15 Jun 2019',
  [ColumnViewType.InvoicesInvoiceStatus]: 'E.g. Scheduled',
  [ColumnViewType.InvoicesPaymentStatus]: 'E.g. Paid',
};

export const EditColumns = observer(() => {
  const { tableViewDefsStore } = useStore();

  const tableViewDef = tableViewDefsStore.getById('5');

  const columns =
    tableViewDef?.value?.columns.map((c) => ({
      ...c,
      label: invoicesOptionsMap[c.columnType as InvoicesColumnType],
      helperText: invoicesHelperTextMap[c.columnType as InvoicesColumnType],
    })) ?? [];

  const handleDragEnd: OnDragEndResponder = (res) => {
    const sourceIndex = res.source.index;
    const destIndex = res?.destination?.index as number;
    const destination = res.destination;

    if (!destination) return;
    if (sourceIndex === destIndex) return;

    tableViewDef?.reorderColumn(sourceIndex, destIndex);
  };

  return (
    <>
      <Menu
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
        <DragDropContext onDragEnd={handleDragEnd}>
          <MenuList className='w-[300px]'>
            <ColumnItem
              isPinned
              noPointerEvents
              label={columns[0].label}
              visible={columns[0].visible}
              columnType={columns[0].columnType}
            />
            <Droppable
              key='active-columns'
              droppableId='active-columns'
              renderClone={(provided, snapshot, rubric) => {
                return (
                  <ColumnItem
                    provided={provided}
                    snapshot={snapshot}
                    helperText={columns[rubric.source.index].helperText}
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
              {(provided, { isDraggingOver }) => (
                <>
                  <MenuGroup
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                  >
                    {columns.map(
                      (col, index) =>
                        index > 0 && (
                          <DraggableColumnItem
                            index={index}
                            label={col.label}
                            visible={col?.visible}
                            helperText={col.helperText}
                            noPointerEvents={isDraggingOver}
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
                        ),
                    )}
                    {provided.placeholder}
                  </MenuGroup>
                </>
              )}
            </Droppable>
          </MenuList>
        </DragDropContext>
      </Menu>
    </>
  );
});
