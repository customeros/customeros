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

export type SerializedHashtagNode = Spread<
  {
    hashtagId: string;
    hashtagName: string;
  },
  SerializedTextNode
>;

function $convertHashtagElement(
  domNode: HTMLElement,
): DOMConversionOutput | null {
  const textContent = domNode.textContent;
  const id = domNode.getAttribute('data-hashtag-id');

  if (textContent !== null && id !== null) {
    const node = $createHashtagNode({ label: textContent, value: id });

    return {
      node,
    };
  }

  return null;
}

export class HashtagNode extends TextNode {
  __hashtag: SelectOption;

  static getType(): string {
    return 'hashtag';
  }

  static clone(node: HashtagNode): HashtagNode {
    return new HashtagNode(node.__hashtag, node.__text, node.__key);
  }
  static importJSON(serializedNode: SerializedHashtagNode): HashtagNode {
    const node = $createHashtagNode({
      label: serializedNode.hashtagName,
      value: serializedNode.hashtagId,
    });
    node.setTextContent(serializedNode.text);
    node.setFormat(serializedNode.format);
    node.setDetail(serializedNode.detail);
    node.setMode(serializedNode.mode);
    node.setStyle(serializedNode.style);

    return node;
  }

  constructor(hashtag: SelectOption, text?: string, key?: NodeKey) {
    super(text ?? hashtag.label, key);
    this.__hashtag = hashtag;
  }

  exportJSON(): SerializedHashtagNode {
    return {
      ...super.exportJSON(),
      hashtagName: this.__hashtag.label,
      hashtagId: this.__hashtag.value,
      type: 'hashtag',
      version: 1,
    };
  }

  createDOM(config: EditorConfig): HTMLElement {
    const dom = super.createDOM(config);
    dom.className = 'hashtag text-primary-600';
    dom.setAttribute('data-hashtag-id', this.__hashtag.value);

    return dom;
  }

  exportDOM(): DOMExportOutput {
    const element = document.createElement('span');
    element.setAttribute('data-lexical-hashtag', 'true');
    element.textContent = this.__text;
    element.setAttribute('data-hashtag-id', this.__hashtag.value);

    return { element };
  }

  static importDOM(): DOMConversionMap | null {
    return {
      span: (domNode: HTMLElement) => {
        if (!domNode.hasAttribute('data-lexical-hashtag')) {
          return null;
        }

        return {
          conversion: $convertHashtagElement,
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

export function $createHashtagNode(hashtag: SelectOption): HashtagNode {
  const hashtagNode = new HashtagNode(hashtag);
  hashtagNode.setMode('segmented').toggleDirectionless();

  return $applyNodeReplacement(hashtagNode);
}

export function $isHashtagNode(
  node: LexicalNode | null | undefined,
): node is HashtagNode {
  return node instanceof HashtagNode;
}
