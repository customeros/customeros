import { WaitNode } from './WaitNode.tsx';
import { ActionNode } from './ActionNode';
import { ControlNode } from './ControlNode';
import { TriggerNode } from './TriggerNode';

export const nodeTypes = {
  trigger: TriggerNode,
  action: ActionNode,
  control: ControlNode,
  wait: WaitNode,
};
