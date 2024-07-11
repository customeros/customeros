import { useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import {
  Droppable,
  DragDropContext,
  OnDragEndResponder,
} from '@hello-pangea/dnd';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Columns02 } from '@ui/media/icons/Columns02';
import { TableViewType, ColumnViewType } from '@graphql/types';
import { Menu, MenuList, MenuGroup, MenuButton } from '@ui/overlay/Menu/Menu';

import { ColumnItem, DraggableColumnItem } from './ColumnItem';
import {
  invoicesOptionsMap,
  renewalsOptionsMap,
  contactsOptionsMap,
  invoicesHelperTextMap,
  renewalsHelperTextMap,
  contactsHelperTextMap,
  organizationsOptionsMap,
  organizationsHelperTextMap,
} from './columnOptions';

interface EditColumnsProps {
  type: TableViewType;
}

export const EditColumns = observer(({ type }: EditColumnsProps) => {
  const store = useStore();
  const [searchParams] = useSearchParams();
  const preset = searchParams?.get('preset');

  const [optionsMap, helperTextMap] = useMemo(() => {
    return [
      type === TableViewType.Contacts
        ? contactsOptionsMap
        : type === TableViewType.Invoices
        ? invoicesOptionsMap
        : type === TableViewType.Renewals
        ? renewalsOptionsMap
        : organizationsOptionsMap,
      type === TableViewType.Contacts
        ? contactsHelperTextMap
        : type === TableViewType.Invoices
        ? invoicesHelperTextMap
        : type === TableViewType.Renewals
        ? renewalsHelperTextMap
        : organizationsHelperTextMap,
    ];
  }, [type]);

  const tableViewDef = store.tableViewDefs.getById(preset ?? '0');

  const columns =
    tableViewDef?.value?.columns.map((c) => ({
      ...c,
      label: optionsMap[c.columnType],
      helperText: helperTextMap[c.columnType],
    })) ?? [];

  const handleDragEnd: OnDragEndResponder = (res) => {
    const sourceIndex = res.source.index;
    const destIndex = res?.destination?.index as number;
    const destination = res.destination;

    if (!destination) return;
    if (sourceIndex === destIndex) return;

    tableViewDef?.reorderColumn(sourceIndex, destIndex);
  };

  const pinnedColumns =
    tableViewDef?.value.tableType === TableViewType.Organizations
      ? columns.filter((e) =>
          [
            ColumnViewType.OrganizationsAvatar,
            ColumnViewType.OrganizationsName,
          ].includes(e.columnType),
        )
      : tableViewDef?.value.tableType === TableViewType.Contacts
      ? columns.filter((e) =>
          [ColumnViewType.ContactsAvatar, ColumnViewType.ContactsName].includes(
            e.columnType,
          ),
        )
      : [columns[0]];
  const showDraggable = (index: number) => {
    if (
      tableViewDef?.value &&
      [TableViewType.Organizations, TableViewType.Contacts].includes(
        tableViewDef.value.tableType,
      )
    ) {
      return index > 1;
    }

    return index > 0;
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
          <MenuList className='w-[350px]'>
            {pinnedColumns.map((col, i) => (
              <ColumnItem
                key={`${col?.columnType}-${i}`}
                isPinned
                noPointerEvents
                label={col?.label}
                visible={col?.visible}
                columnType={col?.columnType}
              />
            ))}

            <Droppable
              key='active-columns'
              droppableId='active-columns'
              renderClone={(provided, snapshot, rubric) => {
                return (
                  <ColumnItem
                    provided={provided}
                    snapshot={snapshot}
                    helperText={columns?.[rubric.source.index]?.helperText}
                    columnType={columns?.[rubric.source.index]?.columnType}
                    visible={columns?.[rubric.source.index]?.visible}
                    onCheck={() => {
                      tableViewDef?.update((value) => {
                        value.columns[rubric.source.index].visible =
                          !value?.columns?.[rubric.source.index]?.visible;

                        return value;
                      });
                    }}
                    label={columns[rubric.source.index]?.label}
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
                        showDraggable(index) && (
                          <DraggableColumnItem
                            index={index}
                            label={col?.label}
                            visible={col?.visible}
                            helperText={col?.helperText}
                            noPointerEvents={isDraggingOver}
                            key={col?.columnType}
                            onCheck={() => {
                              tableViewDef?.update((value) => {
                                value.columns[index].visible =
                                  !value.columns[index].visible;

                                return value;
                              });
                            }}
                            columnType={col?.columnType}
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
