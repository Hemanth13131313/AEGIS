import React, { useState } from 'react';
import { Shield, Search, User, Menu } from 'lucide-react';
import NavRail from './NavRail';

interface AppLayoutProps {
  children: React.ReactNode;
  theme: 'dark' | 'light';
  toggleTheme: () => void;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children, theme, toggleTheme }) => {
  const [isRailExpanded, setIsRailExpanded] = useState(true);

  return (
    <div style={{ display: 'flex', height: '100vh', overflow: 'hidden' }}>
      <NavRail expanded={isRailExpanded} onToggle={() => setIsRailExpanded(!isRailExpanded)} />
      
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
        {/* Top bar */}
        <header style={{
          height: '64px',
          borderBottom: '1px solid var(--color-border-subtle)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '0 var(--space-4)',
          backgroundColor: 'var(--color-bg-surface)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-3)' }}>
            {!isRailExpanded && (
              <button 
                onClick={() => setIsRailExpanded(true)}
                style={{ background: 'none', border: 'none', color: 'var(--color-text-primary)', cursor: 'pointer' }}
              >
                <Menu size={20} />
              </button>
            )}
            <div style={{ 
              padding: 'var(--space-1) var(--space-2)', 
              borderRadius: 'var(--radius-md)', 
              backgroundColor: 'var(--color-bg-surface-raised)',
              fontSize: 'var(--text-sm)',
              fontWeight: 500,
              display: 'flex',
              alignItems: 'center',
              gap: 'var(--space-2)'
            }}>
              <Shield size={16} color="var(--color-accent-primary)" />
              Acme Corp / Prod Gateway
            </div>
          </div>
          
          <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-4)' }}>
            <div style={{ position: 'relative' }}>
              <Search size={16} style={{ position: 'absolute', left: '10px', top: '50%', transform: 'translateY(-50%)', color: 'var(--color-text-secondary)' }} />
              <input 
                type="text" 
                placeholder="Search..." 
                style={{
                  padding: 'var(--space-2) var(--space-3) var(--space-2) 36px',
                  borderRadius: 'var(--radius-md)',
                  border: '1px solid var(--color-border-subtle)',
                  backgroundColor: 'var(--color-bg-canvas)',
                  color: 'var(--color-text-primary)',
                  fontSize: 'var(--text-sm)'
                }}
              />
            </div>
            
            <button className="btn btn-secondary" onClick={toggleTheme}>
              {theme === 'dark' ? '☀️ Light' : '🌙 Dark'}
            </button>
            
            <button style={{ 
              background: 'var(--color-bg-surface-raised)', 
              border: '1px solid var(--color-border-subtle)', 
              borderRadius: '50%', 
              width: '32px', 
              height: '32px', 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              cursor: 'pointer',
              color: 'var(--color-text-primary)'
            }}>
              <User size={16} />
            </button>
          </div>
        </header>

        {/* Main Content */}
        <main style={{ flex: 1, overflowY: 'auto', padding: 'var(--space-6)' }}>
          {children}
        </main>
      </div>
    </div>
  );
};

export default AppLayout;
