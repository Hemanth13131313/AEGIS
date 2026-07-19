import React, { useState } from 'react';

const CompliancePage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'slo' | 'eu-ai-act' | 'iso-42001' | 'evidence'>('slo');

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-6)' }}>
      <h1 style={{ fontSize: 'var(--text-xl)', fontWeight: 600, margin: 0 }}>Compliance & SLO</h1>

      <div style={{ display: 'flex', gap: 'var(--space-4)', borderBottom: '1px solid var(--color-border-subtle)', paddingBottom: 'var(--space-2)' }}>
        {['slo', 'eu-ai-act', 'iso-42001', 'evidence'].map((tab) => (
          <button 
            key={tab} 
            style={{ 
              background: 'none', border: 'none', cursor: 'pointer',
              fontWeight: activeTab === tab ? 'bold' : 'normal',
              color: activeTab === tab ? 'var(--color-text)' : 'var(--color-text-secondary)',
              padding: 'var(--space-2)'
            }}
            onClick={() => setActiveTab(tab as any)}
          >
            {tab === 'slo' ? 'SLO Dashboard' : tab === 'eu-ai-act' ? 'EU AI Act' : tab === 'iso-42001' ? 'ISO 42001' : 'Evidence Bundle'}
          </button>
        ))}
      </div>

      {activeTab === 'slo' && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 'var(--space-4)' }}>
          {[
            { title: 'Gateway Latency', target: 'P95 ≤10ms', current: '6ms', status: '✔' },
            { title: 'Scanner Latency', target: 'P95 ≤500ms', current: '312ms', status: '✔' },
            { title: 'Availability', target: '≥99.9%', current: '99.97%', status: '✔' },
            { title: 'False Positive Rate', target: '<2%', current: '0.8%', status: '✔' },
            { title: 'Red Team Coverage', target: '>80%', current: '96%', status: '✔' }
          ].map(slo => (
            <div key={slo.title} className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
              <div style={{ fontSize: 'var(--text-md)', fontWeight: 600 }}>{slo.title}</div>
              <div style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)' }}>Target: {slo.target}</div>
              <div style={{ fontSize: 'var(--text-2xl)', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 'var(--space-2)' }}>
                {slo.current} <span style={{ color: 'green', fontSize: 'var(--text-lg)' }}>{slo.status}</span>
              </div>
              <div style={{ height: '4px', background: 'green', borderRadius: '2px', width: '80%' }}></div>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'eu-ai-act' && (
        <div className="card">
          <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, marginBottom: 'var(--space-4)' }}>EU AI Act Coverage</h2>
          <table style={{ width: '100%', textAlign: 'left', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--color-border-subtle)' }}>
                <th style={{ padding: 'var(--space-2)' }}>Article</th>
                <th style={{ padding: 'var(--space-2)' }}>Implementation</th>
                <th style={{ padding: 'var(--space-2)' }}>Status</th>
              </tr>
            </thead>
            <tbody>
              {[
                { art: 'Article 9 (Risk Management System)', imp: 'OPA policy engine', stat: '✔ Implemented' },
                { art: 'Article 10 (Data Governance)', imp: 'RAG monitor covers data quality', stat: '⚠ Partial' },
                { art: 'Article 11 (Technical Documentation)', imp: 'ADRs + this dashboard', stat: '✔ Implemented' },
                { art: 'Article 13 (Transparency)', imp: 'Trace explorer', stat: '✔ Implemented' },
                { art: 'Article 15 (Accuracy & Robustness)', imp: 'Red team automation', stat: '✔ Implemented' },
                { art: 'Article 17 (Quality Management)', imp: 'CI/CD covers versioning', stat: '⚠ Partial' }
              ].map(row => (
                <tr key={row.art} style={{ borderBottom: '1px solid var(--color-border-subtle)' }}>
                  <td style={{ padding: 'var(--space-2)' }}>{row.art}</td>
                  <td style={{ padding: 'var(--space-2)' }}>{row.imp}</td>
                  <td style={{ padding: 'var(--space-2)' }}>{row.stat}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <button className="btn btn-secondary" style={{ marginTop: 'var(--space-4)' }}>Generate Compliance Report</button>
        </div>
      )}

      {activeTab === 'iso-42001' && (
        <div className="card">
          <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, marginBottom: 'var(--space-4)' }}>ISO 42001:2023 Coverage</h2>
          <div style={{ color: 'var(--color-text-secondary)' }}>Coverage details for ISO 42001 are similarly structured to the EU AI Act.</div>
        </div>
      )}

      {activeTab === 'evidence' && (
        <div className="card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)', alignItems: 'flex-start' }}>
          <h2 style={{ fontSize: 'var(--text-lg)', fontWeight: 600, margin: 0 }}>Evidence Bundle</h2>
          <div style={{ color: 'var(--color-text-secondary)' }}>Exports policy versions, red team results, SLO data, and ADRs as a signed ZIP archive.</div>
          <button className="btn btn-primary">Generate Evidence Bundle</button>
          <div style={{ fontSize: 'var(--text-sm)', color: 'var(--color-text-secondary)' }}>Last generated: Never</div>
        </div>
      )}
    </div>
  );
};

export default CompliancePage;
