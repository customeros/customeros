import React, { Component, ReactElement } from 'react';

import styled from '@emotion/styled';

const Container = styled.div`
  display: flex;
  flex-direction: column;
`;

import type {
  DroppableProvided,
  DraggableProvided,
  DroppableStateSnapshot,
  DraggableStateSnapshot,
} from '@hello-pangea/dnd';

import { Droppable, Draggable } from '@hello-pangea/dnd';

import { Quote } from '../utils';

export const getBackgroundColor = (
  isDraggingOver: boolean,
  isDraggingFrom: boolean,
): string => {
  if (isDraggingOver) {
    return 'lightblue';
  }
  if (isDraggingFrom) {
    return 'red';
  }

  return 'yellow';
};

interface WrapperProps {
  isDraggingOver: boolean;
  isDraggingFrom: boolean;
  isDropDisabled: boolean;
}

const scrollContainerHeight = 250;

const DropZone = styled.div`
  /* stop the list collapsing when empty */
  min-height: ${scrollContainerHeight}px;

  /*
    not relying on the items for a margin-bottom
    as it will collapse when the list is empty
  */
`;

interface Props {
  title?: string;
  listId?: string;
  quotes: Quote[];
  listType?: string;
  useClone?: boolean;
  internalScroll?: boolean;
  isDropDisabled?: boolean;
  isCombineEnabled?: boolean;
  // may not be provided - and might be null
  ignoreContainerClipping?: boolean;
}

interface QuoteListProps {
  quotes: Quote[];
}

function InnerQuoteList(props: QuoteListProps): ReactElement {
  return (
    <>
      {props.quotes.map((quote: Quote, index: number) => (
        <Draggable key={quote.id} draggableId={quote.id} index={index}>
          {(
            dragProvided: DraggableProvided,
            dragSnapshot: DraggableStateSnapshot,
          ) => (
            <div
              key={quote.id}
              quote={quote}
              isDragging={dragSnapshot.isDragging}
              isGroupedOver={Boolean(dragSnapshot.combineTargetFor)}
              provided={dragProvided}
            >
              ABC
            </div>
          )}
        </Draggable>
      ))}
    </>
  );
}

const InnerQuoteListMemo = React.memo<QuoteListProps>(InnerQuoteList);

interface InnerListProps {
  quotes: Quote[];
  dropProvided: DroppableProvided;
  title: string | undefined | null;
}

function InnerList(props: InnerListProps) {
  const { quotes, dropProvided } = props;
  const title = props.title ? <Title>{props.title}</Title> : null;

  return (
    <Container>
      {title}
      <DropZone ref={dropProvided.innerRef}>
        <InnerQuoteListMemo quotes={quotes} />
        {dropProvided.placeholder}
      </DropZone>
    </Container>
  );
}

export function QuoteList(props: Props): ReactElement {
  const {
    ignoreContainerClipping,
    internalScroll,
    scrollContainerStyle,
    isDropDisabled,
    isCombineEnabled,
    listId = 'LIST',
    listType,
    style,
    quotes,
    title,
    useClone,
  } = props;

  return (
    <Droppable
      droppableId={listId}
      type={listType}
      ignoreContainerClipping={ignoreContainerClipping}
      isDropDisabled={isDropDisabled}
      isCombineEnabled={isCombineEnabled}
      renderClone={
        useClone
          ? (provided, snapshot, descriptor) => (
              <QuoteItem
                quote={quotes[descriptor.source.index]}
                provided={provided}
                isDragging={snapshot.isDragging}
                isClone
              />
            )
          : undefined
      }
    >
      {(
        dropProvided: DroppableProvided,
        dropSnapshot: DroppableStateSnapshot,
      ) => (
        <div
          style={style}
          isDraggingOver={dropSnapshot.isDraggingOver}
          isDropDisabled={Boolean(isDropDisabled)}
          isDraggingFrom={Boolean(dropSnapshot.draggingFromThisWith)}
          {...dropProvided.droppableProps}
        >
          {internalScroll ? (
            <div style={scrollContainerStyle}>
              <InnerList
                quotes={quotes}
                title={title}
                dropProvided={dropProvided}
              />
            </div>
          ) : (
            <InnerList
              quotes={quotes}
              title={title}
              dropProvided={dropProvided}
            />
          )}
        </div>
      )}
    </Droppable>
  );
}

interface HeaderProps {
  isDragging: boolean;
}

const Header = styled.div<HeaderProps>`
  display: flex;
  align-items: center;
`;

interface Props {
  title: string;
  index: number;
  quotes: Quote[];
  useClone?: boolean;
  isScrollable?: boolean;
  isCombineEnabled?: boolean;
}

export default class Column extends Component<Props> {
  render(): ReactElement {
    const title: string = this.props.title;
    const quotes: Quote[] = this.props.quotes;
    const index: number = this.props.index;

    return (
      <Draggable draggableId={title} index={index}>
        {(provided: DraggableProvided, snapshot: DraggableStateSnapshot) => (
          <Container ref={provided.innerRef} {...provided.draggableProps}>
            <Header isDragging={snapshot.isDragging}>
              <div
                {...provided.dragHandleProps}
                aria-label={`${title} quote list`}
              >
                {title}
              </div>
            </Header>
            <QuoteList
              listId={title}
              listType='QUOTE'
              style={{
                backgroundColor: snapshot.isDragging ? 'red' : undefined,
              }}
              quotes={quotes}
              internalScroll={this.props.isScrollable}
              isCombineEnabled={Boolean(this.props.isCombineEnabled)}
              useClone={Boolean(this.props.useClone)}
            />
          </Container>
        )}
      </Draggable>
    );
  }
}
