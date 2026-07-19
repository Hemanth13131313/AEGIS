import React from 'react';
import { X, EyeOff, Tag, Check } from 'lucide-react';

type Action = 'allow' | 'block' | 'redact' | 'tag';

export const ActionBadge: React.FC<{ action: Action }> = ({ action }) => {
  const config = {
    allow: { icon: Check, color: 'text-green-500 bg-green-500/10 border-green-500/20' },
    block: { icon: X, color: 'text-red-500 bg-red-500/10 border-red-500/20' },
    redact: { icon: EyeOff, color: 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20' },
    tag: { icon: Tag, color: 'text-orange-500 bg-orange-500/10 border-orange-500/20' },
  };

  const { icon: Icon, color } = config[action];

  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-xs font-medium uppercase tracking-wide ${color}`}>
      <Icon className="w-3.5 h-3.5" />
      {action}
    </span>
  );
};
