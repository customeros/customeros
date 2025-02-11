import { useState, useEffect, useCallback } from 'react';

import { Edge, Node, useReactFlow } from '@xyflow/react';

type UseUndoRedoOptions = {
  maxHistorySize: number;
  enableShortcuts: boolean;
};

type UseUndoRedo = (options?: UseUndoRedoOptions) => {
  undo: () => void;
  redo: () => void;
  canUndo: boolean;
  canRedo: boolean;
  takeSnapshot: () => void;
};

type HistoryItem = {
  nodes: Node[];
  edges: Edge[];
};

const defaultOptions: UseUndoRedoOptions = {
  maxHistorySize: 100,
  enableShortcuts: true,
};

export const useUndoRedo: UseUndoRedo = ({
  maxHistorySize = defaultOptions.maxHistorySize,
  enableShortcuts = defaultOptions.enableShortcuts,
} = defaultOptions) => {
  // the past and future arrays store the states that we can jump to
  const [past, setPast] = useState<HistoryItem[]>([]);
  const [future, setFuture] = useState<HistoryItem[]>([]);

  const { setNodes, setEdges, getNodes, getEdges } = useReactFlow();

  const takeSnapshot = useCallback(() => {
    // push the current graph to the past state
    setPast((past) => [
      ...past.slice(past.length - maxHistorySize + 1, past.length),
      { nodes: getNodes(), edges: getEdges() },
    ]);

    // whenever we take a new snapshot, the redo operations need to be cleared to avoid state mismatches
    setFuture([]);
  }, [getNodes, getEdges, maxHistorySize]);

  const undo = useCallback(() => {
    // get the last state that we want to go back to
    const pastState = past[past.length - 1];

    if (pastState) {
      // first we remove the state from the history
      setPast((past) => past.slice(0, past.length - 1));
      // we store the current graph for the redo operation
      setFuture((future) => [
        ...future,
        { nodes: getNodes(), edges: getEdges() },
      ]);
      // now we can set the graph to the past state
      setNodes(pastState.nodes);
      setEdges(pastState.edges);
    }
  }, [setNodes, setEdges, getNodes, getEdges, past]);

  const redo = useCallback(() => {
    const futureState = future[future.length - 1];

    if (futureState) {
      setFuture((future) => future.slice(0, future.length - 1));
      setPast((past) => [...past, { nodes: getNodes(), edges: getEdges() }]);
      setNodes(futureState.nodes);
      setEdges(futureState.edges);
    }
  }, [setNodes, setEdges, getNodes, getEdges, future]);

  useEffect(() => {
    if (!enableShortcuts) {
      return;
    }

    const keyDownHandler = (event: KeyboardEvent) => {
      if (event.key === 'y' && (event.ctrlKey || event.metaKey)) {
        event.preventDefault();
        redo();
      } else if (event.key === 'z' && (event.ctrlKey || event.metaKey)) {
        undo();
      }
    };

    document.addEventListener('keydown', keyDownHandler);

    return () => {
      document.removeEventListener('keydown', keyDownHandler);
    };
  }, [undo, redo, enableShortcuts]);

  return {
    undo,
    redo,
    takeSnapshot,
    canUndo: !past.length,
    canRedo: !future.length,
  };
};

export default useUndoRedo;
