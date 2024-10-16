import {
  $isTextNode,
  $getSelection,
  LexicalEditor,
  createCommand,
  $setSelection,
  $createTextNode,
  $isRangeSelection,
  COMMAND_PRIORITY_HIGH,
  $createRangeSelection,
} from 'lexical';

import { SelectOption } from '@ui/utils/types';
import { $createVariableNode } from '@ui/form/Editor/nodes/VariableNode';

export const INSERT_VARIABLE_NODE = createCommand('INSERT_VARIABLE_NODE');

export function registerInsertVariableNodeCommand(
  editor: LexicalEditor,
  onAfterInsertVariableNode: () => void,
): () => void {
  return editor.registerCommand(
    INSERT_VARIABLE_NODE,
    (variable: SelectOption) => {
      editor.update(() => {
        const selection = $getSelection();

        if ($isRangeSelection(selection)) {
          const anchor = selection.anchor;
          const focus = selection.focus;
          const isBackward = selection.isBackward();

          const endPoint = isBackward ? anchor : focus;
          const endNode = endPoint.getNode();
          const endOffset = endPoint.offset;

          const variableNode = $createVariableNode(variable, true);
          const spaceBeforeNode = $createTextNode(' ');
          const spaceAfterNode = $createTextNode(' ');

          if ($isTextNode(endNode)) {
            if (endOffset === endNode.getTextContentSize()) {
              endNode.insertAfter(spaceBeforeNode);
            } else {
              const [, rightSplit] = endNode.splitText(endOffset);

              rightSplit.insertBefore(spaceBeforeNode);
            }
          } else {
            const parent = endNode.getParent();

            if (parent) {
              const index = parent.getChildren().indexOf(endNode);

              if (index !== -1) {
                parent.insertAfter(spaceBeforeNode);
              }
            }
          }

          spaceBeforeNode.insertAfter(variableNode);
          variableNode.insertAfter(spaceAfterNode);

          const nodeSelection = $createRangeSelection();

          nodeSelection.anchor.set(variableNode.getKey(), 0, 'element');
          nodeSelection.focus.set(variableNode.getKey(), 1, 'element');
          $setSelection(nodeSelection);
        }
        onAfterInsertVariableNode();
      });

      return true;
    },
    COMMAND_PRIORITY_HIGH,
  );
}
