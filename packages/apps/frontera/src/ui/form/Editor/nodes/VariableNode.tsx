import {
  TextNode,
  type Spread,
  type NodeKey,
  type LexicalNode,
  type EditorConfig,
  type DOMExportOutput,
  $applyNodeReplacement,
  type DOMConversionMap,
  type SerializedTextNode,
  type DOMConversionOutput,
} from 'lexical';

import { SelectOption } from '@ui/utils/types';

export type SerializedVariableNode = Spread<
  {
    variableId: string;
    variableName: string;
  },
  SerializedTextNode
>;

function $convertVariableElement(
  domNode: HTMLElement,
): DOMConversionOutput | null {
  const textContent = domNode.textContent;
  const id = domNode.getAttribute('data-variable-id');

  if (textContent !== null && id !== null) {
    const node = $createVariableNode({ label: textContent, value: id });

    return {
      node,
    };
  }

  return null;
}

export class VariableNode extends TextNode {
  __variable: SelectOption;

  static getType(): string {
    return 'variable';
  }

  static clone(node: VariableNode): VariableNode {
    return new VariableNode(node.__variable, node.__text, node.__key);
  }

  static importJSON(serializedNode: SerializedVariableNode): VariableNode {
    const node = $createVariableNode({
      label: serializedNode.variableName,
      value: serializedNode.variableId,
    });

    node.setTextContent(serializedNode.text);
    node.setFormat(serializedNode.format);
    node.setDetail(serializedNode.detail);
    node.setMode(serializedNode.mode);
    node.setStyle(serializedNode.style);

    return node;
  }

  constructor(variable: SelectOption, text?: string, key?: NodeKey) {
    super(text ?? `{{${variable.label}}}`, key);
    this.__variable = variable;
  }

  exportJSON(): SerializedVariableNode {
    return {
      ...super.exportJSON(),
      variableName: this.__variable.label,
      variableId: this.__variable.value,
      type: 'variable',
      version: 1,
    };
  }

  createDOM(config: EditorConfig): HTMLElement {
    const dom = super.createDOM(config);

    dom.className = 'variable text-gray-500';
    dom.setAttribute('data-variable-id', this.__variable.value);

    return dom;
  }

  exportDOM(): DOMExportOutput {
    const element = document.createElement('span');

    element.setAttribute('data-lexical-variable', 'true');
    element.textContent = this.__text;
    element.setAttribute('data-variable-id', this.__variable.value);

    return { element };
  }

  static importDOM(): DOMConversionMap | null {
    return {
      span: (domNode: HTMLElement) => {
        if (!domNode.hasAttribute('data-lexical-variable')) {
          return null;
        }

        return {
          conversion: $convertVariableElement,
          priority: 1,
        };
      },
    };
  }

  isTextEntity(): true {
    return true;
  }

  canInsertTextBefore(): boolean {
    return false;
  }

  canInsertTextAfter(): boolean {
    return false;
  }
}

export function $createVariableNode(variable: SelectOption): VariableNode {
  const variableNode = new VariableNode(variable);

  variableNode.setMode('segmented').toggleDirectionless();

  return $applyNodeReplacement(variableNode);
}

export function $isVariableNode(
  node: LexicalNode | null | undefined,
): node is VariableNode {
  return node instanceof VariableNode;
}
