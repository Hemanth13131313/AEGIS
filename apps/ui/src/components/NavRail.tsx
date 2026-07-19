import React from 'react';
import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Shield, AlertTriangle, Search, Target, FileCheck, Settings, ChevronLeft, ChevronRight } from 'lucide-react';

interface NavRailProps {
  expanded: boolean;
  onToggle: () => void;
}

const navItems = [
  { path: '/overview', label: 'Overview', icon: LayoutDashboard },
  { path: '/policies', label: 'Policies', icon: Shield },
  { path: '/detections', label: 'Detections', icon: AlertTriangle },
  { path: '/trace-explorer', label: 'Trace Explorer', icon: Search },
  { path: '/red-team', label: 'Red Team', icon: Target },
  { path: '/compliance', label: 'Compliance', icon: FileCheck },
  { path: '/settings', label: 'Settings', icon: Settings },
];

const NavRail: React.FC<NavRailProps> = ({ expanded, onToggle }) => {
  return (
    <div style={{
      width: expanded ? '240px' : '64px',
      backgroundColor: 'var(--color-bg-surface)',
      borderRight: '1px solid var(--color-border-subtle)',
      display: 'flex',
      flexDirection: 'column',
      transition: 'width var(--transition-normal)',
      zIndex: 10,
    }}>
      <div style={{ 
        height: '64px', 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: expanded ? 'flex-start' : 'center',
        padding: expanded ? '0 var(--space-4)' : '0',
        borderBottom: '1px solid var(--color-border-subtle)' 
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-2)' }}>
          <span style={{ fontSize: '24px' }}>🛡️</span>
          {expanded && <span style={{ fontWeight: 600, fontSize: 'var(--text-lg)', letterSpacing: '0.05em' }}>AEGIS</span>}
        </div>
      </div>

      <nav style={{ flex: 1, padding: 'var(--space-3) 0', display: 'flex', flexDirection: 'column', gap: 'var(--space-1)' }}>
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            style={({ isActive }) => ({
              display: 'flex',
              alignItems: 'center',
              justifyContent: expanded ? 'flex-start' : 'center',
              padding: expanded ? 'var(--space-3) var(--space-4)' : 'var(--space-3) 0',
              textDecoration: 'none',
              color: isActive ? 'var(--color-accent-primary)' : 'var(--color-text-secondary)',
              backgroundColor: isActive ? 'rgba(76, 141, 255, 0.1)' : 'transparent',
              borderLeft: `3px solid ${isActive ? 'var(--color-accent-primary)' : 'transparent'}`,
              transition: 'all var(--transition-fast)',
            })}
          >
            <item.icon size={20} />
            {expanded && <span style={{ marginLeft: 'var(--space-3)', fontSize: 'var(--text-sm)', fontWeight: 500 }}>{item.label}</span>}
          </NavLink>
        ))}
      </nav>

      <div style={{ padding: 'var(--space-3)', borderTop: '1px solid var(--color-border-subtle)', display: 'flex', justifyContent: expanded ? 'flex-end' : 'center' }}>
        <button 
          onClick={onToggle}
          style={{
            background: 'none',
            border: 'none',
            color: 'var(--color-text-secondary)',
            cursor: 'pointer',
            padding: 'var(--space-2)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderRadius: 'var(--radius-md)'
          }}
          className="hover:bg-var(--color-bg-surface-raised)"
        >
          {expanded ? <ChevronLeft size={20} /> : <ChevronRight size={20} />}
        </button>
      </div>
    </div>
  );
};

export default NavRail;
