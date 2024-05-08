import type {
  DropResult,
  DraggableLocation,
  DroppableProvided,
} from '@hello-pangea/dnd';

import React, { Component, ReactElement } from 'react';

import styled from '@emotion/styled';
import { Droppable, DragDropContext } from '@hello-pangea/dnd';
import { PartialAutoScrollerOptions } from '@hello-pangea/dnd/src/state/auto-scroller/fluid-scroller/auto-scroller-options-types';

import Column from './KanbanColumn/Col';
import reorder, { Quote, QuoteMap, reorderQuoteMap } from './utils';

interface ParentContainerProps {
  height: string;
}

const ParentContainer = styled.div<ParentContainerProps>`
  height: ${({ height }) => height};
  overflow-x: hidden;
  overflow-y: auto;
`;

const Container = styled.div`
  background-color: 'red';
  min-height: 100vh;
  /* like display:flex but will allow bleeding over the window width */
  min-width: 100vw;
  display: inline-flex;
`;

interface Props {
  initial: QuoteMap;
  useClone?: boolean;
  containerHeight?: string;
  isCombineEnabled?: boolean;
  applyGlobalStyles?: boolean;
  withScrollableColumns?: boolean;
  autoScrollerOptions?: PartialAutoScrollerOptions;
}

interface State {
  columns: QuoteMap;
  ordered: string[];
}

export default class Board extends Component<Props, State> {
  /* eslint-disable react/sort-comp */
  static defaultProps = {
    isCombineEnabled: false,
    applyGlobalStyles: true,
  };

  state: State = {
    columns: this.props.initial,
    ordered: Object.keys(this.props.initial),
  };

  onDragEnd = (result: DropResult): void => {
    if (result.combine) {
      if (result.type === 'COLUMN') {
        const shallow: string[] = [...this.state.ordered];
        shallow.splice(result.source.index, 1);
        this.setState({ ordered: shallow });

        return;
      }

      const column: Quote[] = this.state.columns[result.source.droppableId];
      const withQuoteRemoved: Quote[] = [...column];
      withQuoteRemoved.splice(result.source.index, 1);
      const columns: QuoteMap = {
        ...this.state.columns,
        [result.source.droppableId]: withQuoteRemoved,
      };
      this.setState({ columns });

      return;
    }

    // dropped nowhere
    if (!result.destination) {
      return;
    }

    const source: DraggableLocation = result.source;
    const destination: DraggableLocation = result.destination;

    // did not move anywhere - can bail early
    if (
      source.droppableId === destination.droppableId &&
      source.index === destination.index
    ) {
      return;
    }

    // reordering column
    if (result.type === 'COLUMN') {
      const ordered: string[] = reorder(
        this.state.ordered,
        source.index,
        destination.index,
      );

      this.setState({
        ordered,
      });

      return;
    }

    const data = reorderQuoteMap({
      quoteMap: this.state.columns,
      source,
      destination,
    });

    this.setState({
      columns: data.quoteMap,
    });
  };

  render(): ReactElement {
    const columns: QuoteMap = this.state.columns;
    const ordered: string[] = this.state.ordered;
    const {
      containerHeight,
      useClone,
      isCombineEnabled,
      withScrollableColumns,
      applyGlobalStyles,
    } = this.props;

    const board = (
      <Droppable
        droppableId='board'
        type='COLUMN'
        direction='horizontal'
        ignoreContainerClipping={Boolean(containerHeight)}
        isCombineEnabled={isCombineEnabled}
      >
        {(provided: DroppableProvided) => (
          <Container ref={provided.innerRef} {...provided.droppableProps}>
            {ordered.map((key: string, index: number) => (
              <Column
                key={key}
                index={index}
                title={key}
                quotes={columns[key]}
                isScrollable={withScrollableColumns}
                isCombineEnabled={isCombineEnabled}
                useClone={useClone}
              />
            ))}
            {provided.placeholder}
          </Container>
        )}
      </Droppable>
    );

    return (
      <React.Fragment>
        <DragDropContext
          onDragEnd={this.onDragEnd}
          autoScrollerOptions={this.props.autoScrollerOptions}
        >
          {containerHeight ? (
            <ParentContainer height={containerHeight}>{board}</ParentContainer>
          ) : (
            board
          )}
        </DragDropContext>
      </React.Fragment>
    );
  }
}
