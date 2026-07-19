import React from 'react';
import { SeverityBadge } from './SeverityBadge';
import { Severity } from '../lib/constants';

interface StatCardProps {
  title: string;
  value: string | number;
  trend?: {
    value: number;
    direction: 'up' | 'down';
    label: string;
  };
  severity?: Severity;
}

const StatCard: React.FC<StatCardProps> = ({ title, value, trend, severity }) => {
  return (
    <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <h3 style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)', fontWeight: 500, margin: 0 }}>
          {title}
        </h3>
        {severity && <SeverityBadge severity={severity} />}
      </div>
      
      <div style={{ fontSize: 'var(--text-xl)', fontWeight: 600, color: 'var(--color-text-primary)' }}>
        {value}
      </div>
      
      {trend && (
        <div style={{ 
          fontSize: 'var(--text-xs)', 
          color: trend.direction === 'up' && title.includes('Rate') 
            ? 'var(--color-severity-safe)' 
            : trend.direction === 'up' ? 'var(--color-severity-high)' : 'var(--color-severity-safe)',
          display: 'flex',
          alignItems: 'center',
          gap: 'var(--space-1)'
        }}>
          <span>{trend.direction === 'up' ? '↑' : '↓'} {trend.value}%</span>
          <span style={{ color: 'var(--color-text-disabled)' }}>{trend.label}</span>
        </div>
      )}
    </div>
  );
};

export default StatCard;
