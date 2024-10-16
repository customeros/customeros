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
    isNewlyInserted: boolean;
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
  __isNewlyInserted: boolean;

  static getType(): string {
    return 'variable';
  }

  static clone(node: VariableNode): VariableNode {
    return new VariableNode(
      node.__variable,
      node.__isNewlyInserted,
      node.__text,
      node.__key,
    );
  }

  static importJSON(serializedNode: SerializedVariableNode): VariableNode {
    const node = $createVariableNode(
      {
        label: serializedNode.variableName,
        value: serializedNode.variableId,
      },
      serializedNode.isNewlyInserted,
    );

    node.setTextContent(serializedNode.text);
    node.setFormat(serializedNode.format);
    node.setDetail(serializedNode.detail);
    node.setMode(serializedNode.mode);
    node.setStyle(serializedNode.style);

    return node;
  }

  constructor(
    variable: SelectOption,
    isNewlyInserted: boolean = true,
    text?: string,
    key?: NodeKey,
  ) {
    super(text ?? `{{${variable.label}}}`, key);
    this.__variable = variable;
    this.__isNewlyInserted = isNewlyInserted;
  }

  exportJSON(): SerializedVariableNode {
    return {
      ...super.exportJSON(),
      variableName: this.__variable.label,
      variableId: this.__variable.value,
      isNewlyInserted: this.__isNewlyInserted,
      type: 'variable',
      version: 1,
    };
  }

  createDOM(config: EditorConfig): HTMLElement {
    const dom = super.createDOM(config);

    dom.className = `variable border-dotted ${
      this.__isNewlyInserted
        ? ' text-gray-400  hover:text-gray-400 newly-inserted'
        : ' text-gray-500  hover:text-gray-700 hover:border-gray-700'
    }`;
    dom.setAttribute('data-variable-id', this.__variable.value);
    dom.setAttribute('data-lexical-variable', 'true');

    return dom;
  }

  updateDOM(prevNode: VariableNode, dom: HTMLElement): boolean {
    const isUpdated = super.updateDOM(prevNode, dom, {
      namespace: 'data-lexical-variable',
      theme: {},
    });

    if (
      prevNode.__variable.value !== this.__variable.value ||
      prevNode.__isNewlyInserted !== this.__isNewlyInserted
    ) {
      dom.setAttribute('data-variable-id', this.__variable.value);
      dom.className = `variable border-dotted ${
        this.__isNewlyInserted
          ? ' text-gray-400 hover:text-gray-400 newly-inserted'
          : ' text-gray-500 hover:text-gray-700 hover:border-gray-700'
      }`;

      return true;
    }

    return isUpdated;
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

  setSelected() {
    if (this.__isNewlyInserted) {
      this.__isNewlyInserted = false;
      this.getWritable().__isNewlyInserted = false;
    }
  }

  getLength(): number {
    return this.__text.length;
  }
}

export function $createVariableNode(
  variable: SelectOption,
  isNewlyInserted: boolean = true,
): VariableNode {
  const variableNode = new VariableNode(variable, isNewlyInserted);

  variableNode.setMode('segmented').toggleDirectionless();

  return $applyNodeReplacement(variableNode);
}

export function $isVariableNode(
  node: LexicalNode | null | undefined,
): node is VariableNode {
  return node instanceof VariableNode;
}
