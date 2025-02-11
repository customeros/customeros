import { useRef, useCallback, useLayoutEffect } from 'react';

import { ElkNode, LayoutOptions } from 'elkjs/lib/elk-api';
import ELK, { ElkExtendedEdge } from 'elkjs/lib/elk.bundled';
import {
  Edge,
  Node,
  useStore,
  useReactFlow,
  ReactFlowState,
} from '@xyflow/react';

const elk = new ELK();

const elkOptions: LayoutOptions = {
  'elk.algorithm': 'layered',
  'elk.layered.spacing.nodeNodeBetweenLayers': '100',
  'elk.spacing.nodeNode': '80',
  'elk.direction': 'DOWN',
};

const ANIMATION_DURATION = 200;

const nodeCountSelector = (state: ReactFlowState) => state.nodeLookup.size;

export const useLayout = () => {
  const isInitialLayout = useRef(true);
  const animationFrameId = useRef<number | null>(null);

  const nodeCount = useStore(nodeCountSelector);
  const { getNodes, getEdges, setNodes, fitView } = useReactFlow();

  const getLayoutedElements = useCallback(
    async (
      nodes: Node[],
      edges: Edge[],
      options: LayoutOptions = elkOptions,
    ) => {
      const graph: ElkNode = {
        id: 'root',
        layoutOptions: options,
        children: nodes.map((node) => ({
          ...node,
          width: getNodeWidth(node).width || 150,
          height: getNodeWidth(node).height || 56,
        })),
        edges: edges as unknown as ElkExtendedEdge[],
      };

      const layoutResult = await elk.layout(graph);

      return layoutResult.children || [];
    },
    [],
  );

  const animateLayout = useCallback(
    (targetNodes: Node[]) => {
      const startNodes = getNodes();
      let start: number;

      const animate = (timestamp: number) => {
        if (!start) start = timestamp;
        const progress = (timestamp - start) / ANIMATION_DURATION;

        if (progress < 1) {
          const currentNodes = startNodes.map((startNode) => {
            const targetNode = targetNodes.find((n) => n.id === startNode.id);

            if (!targetNode) return startNode;

            return {
              ...startNode,
              position: {
                x:
                  startNode.position.x +
                  (targetNode.position.x - startNode.position.x) * progress,
                y:
                  startNode.position.y +
                  (targetNode.position.y - startNode.position.y) * progress,
              },
            };
          });

          setNodes(currentNodes);
          animationFrameId.current = requestAnimationFrame(animate);
        } else {
          setNodes(targetNodes);

          if (!isInitialLayout.current) {
            // fitView({ duration: 200, padding: 0.2, maxZoom: 1, minZoom: 0.5 });
          }
          isInitialLayout.current = false;
        }
      };

      animationFrameId.current = requestAnimationFrame(animate);
    },
    [getNodes, setNodes, fitView],
  );

  const organizeLayout = useCallback(async () => {
    const nodes = getNodes();
    const edges = getEdges();

    const layoutedNodes = await getLayoutedElements(nodes, edges);

    if (!layoutedNodes) return;

    const nodesToAnimate = layoutedNodes.map((node) => ({
      ...node,
      position: { x: node.x || 0, y: node.y || 0 },
    })) as Node[];

    animateLayout(nodesToAnimate);
  }, [getNodes, getEdges, getLayoutedElements, animateLayout]);

  useLayoutEffect(() => {
    organizeLayout();

    return () => {
      if (animationFrameId.current !== null) {
        cancelAnimationFrame(animationFrameId.current);
      }
    };
  }, [nodeCount, organizeLayout]);

  return { organizeLayout };
};

export const getNodeWidth = (node: Node) => {
  switch (node.type) {
    case 'trigger':
      return { width: 300, height: 83 };
    case 'control':
      return { width: 156, height: 56 };
    case 'wait':
      return { width: 156, height: 56 };
    case 'action':
      return { width: 300, height: 56 };
    default:
      return { width: 300, height: 56 };
  }
};
