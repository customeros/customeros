import { Edge, Node } from '@xyflow/react';
import { ElkNode, LayoutOptions } from 'elkjs/lib/elk-api';
import ELK, { ElkExtendedEdge } from 'elkjs/lib/elk.bundled.js';

const elk = new ELK();

// Elk has a *huge* amount of options to configure. To see everything you can
// tweak check out:
//
// - https://www.eclipse.org/elk/reference/algorithms.html
// - https://www.eclipse.org/elk/reference/options.html
const elkOptions = {
  'elk.algorithm': 'layered',
  'elk.layered.spacing.nodeNodeBetweenLayers': '100',
  'elk.spacing.nodeNode': '80',
  'elk.direction': 'DOWN',
};
interface useLayoutedElementsReturn {
  getLayoutedElements: (
    nodes: Node[],
    edges: Edge[],
    options?: LayoutOptions,
  ) => Promise<unknown>;
}

export const useLayoutedElements = (): useLayoutedElementsReturn => {
  const getLayoutedElements: useLayoutedElementsReturn['getLayoutedElements'] =
    (nodes, edges, options = elkOptions) => {
      const isHorizontal = options?.['elk.direction'] === 'RIGHT';
      const graph: ElkNode = {
        id: 'root',
        layoutOptions: options,
        children: nodes.map((node) => ({
          ...node,
          // Adjust the target and source handle positions based on the layout
          // direction.
          targetPosition: isHorizontal ? 'left' : 'top',
          sourcePosition: isHorizontal ? 'right' : 'bottom',

          // Hardcode a width and height for elk to use when layouting.
          width: node.width ?? getNodeWidth(node).width,
          height: node.height ?? getNodeWidth(node).height,
          properties: {
            'org.eclipse.elk.portConstraints': 'FIXED_ORDER',
          },
        })),
        edges: edges as unknown as ElkExtendedEdge[],
      };

      return elk
        .layout(graph)
        .then((layoutedGraph) => ({
          nodes: layoutedGraph?.children?.map((node) => ({
            ...node,
            // React Flow expects a position property on the node instead of `x`
            // and `y` fields.
            position: { x: node.x, y: node.y },
          })),

          edges: layoutedGraph.edges,
        }))
        .catch(console.error);
    };

  return { getLayoutedElements };
};

export const getNodeWidth = (node: Node) => {
  if (node.type === 'trigger') {
    return {
      width: 300,
      height: 56,
    };
  }

  if (node.type === 'control') {
    return {
      width: 131,
      height: 56,
    };
  }

  if (node.type === 'wait') {
    return {
      width: 150,
      height: 56,
    };
  }

  if (node.type === 'action') {
    return {
      width: 300,
      height: 56,
    };
  }

  return {
    width: 200,
    height: 75,
  };
};
