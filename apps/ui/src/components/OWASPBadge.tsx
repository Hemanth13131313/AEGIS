import React from 'react';

export const OWASPBadge: React.FC<{ category: string }> = ({ category }) => {
  const getStyle = (cat: string) => {
    if (cat.includes('LLM01')) return 'bg-red-500/10 text-red-400 border-red-500/30';
    if (cat.includes('LLM06')) return 'bg-orange-500/10 text-orange-400 border-orange-500/30';
    if (cat.includes('LLM03')) return 'bg-purple-500/10 text-purple-400 border-purple-500/30';
    return 'bg-gray-500/10 text-gray-400 border-gray-500/30';
  };

  return (
    <span className={`px-2 py-1 rounded text-xs font-mono border ${getStyle(category)}`}>
      {category}
    </span>
  );
};
