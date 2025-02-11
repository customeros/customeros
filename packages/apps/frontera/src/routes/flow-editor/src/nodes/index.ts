import { WaitNode } from './WaitNode.tsx';
import { ActionNode } from './ActionNode.tsx';
import { ControlNode } from './ControlNode.tsx';
import { TriggerNode } from './TriggerNode.tsx';

export const nodeTypes = {
  trigger: TriggerNode,
  action: ActionNode,
  control: ControlNode,
  wait: WaitNode,
};
