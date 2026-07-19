import React, { useState } from 'react';

const SettingsPage: React.FC = () => {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const toggleTheme = () => setTheme(t => t === 'dark' ? 'light' : 'dark');

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-6)', maxWidth: '800px' }}>
      <h1 style={{ fontSize: 'var(--text-xl)', fontWeight: 600, margin: 0 }}>Settings</h1>

      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
        <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>General</h2>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontWeight: 500 }}>Theme preference</div>
            <div style={{ color: 'var(--color-text-secondary)', fontSize: 'var(--text-sm)' }}>Toggle between light and dark mode</div>
          </div>
          <button className="btn btn-secondary" onClick={toggleTheme}>
            {theme === 'dark' ? 'Switch to Light' : 'Switch to Dark'}
          </button>
        </div>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 'var(--space-2)' }}>
          <div>
            <div style={{ fontWeight: 500 }}>Organization Name</div>
          </div>
          <input type="text" className="input" placeholder="Organization" disabled />
        </div>
      </div>

      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
        <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>Integrations</h2>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontWeight: 500 }}>SIEM Endpoint</div>
          </div>
          <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
            <input type="text" className="input" placeholder="https://..." />
            <button className="btn btn-secondary">Save</button>
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontWeight: 500 }}>Slack Webhook (Alerts)</div>
          </div>
          <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
            <input type="text" className="input" placeholder="https://hooks.slack.com/..." />
            <button className="btn btn-secondary">Save</button>
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ fontWeight: 500 }}>OTLP Endpoint</div>
          </div>
          <span style={{ color: 'var(--color-text-secondary)' }}>{import.meta.env?.OTEL_EXPORTER_OTLP_ENDPOINT || 'otelcol:4317'}</span>
        </div>
      </div>

      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
        <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>Scanner</h2>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div><div style={{ fontWeight: 500 }}>Model A Endpoint</div></div>
          <span style={{ color: 'var(--color-text-secondary)' }}>{import.meta.env?.AEGIS_SCANNER_MODEL_A_URL || 'http://model-a:8000'}</span>
        </div>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div><div style={{ fontWeight: 500 }}>Model B Endpoint</div></div>
          <span style={{ color: 'var(--color-text-secondary)' }}>{import.meta.env?.AEGIS_SCANNER_MODEL_B_URL || 'http://model-b:8000'}</span>
        </div>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div><div style={{ fontWeight: 500 }}>Fail Mode</div></div>
          <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
            <label><input type="radio" name="failMode" value="open" /> Open</label>
            <label><input type="radio" name="failMode" value="closed" defaultChecked /> Closed</label>
          </div>
        </div>
      </div>

      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
        <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>Red Team</h2>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div><div style={{ fontWeight: 500 }}>Target URL</div></div>
          <input type="text" className="input" defaultValue="http://gateway:8080" />
        </div>
        
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div><div style={{ fontWeight: 500 }}>Schedule</div></div>
          <span className="badge badge--info">Every 6h</span>
        </div>
        
        <button className="btn btn-primary" style={{ alignSelf: 'flex-start' }}>Run Now</button>
      </div>

      <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
        <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>API Keys</h2>
        <div style={{ color: 'var(--color-text-secondary)', fontSize: 'var(--text-sm)' }}>
          Note: All credentials are managed by Vault. No secrets are stored in this UI.
        </div>
      </div>

    </div>
  );
};

export default SettingsPage;
